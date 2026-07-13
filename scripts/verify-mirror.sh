#!/usr/bin/env bash
# Prove the provider is consumable from the Nix filesystem mirror — the way Nivis
# consumes it — WITHOUT any public registry (BRIEFING.md §7). Builds
# packages.provider-mirror, points a CLI config at it, and runs a pinned-version
# consumer config through `init` + `plan` under both terraform and tofu.
#
# Uses HCLOUDIMAGE_FAKE=1 so no Hetzner token is needed; plan is enough to prove
# the provider was resolved and its schema loaded from the mirror.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="0.1.0"

echo "== building packages.provider-mirror =="
MIRROR="$(nix build "${REPO_ROOT}#provider-mirror" --no-link --print-out-paths | tail -1)"
echo "mirror: ${MIRROR}"

RC_FILE="$(mktemp)"
WORK="$(mktemp -d)"
trap 'rm -rf "${RC_FILE}" "${WORK}"' EXIT

cat >"${RC_FILE}" <<EOF
provider_installation {
  filesystem_mirror {
    path    = "${MIRROR}"
    include = ["registry.terraform.io/nivis-project/hcloudimage", "registry.opentofu.org/nivis-project/hcloudimage"]
  }
  direct {
    exclude = ["registry.terraform.io/nivis-project/hcloudimage", "registry.opentofu.org/nivis-project/hcloudimage"]
  }
}
EOF

# A real consumer config with a PINNED version — the Nivis-shaped usage.
cat >"${WORK}/main.tf" <<EOF
terraform {
  required_providers {
    hcloudimage = {
      source  = "nivis-project/hcloudimage"
      version = "${VERSION}"
    }
  }
}

provider "hcloudimage" {}

resource "hcloudimage_image" "consumed" {
  image_url    = "https://example.com/image.raw.xz"
  architecture = "x86"
  compression  = "xz"
  labels       = { consumed_from = "nix-mirror" }
}
EOF

status=0
for bin in terraform tofu; do
  if ! command -v "${bin}" >/dev/null 2>&1; then
    echo "  ${bin}: not found, skipping" >&2
    continue
  fi
  echo "== ${bin}: init + plan from the mirror =="
  ( cd "${WORK}" && rm -rf .terraform .terraform.lock.hcl )
  if ( cd "${WORK}" \
        && TF_CLI_CONFIG_FILE="${RC_FILE}" HCLOUDIMAGE_FAKE=1 "${bin}" init -backend=false -input=false >/dev/null \
        && TF_CLI_CONFIG_FILE="${RC_FILE}" HCLOUDIMAGE_FAKE=1 "${bin}" plan -input=false >/dev/null ); then
    echo "  ${bin}: OK — provider resolved from the mirror, plan succeeded"
  else
    echo "  ${bin}: FAILED"
    ( cd "${WORK}" && TF_CLI_CONFIG_FILE="${RC_FILE}" HCLOUDIMAGE_FAKE=1 "${bin}" plan -input=false || true )
    status=1
  fi
done

exit "${status}"
