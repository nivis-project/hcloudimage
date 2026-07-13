{
  description = "Terraform & OpenTofu provider for uploading raw disk images to Hetzner Cloud as reusable snapshots";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      # Plain-nix multi-system support — deliberately no flake-utils.
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems =
        f:
        nixpkgs.lib.genAttrs systems (
          system:
          f (import nixpkgs {
            inherit system;
            config.allowUnfree = true; # terraform is BUSL-licensed
          })
        );
    in
    {
      devShells = forAllSystems (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [
            go
            golangci-lint
            terraform
            opentofu
            terraform-plugin-docs
            goreleaser
            gnumake
            hcloud-upload-image
          ];
        };
      });

      # To be added per BRIEFING.md §7 as milestones progress:
      #   packages.default          — provider via buildGoModule, pinned vendorHash (milestone 01)
      #   checks.hermetic-e2e       — NixOS-VM lifecycle test, Linux only (milestone 05)
      #   packages.test-image-x86   — Alpine amd64 fixture, .raw.xz, baked SSH key (milestone 07)
      #   packages.test-image-arm   — Alpine aarch64 fixture (milestone 07)
      #   packages.provider-mirror  — filesystem mirror for registry-less consumption (milestone 09)
    };
}
