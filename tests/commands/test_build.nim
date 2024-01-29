discard """
  output: '''which nom
nom build -L .'''
  exitcode: 0
  joinable: false
"""

import ../../src/hei/commands
import std/envvars
putEnv "HEI_TESTING", "1"
dispatchCommand "build", "/etc/nixos", @[]
