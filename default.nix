{
  pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  ),
  buildGoApplication ? pkgs.buildGoApplication,
}:

let
  version = "0.1";
in
buildGoApplication {
  pname = "hei";
  inherit version;
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  ldflags = [
    "-s -w"
    "-X main.Version=${version}"
  ];
}
