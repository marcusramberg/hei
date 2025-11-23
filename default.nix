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
  version = "0.2.0";
in
buildGoApplication {
  pname = "hei";
  inherit version;
  nativeBuildInputs = with pkgs; [
    installShellFiles
    makeWrapper
  ];
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  ldflags = [
    "-s -w"
    "-X main.Version=${version}"
  ];
  postInstall = ''
    installManPage ./hei.1.man --name hei.1
    installShellCompletion --cmd hei \
      --bash <($out/bin/hei completion bash) \
      --fish <($out/bin/hei completion fish ) \
      --zsh <($out/bin/hei completion zsh)
      wrapProgram $out/bin/hei \
          --prefix PATH : ${
            pkgs.lib.makeBinPath [
              pkgs.nom
              pkgs.nvd
            ]
          }
  '';
}
