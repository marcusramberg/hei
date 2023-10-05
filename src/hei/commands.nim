import std/[parseopt, os, strutils, tables, tempfiles]
var dispatchTable = initTable[string, proc(flakePath: string, args: seq[string])]()

dispatchTable["help"] = proc (flakePath: string, args: seq[string]) =
  echo "help"
dispatchTable["build"] = proc(flakePath: string, args: seq[string]) =

  let res = execShellCmd "nix build -L " & args[0]
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
    of cmdArgument:
      echo "arg"
    of cmdEnd:
      assert(false)
  if all or sys:
    echo "Cleaning up your system profile"
    discard execShellCmd "sudo nix-collect-garbage -d"
    discard execShellCmd "sudo nix-store --optimise"
    discard execShellCmd "sudo nix-env --delete-generations old --profile /nix/var/nix/profiles/system"
    discard execShellCmd "sudo /nix/var/nix/profiles/system/bin/switch-to-configuration switch"
  if all and not sys:
    discard execShellCmd "nix-collect-garbage -d"
  system.quit(0)

dispatchTable["repl"] = proc(flakePath: string, args: seq[string]) =
  let (tmpfile, path) = createTempFile("dotfiles-repl.nix", "_end.tmp")

  tmpfile.write("import " & flakePath & "(builtins.getFlake \"" & flakePath & "\")")
  let res = execShellCmd "nix repl \\<nixpkgs\\> " & path
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
