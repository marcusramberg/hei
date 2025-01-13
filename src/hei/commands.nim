import std/[os, parseopt, re, sequtils, strformat, strutils, tables, envvars]

type CommandProc = proc(flakePath: string, args: seq[string]): int

type Command = object
  description: string
  arg: string
  body: CommandProc

let backupSuffix = ".nix-store-backup"

var dispatchTable = initOrderedTable[string, Command]()

template makeCommand(name: string, help: string, args: string, command: untyped) =
  let cmd: CommandProc = command
  dispatchTable[name] = Command(description: help, arg: args, body: cmd)

proc exec(cmd: string): int =
  if getEnv("HEI_TESTING") == "1":
    return execShellCmd &"echo {cmd}"
  echo &"Running: {cmd}"
  return execShellCmd cmd

makeCommand(
  "help", help = "Call help with a command to get more info", args = "[SUBCOMMAND]"
):
  proc(flakePath: string, args: seq[string]): int =
    if args.len > 0:
      echo &"Forwarding to {args[0]} --help"
      return dispatchTable[args[0]].body(flakePath, @["--help"])
    echo """
    usage:  hei [global-options] [command] [sub-options]

    A simpler nix experience (inspired by hey by hlissner)

    Available commands: """

    for cmd in dispatchTable.keys:
      echo fmt"  {cmd:<12}  {dispatchTable[cmd].arg:<15}  {dispatchTable[cmd].description}"
    echo """

    Options:
        -d, --dry-run                     Don't change anything; perform dry run
        -D, --debug                      Show trace on nix errors
        -f, --flake URI                  Change target flake to URI
        -h, --help                       Display this help, or help for a specific command
        -i, -A, -q, -e, -p               Forward to nix-env

    """

makeCommand("build", help = "Run build with full logs", args = "<TARGET|.>"):
  proc(flakePath: string, args: seq[string]): int =
    var buildCommand = "nix"
    if exec("which nom") == 0:
      buildCommand = "nom"
    var argStr = args.join(" ")
    if argStr == "":
      argStr = "."
    exec &"{buildCommand} build -L {argStr}"

makeCommand("check", help = "Run 'nix flake check' on your flake", args = ""):
  proc(flakePath: string, args: seq[string]): int =
    exec "nix flake check " & flakePath

makeCommand(
  "completions", help = "Generate shell completions for hei", args = "[zsh|bash|fish]"
):
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
      echo &"complete -c hei {isRoot} -s d -l dry-run -d \"Don't change anything; perform dry run\""
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

makeCommand("gen", help = "Work with generations", args = "list|diff|show|switch"):
  proc(flakePath: string, args: seq[string]): int =
    for kind, key, val in getOpt(args):
      case kind
      of cmdArgument:
        case key
        of "list":
          return exec "sudo nix-env --list-generations --profile /nix/var/nix/profiles/system"
        of "delete":
          if args.len != 2:
            echo "Usage: dotfiles gen delete <generation>"
            return 1
          return exec "sudo nix-env --delete-generations --profile /nix/var/nix/profiles/system " &
            args[1]
        of "diff":
          echo "diff"
          return 1
        of "help":
          break
      of cmdLongOption, cmdShortOption:
        case key
        of "h", "help":
          break
      of cmdEnd:
        assert(false)
    echo """Usage: dotfiles gen <command>

  list    List generations
  delete  Delete a generation
  diff    Diff two generations"""
    return 0

makeCommand("gc", help = "Garbage collect & optimize nix store", args = "[-a] [-s]"):
  proc(flakePath: string, args: seq[string]): int =
    var all, sys = false
    for kind, key, val in getOpt(args):
      case kind
      of cmdLongOption, cmdShortOption:
        case key
        of "a", "all":
          all = true
        of "s", "system":
          sys = true
        of "h", "help":
          echo """Usage: dotfiles gc [options]
  -a, --all     Clean up all profiles"
  -s, --system  Clean up the system profile
"""
          return 0
      of cmdArgument:
        continue
      of cmdEnd:
        assert(false)
    if all or sys:
      echo "Cleaning up your system profile"
      discard exec "sudo nix-collect-garbage -d"
      discard exec "sudo nix-store --optimise"
      discard exec "sudo nix-env --delete-generations old --profile /nix/var/nix/profiles/system"
      discard
        exec "sudo /nix/var/nix/profiles/system/bin/s,witch-to-configuration switch"
    if all or not sys:
      discard exec "nix-collect-garbage -d"
    return 0

makeCommand("p", help = "nix profile commands", args =""):
  proc(flakePath: string, args: seq[string]): int =
    system.quit execShellCmd "nix profile " & args.join(" ")

