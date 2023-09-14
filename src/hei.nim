import std / [parseopt, os, envvars]

import strformat, strutils, sequtils
import utils
let flakePaths = @["/etc/nixos", "~/.config/nix-darwin"]

if existsEnv("NIX_SYSTEM_FLAKE"):
  let envFlake = getEnv("NIX_SYSTEM_FLAKE")
  flakePaths.insert(@[envFlake])
let flake = findFlake(flakePaths)
if flake == "":
  echo "No system flake found in ", flakePaths
  quit()

type Command = object

  name: string
  description: string
  args: string

const commands: seq[Command] = @[
  Command(name: "check", description: "Run 'nix flake check' on your dotfiles"),
  Command(name: "gc", description: "Garbage collect & optimize nix store"),
  Command(name: "help [SUBCOMMAND]", description: "Show usage information for this script or a subcommand"),
  Command(name: "generations", description: "Explore, manage, diff across generations"),
  Command(name: "info", args: "REPO [QUERY]",
      description: "Retrieve details (including SHA) for a REPO."),
  Command(name: "rebuild", description: "Rebuild the current system's flake"),
  Command(name: "repl", description: "Open a nix-repl with nixpkgs and dotfiles preloaded"),
  Command(name: "rollback", description: "Roll back to last generation"),
  Command(name: "search", description: "Search nixpkgs for a package"),
  Command(name: "show", args: "[ARGS...]", description: "Show your flake"),
  Command(name: "ssh", args: "HOST [COMMAND]",
      description: "Run a bin/hey command on a remote NixOS system"),
  Command(name: "swap", args: "PATH [PATH...]",
      description: "Recursively swap nix-store symlinks with copies (or back)."),
  Command(name: "test", description: "Quickly rebuild, for quick iteration"),
  Command(name: "theme", args: "THEME_NAME",
      description: "Quickly swap to another theme module"),
  Command(name: "upgrade", description: "Update all flakes and rebuild system"),
  Command(name: "update", args: "[ INPUT...]",
      description: "Update specific flakes or all of them"),
]

when isMainModule:
  for kind, key, value in getOpt():
    case kind
    of cmdEnd: break
    of cmdShortOption, cmdLongOption:
      if key == "v" or key == "version":
        echo "hei 0.0.1"
        quit()
      elif key == "h" or key == "help":
        break
        # run nix-env with the original command line arguments
      elif ["i", "A", "q", "e", "p"].contains(key):
        let res = execShellCmd "nix-env " & commandLineParams().join(" ")
        echo "Fowarding to nix-env ..."
        system.quit(res)
      else:
        echo "Unknown option: ", key, ". Run `hei` for help."
        quit()
    of cmdArgument:
      echo "command: ", key
  echo """hei - Welcome to a simpler nix experience

  Note: `hei` can also be used as a shortcut for nix-env:
    hei -q

    hei -iA nixos.htop
    hei -e htop

    Available commands: """
  for cmd in commands:
    echo fmt"{cmd.name:<10} - {cmd.description}"
  echo """
    Options:
    -d, --dryrun                     Don't change anything; perform dry run
    -D, --debug                      Show trace on nix errors
    -f, --flake URI                  Change target flake to URI
    -h, --help                       Display this help, or help for a specific command
    -i, -A, -q, -e, -p               Forward to nix-env

"""
