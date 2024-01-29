discard """
  output: '''sudo nixos-rebuild switch --flake /etc/nixos'''
  exitcode: 0
  joinable: false
"""

import ../../src/hei/commands
import std/envvars
putEnv "HEI_TESTING", "1"
dispatchCommand "rebuild", "/etc/nixos", @[]