makeCommand(
  "rebuild",
  help = "Rebuild the current system's flake",
  args = "[-o|--offline] [-r|--rollback] [-f|--fast] [switch|boot]",
):
  proc(flakePath: string, args: seq[string]): int =
    var
      options = &"--flake {flakePath}"
      argument = "switch"
    if args.len > 0:
      for kind, key, val in getOpt(args):
        case kind
        of cmdLongOption, cmdShortOption:
          case key
          of "o", "offline":
            options &= " --option substitute false"
          of "r", "rollback":
            options &= " --rollback"
          of "f", "fast":
            options &= " --fast"
          of "h", "help":
            echo """Usage: rebuild [options]
    [-o|--offline]
    [-r|--rollback]
    [-f|--fast] [switch|boot]
  """
            return 0
          else:
            var sep = if kind == cmdShortOption: "-" else: "--"
            options &= &" {sep}{key}"
        of cmdArgument:
          if val != "":
            argument = val
        of cmdEnd:
          assert(false)
    let rebuildCommand =
      if hostOs == "macosx": "darwin-rebuild" else: "sudo nixos-rebuild"
    exec &"{rebuildCommand} {argument} {options}"

makeCommand("repl", help = "Open a nix-repl with our system flake preloaded", args = ""):
  proc(flakePath: string, args: seq[string]): int =
    exec "nix repl --file " & flakePath & "/flake.nix"

makeCommand("rollback", help = "Roll back to previous generation", args = ""):
  proc(flakePath: string, args: seq[string]): int =
    dispatchTable["rebuild"].body(flakePath, @["--rollback"])

makeCommand("search", help = "Search nixpkgs for a package", args = "[package]"):
  proc(flakePath: string, args: seq[string]): int =
    exec "nix search nixpkgs " & args.join(" ")

makeCommand("show", help = "Show your flake", args = "[ARGS...]"):
  proc(flakePath: string, args: seq[string]): int =
    exec "nix flake show " & flakePath

makeCommand(
  "ssh", help = "Run a hei command on a remote NixOS system", args = "HOST [COMMAND]"
):
  proc(flakePath: string, args: seq[string]): int =
    exec "ssh " & args[0] & " hei " & args[1 .. args.high].join(" ")

makeCommand(
  "swap",
  help = "Recursively swap nix-store symlinks with copies (or back)",
  args = "PATH [PATH...]",
):
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
        discard exec &"mv -i {target}{backupSuffix} {target}"
      elif target.fileExists:
        if target.symlinkExists and target.expandSymlink.contains(re"^/nix/"):
          echo &"Swapping {target}"
          discard exec &"mv {target} {target}{backupSuffix}"
          discard exec &"cp {target}{backupSuffix} {target}"
        else:
          echo &"Not swapping {target} because it is not in the nix store"
      else:
        echo &"No such file or directory: {target}"
        return 1
    0

makeCommand("test", help = "Quickly rebuild, for quick iteration", args = ""):
  proc(flakePath: string, args: seq[string]): int =
    dispatchTable["rebuild"].body(flakePath, @["--fast"])

makeCommand("upgrade", help = "Update all flakes and rebuild system", args = ""):
  proc(flakePath: string, args: seq[string]): int =
    var rebuildArgs: seq[string] = @[]
    for kind, key, val in getOpt(args):
      case kind
      of cmdLongOption, cmdShortOption:
        case key
        of "h", "help":
          echo """Usage: dotfiles upgrade [options]
  -p, --pull     Pull from git before updating
  -h, --help     Display this help
          """
          return 0
        of "p", "pull":
          discard exec &"git -C {flakePath} pull --rebase"
        else:
          var sep = if kind == cmdShortOption: "-" else: "--"
          rebuildArgs.add(&"{sep}{key}")
      of cmdArgument:
        rebuildArgs.add(val)
      of cmdEnd:
        assert(false)
    if exec("nix flake update --flake " & flakePath) == 0:
      return dispatchTable["rebuild"].body(flakePath, rebuildArgs)
    echo "Update failed, not rebuilding."
    1

makeCommand(
  "update", help = "Update specific flakes or all of them", args = "[ INPUT...]"
):
  proc(flakePath: string, args: seq[string]): int =
    if args.len == 0:
      return exec "nix flake update " & flakePath
    return exec "nix flake lock " & flakePath &
      join(
        map(
          args,
          proc(arg: string): string =
            fmt" --update-input {arg}",
        ),
        " ",
      )

proc dispatchCommand*(cmd: string, flakePath: string, args: seq[string]) =
  if dispatchTable.hasKey(cmd):
    quit(dispatchTable[cmd].body(flakePath, args))
  echo(&"\n  Unknown command: {cmd}\n")
  system.quit dispatchTable["help"].body(flakePath, @[])
