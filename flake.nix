rec {
  description = "Description for the project";

  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux"];
      perSystem = {
        system,
        pkgs,
        lib,
        ...
      }: let
        info = {
          projectName = "reactions";
        };
        libs = with pkgs; [
          pango
          cairo
          glib
          gtk3
        ];
      in
        (
          {
            projectName,
            moduleName ? projectName,
          }: rec {
            devShells.default = pkgs.mkShell {
              packages = with pkgs; [
                go
                delve
                pkg-config
              ] ++ libs;

              LD_LIBRARY_PATH = lib.makeLibraryPath libs;
              PKG_CONFIG_PATH = lib.makeSearchPath "lib/pkgconfig" libs;
            };
            packages = {
              ${projectName} = pkgs.buildGoModule {
                pname = projectName;
                version = "0.1";

                src = ./.;

                vendorHash = null;

                meta = {
                  inherit description;
                  # homepage = "";
                  # license = lib.licenses.;
                  # maintainers = with lib.maintainers; [  ];
                };
              };
              default = packages.${projectName};
            };
          }
        )
        info;
    };
}
