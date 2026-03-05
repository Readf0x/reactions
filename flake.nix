rec {
  description = "Integral Prompt for zsh";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    flake-utils,
    nixpkgs,
    ...
  }:
    flake-utils.lib.eachSystem ["x86_64-linux"] (system: let
      pkgs = import nixpkgs { inherit system; };
      projectName = "reactions";
      libs = with pkgs; [
        pango
        cairo
        glib
        gtk3
      ];
    in {
      devShells.default = pkgs.mkShell {
        packages = with pkgs;
          [
            go
            delve
            pkg-config
          ]
          ++ libs;

        GSETTINGS_SCHEMA_DIR = "${pkgs.gtk3}/share/gsettings-schemas/${pkgs.gtk3.name}/glib-2.0/schemas";
      };
      packages = {
        ${projectName} = pkgs.buildGoModule {
          pname = projectName;
          version = "0.1";

          src = ./.;

          vendorHash = "sha256-jK87vZYfUe8znk65SmJ1mN8qP5K3dtt950hKGWTYXs4=";

          nativeBuildInputs = [pkgs.pkg-config];
          buildInputs = libs;

          meta = {
            inherit description;
            # homepage = "";
            # license = lib.licenses.;
            # maintainers = with lib.maintainers; [  ];
          };
        };
        default = self.packages.${system}.${projectName};
      };
    });
}
