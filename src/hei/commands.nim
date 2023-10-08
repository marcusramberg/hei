import std/[
  os,
  parseopt,
  re,
  sequtils,
  strformat,
  strutils,
  tables,
  tempfiles
]

type Command = object
  name: string
  description: string
  args: string

const commandsHelp: seq[Command] = @[
  Command(name: "build", description: "Run build with full logs"),
  Command(name: "check", description: "Run 'nix flake check' on your dotfiles"),
  Command(name: "gc", description: "Garbage collect & optimize nix store"),
  Command(name: "help", args: "[SUBCOMMAND]",
      description: "Show usage information for this script or a subcommand"),
  Command(name: "generations", description: "Explore, manage, diff across generations"),
  Command(name: "info", args: "REPO [QUERY]",
      description: "Retrieve details (including SHA) for a REPO."),
  Command(name: "rebuild", description: "Rebuild the current system's flake"),
  Command(name: "repl", description: "Open a nix-repl with nixpkgs and dotfiles preloaded"),
  Command(name: "rollback", description: "Roll back to last generation"),
  Command(name: "search", description: "Search nixpkgs for a package"),
  Command(name: "show", args: "[ARGS...]", description: "Show your flake"),
  Command(name: "ssh", args: "HOST [COMMAND]",
      description: "Run a hei command on a remote NixOS system"),
  Command(name: "swap", args: "PATH [PATH...]",
      description: "Recursively swap nix-store symlinks with copies (or back)."),
  Command(name: "test", description: "Quickly rebuild, for quick iteration"),
  # Command(name: "theme", args: "THEME_NAME",
    #     description: "Quickly swap to another theme module"),
  Command(name: "upgrade", description: "Update all flakes and rebuild system"),
  Command(name: "update", args: "[ INPUT...]",
      description: "Update specific flakes or all of them"),
]
let backupSuffix = ".nix-store-backup"

var dispatchTable = initTable[string, proc(flakePath: string, args: seq[string])]()

dispatchTable["help"] = proc (flakePath: string, args: seq[string]) =
  if args.len > 0:
    echo &"Forwarding to {args[0]} --help"
    dispatchTable[args[0]](flakePath, @["--help"])
    return
echo """
 usage:  hei [global-options] [command] [sub-options]

Welcome to a simpler nix experience (inspired by hey by hlissner)

Note: `hei` can also be used as a shortcut for nix-env:

  hei -q
  hei -iA nixos.htop
  hei -e htop


Available commands: """

for cmd in commandsHelp:
  echo fmt"  {cmd.name:<12}  {cmd.args:<15}  {cmd.description}"
echo """

 Options:
    -d, --dryrun                     Don't change anything; perform dry run
    -D, --debug                      Show trace on nix errors
    -f, --flake URI                  Change target flake to URI
    -h, --help                       Display this help, or help for a specific command
    -i, -A, -q, -e, -p               Forward to nix-env

"""

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

  if args.len == 0:
    let res = execShellCmd "nix flake update " & flakePath & " " & args[0]
    system.quit(res)
  else:
    let res = execShellCmd "nix flake lock " & flakePath &
      join(map(args, proc(arg: string): string = fmt" --update-input {arg}"), " ")
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
