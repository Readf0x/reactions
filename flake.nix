{
  description = "";
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  outputs = inputs @ {flake-parts, self, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux"];
      perSystem = { system, pkgs, lib, ... }: let
        dependencies = with pkgs; [
          cairo
          glib
          # gtk3
          gtk4
          pango
        ];
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            vala
            vala-language-server
            gnumake
            pkg-config
          ];
          buildInputs = dependencies;
        };
        packages = rec {
          reactions = pkgs.stdenv.mkDerivation (final: {
            pname = "reactions";
            version = "0.0.0";
            src = ./.;

            nativeBuildInputs = with pkgs; [ vala pkg-config ];
            buildInputs = dependencies;

            DESTDIR = placeholder "out";
            PREFIX = "";

            meta = {
              description = "";
              homepage = "https://git.gay/readf0x/reactions/";
              license = lib.licenses.gpl3Only;
            };
          });
          default = reactions;
        };
      };
    };
}
