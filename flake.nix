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
        packages.hei = pkgs.nimblePackages.${system}.hei;
        devShells.default = pkgs.mkShell {
          name = "hei";
          buildInputs = with pkgs; [
            nim2
            nim2Packages.safeseq
            nimble-unwrapped
          ];
        };
      }
    );
}
