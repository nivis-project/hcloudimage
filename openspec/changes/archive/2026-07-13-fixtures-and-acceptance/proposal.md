## Why

Milestone 07 (BRIEFING.md §13.7, §8.3) adds the billable, real-Hetzner acceptance layer:
reproducible Alpine test-image fixtures (`packages.test-image-{x86,arm}`), acceptance
tests that compose the official `hcloud` provider and assert real guest reachability over
SSH, and the CI workflows (`acceptance.yml`, `cleanup.yml`) with mandatory cost/safety
controls. The acceptance *run* needs a real `HCLOUD_TOKEN` and a throwaway SSH key — those
are human-provided secrets (BRIEFING.md §14) — so the code is complete and gated; the live
run is a documented human step.

## What Changes

- `packages.test-image-x86` / `packages.test-image-arm`: reproducible Nix derivations that
  take a pinned Alpine generic-cloud raw image (amd64/aarch64), bake a throwaway SSH public
  key into `/root/.ssh/authorized_keys`, enable `sshd` + DHCP, and recompress to `.raw.xz`.
- Acceptance tests (`TF_ACC=1`, real uploader) that compose `hcloudimage_image` +
  `hcloud_server`, boot from the snapshot id, and assert reachability by **SSHing from the
  runner into the server** with the baked key and reading `/etc/os-release` — not merely
  `running`. Cover both `x86` (cx22) and `arm` (cax11). Deferred cleanup runs even on
  failure.
- `.github/workflows/acceptance.yml`: `workflow_dispatch` + push-to-main + nightly
  schedule; never on fork PRs; concurrency-limited; arm gated behind a dispatch input;
  always-run cleanup; uses `HCLOUD_TOKEN` secret. Composes both providers.
- `.github/workflows/cleanup.yml`: nightly label-scoped orphan sweep via
  `hcloud-upload-image cleanup`.
- Document the aarch64 build path (native runner / remote builder / binfmt+QEMU) and which
  the project uses.

## Capabilities

### New Capabilities
- `test-image-fixtures`: reproducible Alpine `.raw.xz` fixtures with a baked SSH key.
- `acceptance-tests`: real-Hetzner acceptance suite + gated `acceptance.yml`/`cleanup.yml`.

### Modified Capabilities
- `provider-scaffold`: flake exposes `packages.test-image-{x86,arm}`.

## Impact

- New: `nix/test-image.nix`, `test/fixtures/**` (SSH key gen + Alpine overlay),
  `internal/provider/*_acc_test.go`, `.github/workflows/acceptance.yml`,
  `.github/workflows/cleanup.yml`, docs on the aarch64 path and running acceptance.
- Modified: `flake.nix` (test-image packages), README (acceptance + cost controls).
- Gate: fixtures build reproducibly to `.raw.xz`; acceptance tests **compile** and are
  correctly gated (real run is a documented human step needing `HCLOUD_TOKEN`); workflows
  are valid and encode all §8.3 cost/safety controls.
