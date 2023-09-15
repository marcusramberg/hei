import std/[os, strutils, tables, tempfiles]


var dispatchTable = initTable[string, proc(args: seq[string])]()

dispatchTable["help"] = proc (args: seq[string]) =
  echo "help"

dispatchTable["build"] = proc(args: seq[string]) =
  echo $args
  let res = execShellCmd "nix build -L " & args[1]
  system.quit(res)
dispatchTable["check"] = proc(args: seq[string]) =
  let res = execShellCmd "nix flake check " & args[0]
  system.quit(res)

dispatchTable["show"] = proc(args: seq[string]) =
  let res = execShellCmd "nix flake show " & args.join(" ")
  system.quit(res)

dispatchTable["gc"] = proc(args: seq[string]) =
  let res = execShellCmd "nix gc"
  system.quit(res)

dispatchTable["repl"] = proc(args: seq[string]) =
  let (tmpfile, path) = createTempFile("dotfiles-repl.nix", "_end.tmp")

  tmpfile.write("import " & args[0] & "(builtins.getFlake \"" & args[0] & "\")")
  let res = execShellCmd "nix repl \\<nixpkgs\\> " & path
  system.quit(res)

dispatchTable["update"] = proc(args: seq[string]) =
  let res = execShellCmd "nix flake update " & args[0]
  system.quit(res)

dispatchTable["rebuild"] = proc(args: seq[string]) =
  if hostOs == "macosx":
    let res = execShellCmd "darwin-rebuild switch --flake " & args[0]
    system.quit(res)
  else:
    let res = execShellCmd "sudo nixos-rebuild switch --flake " & args.join(" ")
    system.quit(res)

dispatchTable["ssh"] = proc(args: seq[string]) =
  let res = execShellCmd "ssh " & args[1] & " hei " & args[2..args.high].join(" ")
  system.quit(res)

dispatchTable["test"] = proc(args: seq[string]) =
  dispatchTable["rebuild"](@[args[0], "--fast"])

dispatchTable["upgrade"] = proc(args: seq[string]) =
  discard execShellCmd "nix flake update " & args[0]
  dispatchTable["rebuild"](args)

proc dispatchCommand*(cmd: string, args: seq[string]) =
  if dispatchTable.hasKey(cmd):
    dispatchTable[cmd](args)
  else:
    echo("Unknown command: "&cmd)
