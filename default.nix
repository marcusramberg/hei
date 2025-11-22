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
  nativeBuildInputs = with pkgs; [
    nom
    nvd
    installShellFiles
  ];
  pwd = ./.;
  src = ./.;
  modules = ./gomod2nix.toml;
  ldflags = [
    "-s -w"
    "-X main.Version=${version}"
  ];
  postInstall = ''
    installManPage ./hei.1.man
    installShellCompletion --cmd hei \
      --bash <($out/bin/hei completion bash) \
      --fish <($out/bin/hei completion fish ) \
      --zsh <($out/bin/hei completion zsh)
  '';
}
