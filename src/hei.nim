import std / [parseopt, os]

import shell
import strutils
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
      elif ["i", "A", "q", "e", "p"].contains(key):
        shell:
          "nix-env " & commandLineParams().join(" ")
        system.quit()
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
  echo """
    Options:
    -d, --dryrun                     Don't change anything; perform dry run
    -D, --debug                      Show trace on nix errors
    -f, --flake URI                  Change target flake to URI
    -h, --help                       Display this help, or help for a specific command
    -i, -A, -q, -e, -p               Forward to nix-env

"""
