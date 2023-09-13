{
  description = "hei - hey ported to nim";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, utils, ... }:

    utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system}.pkgs;
      in
      {
        packages.hei = pkgs.nim2Packages.buildNimPackage {
          name = "hei";
          version = 0.1;
          src = ./.;
          buildinputs = with pkgs; [
            nim2
            nimble-unwrapped
          ];
        };
        devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
            nim2
            nim2Packages.safeseq
            nimble-unwrapped
          ];
        };
      }
    );
}
