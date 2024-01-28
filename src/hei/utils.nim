# Take a list of folders, expand tilde and return the first that has a flake.nix
import os
import strutils


proc hasFlake(flake: string): bool =
  let strippedFlake = flake.strip(leading = false, chars = {'/'})
  return fileExists(strippedFlake & "/flake.nix")

proc findFlake*(folders: seq[string]): string =
  for folder in folders:
    let expanded = folder.expandTilde
    if expanded.hasFlake:
      return folder
