# Builds a Terraform/OpenTofu filesystem-mirror layout for the provider, so it can
# be consumed without a public registry (BRIEFING.md §7). Shared by the hermetic
# test (milestone 05) and packages.provider-mirror (milestone 09).
#
#   <out>/registry.terraform.io/nivis-project/hcloudimage/<version>/<os>_<arch>/terraform-provider-hcloudimage_v<version>
#   <out>/registry.opentofu.org/...   (same, so tofu resolves it too)
{
  pkgs,
  provider, # packages.default (buildGoModule result)
  version ? "0.1.0",
}:
let
  inherit (pkgs.stdenv.hostPlatform) system;
  # Map Nix system to Terraform's <os>_<arch>.
  platform =
    {
      "x86_64-linux" = "linux_amd64";
      "aarch64-linux" = "linux_arm64";
      "x86_64-darwin" = "darwin_amd64";
      "aarch64-darwin" = "darwin_arm64";
    }
    .${system};
in
pkgs.runCommand "terraform-provider-hcloudimage-mirror-${version}" { } ''
  for host in registry.terraform.io registry.opentofu.org; do
    dir="$out/$host/nivis-project/hcloudimage/${version}/${platform}"
    mkdir -p "$dir"
    cp ${provider}/bin/terraform-provider-hcloudimage "$dir/terraform-provider-hcloudimage_v${version}"
  done
''
