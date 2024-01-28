import std / [parseopt, os, envvars]

import strformat, strutils, sequtils
import hei/[commands, utils]

const version = staticExec("git describe --tags HEAD")

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
      of cmdEnd: break
      of cmdShortOption, cmdLongOption:
        case p.key
        of "f", "flake":
          flakePath = p.val
        of "d", "dry-run":
          putEnv("NIX_DEBUG", "1")
        of "D", "debug":
          putEnv("NIX_DEBUG", "1")
          putEnv("NIX_SHOW_TRACE", "1")
        of "v", "version":
          echo fmt"hei {version} - running on {hostOs}({hostCpu})"
          quit()
        of "h", "help":
          break
        of "i", "a", "q", "e", "p":
          # run nix-env with the original command line arguments
          echo "forwarding to nix-env ..."
          let res = execShellCmd "nix-env " & commandLineParams().join(" ")
          system.quit(res)
        else:
          echo "Unknown option: ", p.key, ". run `hei` for help."
          quit()
      of cmdArgument:
        dispatchCommand(p.key, flakePath, p.remainingArgs)
        quit()

  echo "Error: No command specified."
  dispatchCommand("help", flakePath, @[])
