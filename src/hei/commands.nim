import std/[os, tables]

var dispatchTable = initTable[string, proc(args: seq[string])]()

dispatchTable["help"] = proc (args: seq[string]) =
  echo "help"

dispatchTable["check"] = proc(args: seq[string]) =
  let res = execShellCmd "nix flake check " & args[0]
  system.quit(res)

dispatchTable["show"] = proc(args: seq[string]) =
  let res = execShellCmd "nix flake show " & args[0]
  system.quit(res)

dispatchTable["gc"] = proc(args: seq[string]) =
  let res = execShellCmd "nix gc"
  system.quit(res)

proc dispatchCommand*(cmd: string, args: seq[string]) =
  if dispatchTable.hasKey(cmd):
    dispatchTable[cmd](args)
  else:
    echo("Unknown command: "&cmd)
