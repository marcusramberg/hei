{
  description = "hei - hey ported to nim";
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { nixpkgs, utils, self, ... }:

    utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system}.pkgs;
      in
      {
        packages.default = pkgs.nim2Packages.buildNimPackage {
          name = "hei";
          version = 0.1;
          nimBinOnly = true;
          src = ./.;
          nativeBuildInputs = with pkgs; [
            (writeScriptBin "git" ''
              echo ${ self.ref or "dirty" }
            '')
            installShellFiles
          ];
          buildinputs = with pkgs; [
            nim
            nix-output-monitor
          ];
          postInstall = ''
            export NIX_SYSTEM_FLAKE="."
            installShellCompletion --cmd hei \
              --bash <($out/bin/hei completions bash) \
              --fish <($out/bin/hei completions fish) \
              --zsh <($out/bin/hei completions zsh)
          '';

        };
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            nim
            nimlsp
            nim2Packages.safeseq
            nix-output-monitor
          ];
        };
      }
    );
}
