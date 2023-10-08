import std/[
  os,
  parseopt,
  re,
  strformat,
  strutils,
  tables,
  tempfiles
]

let backupSuffix = ".nix-store-backup"

var dispatchTable = initTable[string, proc(flakePath: string, args: seq[string])]()

dispatchTable["help"] = proc (flakePath: string, args: seq[string]) =
  echo &"Forwarding to {args[0]} --help"
  dispatchTable[args[0]](flakePath, @["--help"])

dispatchTable["build"] = proc(flakePath: string, args: seq[string]) =
  var argstr = args.join(" ")
  if argstr == "": argstr = "."
  let res = execShellCmd &"nix build -L {argstr}"
  system.quit(res)

dispatchTable["check"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix flake check " & flakePath
  system.quit(res)

dispatchTable["show"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix flake show " & flakePath
  system.quit(res)

dispatchTable["gc"] = proc(flakePath: string, args: seq[string]) =
  var all, sys = false
  for kind, key, val in getopt(args):
    case kind
    of cmdLongOption, cmdShortOption:
      case key
      of "a", "all": all = true
      of "s", "system": sys = true
      of "h", "help":
        echo "Usage: dotfiles gc [options]"
        echo ""
        echo "  -a, --all     Clean up all profiles"
        echo "  -s, --system  Clean up the system profile"
        system.quit(0)
    of cmdArgument:
      continue
    of cmdEnd:
      assert(false)
  if all or sys:
    echo "Cleaning up your system profile"
    discard execShellCmd "sudo nix-collect-garbage -d"
    discard execShellCmd "sudo nix-store --optimise"
    discard execShellCmd "sudo nix-env --delete-generations old --profile /nix/var/nix/profiles/system"
    discard execShellCmd "sudo /nix/var/nix/profiles/system/bin/s,witch-to-configuration switch"
  if all or not sys:
    discard execShellCmd "nix-collect-garbage -d"

  system.quit(0)

dispatchTable["repl"] = proc(flakePath: string, args: seq[string]) =
  let (tmpfile, path) = createTempFile("dotfiles-repl.nix", "_end.tmp")
  tmpfile.write("import " & flakePath & "(builtins.getFlake \"" & flakePath & "\")")
  let res = execShellCmd "nix repl \\<nixpkgs\\> " & path
  system.quit(res)

dispatchTable["search"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix search nixpkgs " & args.join(" ")
  system.quit(res)

dispatchTable["update"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix flake update " & flakePath
  system.quit(res)

dispatchTable["rebuild"] = proc(flakePath: string, args: seq[string]) =
  if hostOs == "macosx":
    let res = execShellCmd "darwin-rebuild switch --flake " & flakePath
    system.quit(res)
  else:
    let res = execShellCmd "sudo nixos-rebuild switch --flake " & flakePath
    system.quit(res)

dispatchTable["swap"] = proc(flakePath: string, args: seq[string]) =
  for target in args:
    if target.dirExists:
      var foundBackups: seq[string] = @[]
      echo fmt"Checking {target}"
      for path in walkDirRec(target, checkDir = true, yieldFilter = {pcLinkToFile}):
        if path.contains(re".nix-store-backup$"):
          foundBackups.add(path.replace(backupSuffix, ""))
      if foundBackups.len > 0:
        echo "Backups found, swapping back"
        dispatchTable["swap"](flakePath, foundBackups)
      else:
        var targets: seq[string] = @[]
        for path in walkDirRec(target, checkDir = true, yieldFilter = {pcLinkToFile}):
          targets.add(path)
        dispatchTable["swap"](flakePath, targets)
    elif fmt"{target}{backupSuffix}".fileExists:
      echo &"Unswapping {target}"
      discard execShellCmd &"mv -i {target}{backupSuffix} {target}"
    elif target.fileExists:
      if target.symlinkExists and target.expandSymlink.contains(re"^/nix/"):
        echo &"Swapping {target}"
        discard execShellCmd &"mv {target} {target}{backupSuffix}"
        discard execShellCmd &"cp {target}{backupSuffix} {target}"
      else:
        echo &"Not swapping {target} because it is not in the nix store"
    else:
      echo &"No such file or directory: {target}"
      system.quit(1)

dispatchTable["ssh"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "ssh " & args[0] & " hei " & args[1..args.high].join(" ")
  system.quit(res)

dispatchTable["test"] = proc(flakePath: string, args: seq[string]) =
  dispatchTable["rebuild"](flakePath, @["--fast"])

dispatchTable["upgrade"] = proc(flakePath: string, args: seq[string]) =
  discard execShellCmd "nix flake update " & flakePath
  dispatchTable["rebuild"](flakePath, args)

proc dispatchCommand*(cmd: string, flakePath: string, args: seq[string]) =
  if dispatchTable.hasKey(cmd):
    dispatchTable[cmd](flakePath, args)
  else:
    echo("Unknown command: "&cmd)
