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

      packages = forAllSystems (pkgs: {
        default = pkgs.buildGoModule {
          pname = "terraform-provider-hcloudimage";
          version = "0.1.0-dev";
          src = ./.;

          # Bump whenever go.mod / go.sum change:
          #   set to pkgs.lib.fakeHash, run `nix build .#default`, copy the "got:" hash.
          vendorHash = "sha256-HxlNF8o+jzQlr0Lfv5udQUSG2p7HqDkvTmWaM+32Txw=";

          # Stamp the version into the binary the way goreleaser does (BRIEFING.md §10).
          ldflags = [
            "-s"
            "-w"
            "-X main.version=0.1.0-dev"
          ];

          # The package build only compiles the provider. Tests run in the
          # devShell / CI (unit) and, for the protocol layer, the hermetic
          # NixOS-VM check (milestone 05) — they need a terraform binary and
          # TF_ACC, which the buildGoModule sandbox deliberately lacks.
          doCheck = false;

          meta = {
            description = "Terraform & OpenTofu provider: upload raw disk images to Hetzner Cloud as snapshots";
            homepage = "https://github.com/nivis-project/terraform-provider-hcloudimage";
            license = pkgs.lib.licenses.mpl20;
            mainProgram = "terraform-provider-hcloudimage";
          };
        };
      });

      # To be added per BRIEFING.md §7 as milestones progress:
      #   checks.hermetic-e2e       — NixOS-VM lifecycle test, Linux only (milestone 05)
      #   packages.test-image-x86   — Alpine amd64 fixture, .raw.xz, baked SSH key (milestone 07)
      #   packages.test-image-arm   — Alpine aarch64 fixture (milestone 07)
      #   packages.provider-mirror  — filesystem mirror for registry-less consumption (milestone 09)
    };
}
