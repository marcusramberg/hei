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

type CommandProc = proc (flakePath: string, args: seq[string]): int

type Command = object
  description: string
  arg: string
  body: CommandProc

let backupSuffix = ".nix-store-backup"

var dispatchTable = initOrderedTable[string, Command]()

template makeCommand(name: string, help: string, args: string,
    command: untyped) =
  let cmd: CommandProc = command
  dispatchTable[name] = Command(description: help, arg: args, body: cmd)

makeCommand("help",
  help = "Call help with a command to get more info",
  args = "[SUBCOMMAND]"):
  proc(flakePath: string, args: seq[string]): int =
    if args.len > 0:
      echo &"Forwarding to {args[0]} --help"
      return dispatchTable[args[0]].body(flakePath, @["--help"])
    echo """
    usage:  hei [global-options] [command] [sub-options]

    A simpler nix experience (inspired by hey by hlissner)

    Note: `hei` can also be used as a shortcut for nix-env:

      hei -q
      hei -iA nixos.htop
      hei -e htop

    Available commands: """

    for cmd in dispatchTable.keys:
      echo fmt"  {cmd:<12}  {dispatchTable[cmd].arg:<15}  {dispatchTable[cmd].description}"
    echo """

    Options:
        -d, --dryrun                     Don't change anything; perform dry run
        -D, --debug                      Show trace on nix errors
        -f, --flake URI                  Change target flake to URI
        -h, --help                       Display this help, or help for a specific command
        -i, -A, -q, -e, -p               Forward to nix-env

    """

makeCommand("build",
  help = "Run build with full logs",
  args = "<TARGET|.>"):
  proc(flakePath: string, args: seq[string]): int =
    var buildCommand = "nix"
    if execShellCmd("which nom") == 0: buildCommand = "nom"
    var argstr = args.join(" ")
    if argstr == "": argstr = "."
    execShellCmd &"{buildCommand} build -L {argstr}"

makeCommand("check",
  help = "Run 'nix flake check' on your flake",
  args = ""):
  proc(flakePath: string, args: seq[string]): int =
    execShellCmd "nix flake check " & flakePath

makeCommand("completions",
  help = "Generate shell completions for hei",
  args = "[zsh|bash|fish]"):
  proc(flakePath: string, args: seq[string]): int =
    if args.len == 0:
      echo "Usage: hei completion [zsh|bash|fish]"
      return 1
    case args[0]
      of "fish":
        let commands = dispatchTable.keys.toSeq.join(" ")
        let isRoot = &" -n \"not __fish_seen_subcommand_from {commands}\""
        echo "complete -c hei -f"
        for cmd in dispatchTable.keys:
          echo &"complete -c hei {isRoot} -a \"{cmd}\" -d \"{dispatchTable[cmd].description}\""
        # FIXME: Generate these from code
        echo &"complete -c hei {isRoot} -s d -l dryrun -d \"Don't change anything; perform dry run\""
        echo &"complete -c hei {isRoot} -s D -l debug -d \"Show trace on nix errors\""
        echo &"complete -c hei {isRoot} -s f -l flake -d \"Change target flake to URI\""
        echo &"complete -c hei {isRoot} -s h -l help -d \"Display this help, or help for a specific command\""
        echo &"complete -c hei {isRoot} -s i -s A -s q -s e -s p -d \"Forward to nix-env\""
      of "zsh":
        echo "compctl -k '(hei help | awk \"{print $2}\")' hei"
      of "bash":
        echo "complete -W \"$(hei help | awk '{print $2}')\" hei"
      else:
        echo "Unknown shell: {args[0]}"
        return 1

makeCommand("gen",
  help = "Work with generations",
  args = "list|diff|show|switch"):
  proc(flakePath: string, args: seq[string]): int =
    for kind, key, val in getopt(args):
      case kind
      of cmdArgument:
        case key
        of "list":
          return execShellCmd "sudo nix-env --list-generations --profile /nix/var/nix/profiles/system"
        of "delete":
          return execShellCmd "sudo nix-env --delete-generations --profile /nix/var/nix/profiles/system " & val
        of "diff":
          echo "diff"
          return 1
        of "help": break
      of cmdLongOption, cmdShortOption:
        case key
        of "h", "help": break
      of cmdEnd:
        assert(false)
    echo "Usage: dotfiles gen <command>"
    echo ""
    echo "  list    List generations"
    echo "  delete  Delete a generation"
    echo "  diff    Diff two generations"
    return 0

