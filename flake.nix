{
  description = "Cue terminal media client for Plex and Jellyfin";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      overlays.default = final: _prev: {
        cue = final.callPackage ./package.nix { };
      };

      packages = forAllSystems (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            overlays = [ self.overlays.default ];
          };
        in
        {
          inherit (pkgs) cue;
          default = pkgs.cue;
        }
      );

      apps = forAllSystems (system: {
        cue = {
          type = "app";
          program = "${self.packages.${system}.cue}/bin/cue";
          meta.description = self.packages.${system}.cue.meta.description;
        };
        default = self.apps.${system}.cue;
      });

      checks = forAllSystems (system: {
        cue = self.packages.${system}.cue;
      });
    };
}
