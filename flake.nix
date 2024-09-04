{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };

  outputs = inputs:
    with inputs; let
      forEachSupportedSystem = let
        supportedSystems = [
          "x86_64-linux"
          "aarch64-linux"
          "x86_64-darwin"
          "aarch64-darwin"
        ];
      in (f:
        nixpkgs.lib.genAttrs supportedSystems (system:
          f {
            pkgs = let
              overlays = [];
            in
              import nixpkgs {inherit overlays system;};
          }));
    in {
      devShells = forEachSupportedSystem ({pkgs}:
        with pkgs; {
          default = mkShell {
            packages = [
              go

              goose
              go-jet
              just
              hurl
            ];
          };
        });
    };
}