makeCommand("gc",
  help = "Garbage collect & optimize nix store",
  args = "[-a] [-s]"):
  proc(flakePath: string, args: seq[string]): int =
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
          return 0
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
    return 0

makeCommand("rebuild",
  help = "Rebuild the current system's flake",
  args = ""):
  proc(flakePath: string, args: seq[string]): int =
    let rebuildCommand = if hostOs == "macosx": "darwin-rebuild" else: "sudo nixos-rebuild"
    if args.len == 0:
      return execShellCmd &"{rebuildCommand} switch --flake {flakePath}"
    execShellCmd &"""{rebuildCommand} {args.join(" ")} --flake {flakePath}"""

makeCommand("repl",
  help = "Open a nix-repl with nixpkgs and dotfiles preloaded",
  args = ""):
  proc(flakePath: string, args: seq[string]): int =
    let (tmpfile, path) = createTempFile("dotfiles-repl.nix", "_end.tmp")
    tmpfile.write("import " & flakePath & "(builtins.getFlake \"" & flakePath & "\")")
    execShellCmd "nix repl \\<nixpkgs\\> " & path

makeCommand("rollback",
  help = "Roll back to previous generation",
  args = ""):
  proc(flakePath: string, args: seq[string]): int =
    dispatchTable["rebuild"].body(flakePath, @["--rollback", "switch"])

makeCommand("search",
  help = "Search nixpkgs for a package",
  args = "[package]"):
  proc(flakePath: string, args: seq[string]): int =
    execShellCmd "nix search nixpkgs " & args.join(" ")

makeCommand("show",
  help = "Show your flake",
  args = "[ARGS...]"):
  proc(flakePath: string, args: seq[string]): int =
    execShellCmd "nix flake show " & flakePath

makeCommand("ssh",
  help = "Run a hei command on a remote NixOS system",
  args = "HOST [COMMAND]"):
  proc(flakePath: string, args: seq[string]): int =
    execShellCmd "ssh " & args[0] & " hei " & args[1..args.high].join(" ")

makeCommand("swap",
  help = "Recursively swap nix-store symlinks with copies (or back)",
  args = "PATH [PATH...]"):
  proc(flakePath: string, args: seq[string]): int =
    for target in args:
      if target.dirExists:
        var foundBackups: seq[string] = @[]
        echo fmt"Checking {target}"
        for path in walkDirRec(target, checkDir = true, yieldFilter = {pcLinkToFile}):
          if path.contains(re".nix-store-backup$"):
            foundBackups.add(path.replace(backupSuffix, ""))
        if foundBackups.len > 0:
          echo "Backups found, swapping back"
          discard dispatchTable["swap"].body(flakePath, foundBackups)
        else:
          var targets: seq[string] = @[]
          for path in walkDirRec(target, checkDir = true, yieldFilter = {pcLinkToFile}):
            targets.add(path)
          discard dispatchTable["swap"].body(flakePath, targets)
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
        return 1
    0

makeCommand("test",
  help = "Quickly rebuild, for quick iteration",
  args = ""):
  proc(flakePath: string, args: seq[string]): int =
    dispatchTable["rebuild"].body(flakePath, @["--fast", "switch"])

makeCommand("upgrade",
  help = "Update all flakes and rebuild system",
  args = ""):
  proc(flakePath: string, args: seq[string]): int =
    if execShellCmd("nix flake update " & flakePath) == 0:
      return dispatchTable["rebuild"].body(flakePath, args)
    echo "Update failed, not rebuilding."
    1

makeCommand("update",
  help = "Update specific flakes or all of them",
  args = "[ INPUT...]"):
  proc(flakePath: string, args: seq[string]): int =
    if args.len == 0:
      return execShellCmd "nix flake update " & flakePath
    return execShellCmd "nix flake lock " & flakePath &
      join(map(args, proc(arg: string): string = fmt" --update-input {arg}"), " ")

proc dispatchCommand*(cmd: string, flakePath: string, args: seq[string]) =
  if dispatchTable.hasKey(cmd):
    quit(dispatchTable[cmd].body(flakePath, args))
  echo(&"\n  Unknown command: {cmd}\n")
  system.quit dispatchTable["help"].body(flakePath, @[])
