discard """
  output: "nix-collect-garbage -d"
  exitcode: 0
  joinable: false
"""

import ../../src/hei/commands
import std/envvars
putEnv "HEI_TESTING", "1"
dispatchCommand "gc", "/etc/nixos", @[]
