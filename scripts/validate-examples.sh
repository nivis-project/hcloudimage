#!/usr/bin/env bash
# Validate every example under examples/ with both terraform and tofu.
#
# Installs the locally built provider through a filesystem mirror (so both
# terraform and tofu resolve nivis-project/hcloudimage without a public
# registry), while other providers (hetznercloud/hcloud) install normally from
# their registry. This mirrors how Nivis consumes the provider (milestone 09)
# and validates the composed example under both binaries.
#
# Requires: terraform, tofu, and a built provider binary. Pass the directory
# holding terraform-provider-hcloudimage as $1 (defaults to ./result/bin).
set -euo pipefail

PROVIDER_BIN_DIR="${1:-$(pwd)/result/bin}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VERSION="0.1.0"

PROVIDER_BIN="${PROVIDER_BIN_DIR}/terraform-provider-hcloudimage"
if [[ ! -x "${PROVIDER_BIN}" ]]; then
  echo "provider binary not found at ${PROVIDER_BIN}; run 'nix build .#default' first" >&2
  exit 1
fi

os="$(go env GOOS 2>/dev/null || uname -s | tr '[:upper:]' '[:lower:]')"
arch="$(go env GOARCH 2>/dev/null || echo amd64)"

MIRROR="$(mktemp -d)"
RC_FILE="$(mktemp)"
trap 'rm -rf "${MIRROR}" "${RC_FILE}"' EXIT

for host in registry.terraform.io registry.opentofu.org; do
  dest="${MIRROR}/${host}/nivis-project/hcloudimage/${VERSION}/${os}_${arch}"
  mkdir -p "${dest}"
  cp "${PROVIDER_BIN}" "${dest}/terraform-provider-hcloudimage_v${VERSION}"
done

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

status=0
for dir in \
  "${REPO_ROOT}/examples/provider" \
  "${REPO_ROOT}/examples/resources/hcloudimage_image" \
  "${REPO_ROOT}/examples/data-sources/hcloudimage_snapshot"; do

  echo "== validating ${dir#"${REPO_ROOT}/"} =="

  for bin in terraform tofu; do
    if ! command -v "${bin}" >/dev/null 2>&1; then
      echo "  ${bin}: not found, skipping" >&2
      continue
    fi
    ( cd "${dir}" && rm -rf .terraform .terraform.lock.hcl )
    if ( cd "${dir}" \
      && TF_CLI_CONFIG_FILE="${RC_FILE}" "${bin}" init -backend=false -input=false >/dev/null 2>&1 \
      && TF_CLI_CONFIG_FILE="${RC_FILE}" "${bin}" validate >/dev/null ); then
      echo "  ${bin}: OK"
    else
      echo "  ${bin}: FAILED"
      ( cd "${dir}" && TF_CLI_CONFIG_FILE="${RC_FILE}" "${bin}" validate || true )
      status=1
    fi
  done
done

exit "${status}"
