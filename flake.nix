{
  description = "sqlc-gen-queries - SQLC Queries Generator";

  nixConfig = {
    extra-substituters = [ "https://ralch.cachix.org" ];
    extra-trusted-public-keys = [ "ralch.cachix.org-1:ti5RSqCJwd+VIv7dlxJEAVD98KIJDLxj1yJY+S80Q3A=" ];
  };

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        version = (pkgs.lib.importJSON ./.github/config/release-please-manifest.json).".";
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "sqlc-gen-queries";
          inherit version;
          src = pkgs.lib.cleanSource ./.;
          subPackages = [ "cmd/sqlc-gen-queries" ];
          vendorHash = "sha256-t8yFhyd0oMLeK/AjIRzKaZXXjG0zdpJQTlPye/AcKj4=";
          doInstallCheck = true;
          installCheckPhase = ''
            $out/bin/sqlc-gen-queries --help
          '';
          meta = with pkgs.lib; {
            description = "SQLC Queries Generator";
            license = licenses.mit;
            mainProgram = "sqlc-gen-queries";
          };
        };

        devShells.default = pkgs.mkShell {
          name = "sqlc-gen-queries";
          packages = [
            pkgs.go
          ];
        };
      }
    );
}
