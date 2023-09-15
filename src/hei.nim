import std / [parseopt, os, envvars]

import strformat, strutils, sequtils
import hei/[commands, utils]

var flakePaths = @["/etc/nixos", "~/.config/nix-darwin"]
if existsEnv("NIX_SYSTEM_FLAKE"):
  let envFlake = getEnv("NIX_SYSTEM_FLAKE")
  flakePaths.insert(@[envFlake])
var flakePath = findFlake(flakePaths)
if flakePath == "":
  echo "No system flake found in ", flakePaths
  echo "If your system flake is in a different location, set the NIX_SYSTEM_FLAKE environment variable"
  quit()

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

when isMainModule:
  var p = initOptParser()
  while true:
    p.next()
    case p.kind
      of cmdend: break
      of cmdshortoption, cmdlongoption:
        if p.key == "f" or p.key == "flake":
          flakePath = p.val
        elif p.key == "d" or p.key == "dryrun":
          putEnv("NIX_DEBUG", "1")
        elif p.key == "D" or p.key == "debug":
          putEnv("NIX_DEBUG", "1")
          putEnv("NIX_SHOW_TRACE", "1")
        elif p.key == "v" or p.key == "version":
          echo fmt"hei 0.0.1 - running on {hostos}({hostcpu})"
          quit()
        elif p.key == "h" or p.key == "help":
          break
        elif ["i", "a", "q", "e", "p"].contains(p.key):
          # run nix-env with the original command line arguments
          echo "forwarding to nix-env ..."
          let res = execshellcmd "nix-env " & commandlineparams().join(" ")
          system.quit(res)
        else:
          echo "Unknown option: ", p.key, ". run `hei` for help."
          quit()
      of cmdargument:
        dispatchcommand(p.key, flakePath, p.remainingArgs)
        quit()

  echo """Error: No command specified.

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
