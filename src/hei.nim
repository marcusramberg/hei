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

when isMainModule:
  var p = initOptParser()
  while true:
    p.next()
    case p.kind
      of cmdend: break
      of cmdshortoption, cmdlongoption:
        case p.key
        of "f", "flake":
          flakePath = p.val
        of "d", "dryrun":
          putEnv("NIX_DEBUG", "1")
        of "D", "debug":
          putEnv("NIX_DEBUG", "1")
          putEnv("NIX_SHOW_TRACE", "1")
        of "v", "version":
          echo fmt"hei 0.0.1 - running on {hostos}({hostcpu})"
          quit()
        of "h", "help":
          break
        of "i", "a", "q", "e", "p":
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

  echo "Error: No command specified."
  dispatchcommand("help", flakePath, @[])

