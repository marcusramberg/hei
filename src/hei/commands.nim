import std/[os, strutils, tables, tempfiles]

var dispatchTable = initTable[string, proc(flakePath: string, args: seq[string])]()

dispatchTable["help"] = proc (flakePath: string, args: seq[string]) =
  echo "help"

dispatchTable["build"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix build -L " & args.join(" ")
  system.quit(res)

dispatchTable["check"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix flake check " & flakePath
  system.quit(res)

dispatchTable["show"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix flake show " & flakePath
  system.quit(res)

dispatchTable["gc"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix gc"
  system.quit(res)

dispatchTable["repl"] = proc(flakePath: string, args: seq[string]) =
  let (tmpfile, path) = createTempFile("dotfiles-repl.nix", "_end.tmp")
  tmpfile.write("import " & flakePath & "(builtins.getFlake \"" & flakePath & "\")")
  let res = execShellCmd "nix repl \\<nixpkgs\\> " & path
  system.quit(res)

dispatchTable["update"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "nix flake update " & flakePath
  system.quit(res)

dispatchTable["rebuild"] = proc(flakePath: string, args: seq[string]) =
  if hostOs == "macosx":
    let res = execShellCmd "darwin-rebuild switch --flake " & flakePath
    system.quit(res)
  else:
    let res = execShellCmd "sudo nixos-rebuild switch --flake " & flakePath
    system.quit(res)

dispatchTable["ssh"] = proc(flakePath: string, args: seq[string]) =
  let res = execShellCmd "ssh " & args[0] & " hei " & args[1..args.high].join(" ")
  system.quit(res)

dispatchTable["test"] = proc(flakePath: string, args: seq[string]) =
  dispatchTable["rebuild"](flakePath, @["--fast"])

dispatchTable["upgrade"] = proc(flakePath: string, args: seq[string]) =
  discard execShellCmd "nix flake update " & flakePath
  dispatchTable["rebuild"](flakePath, args)

proc dispatchCommand*(cmd: string, flakePath: string, args: seq[string]) =
  if dispatchTable.hasKey(cmd):
    dispatchTable[cmd](flakePath, args)
  else:
    echo("Unknown command: "&cmd)
