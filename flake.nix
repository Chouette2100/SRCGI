# <project-root>/flake.nix
{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            delve
            golangci-lint
            # Web開発に必要なものがあれば追加
            sqlx-cli 
          ];

          shellHook = ''
            echo "Go Development Environment Loaded!"
            go version
          '';
        };
      });
}
