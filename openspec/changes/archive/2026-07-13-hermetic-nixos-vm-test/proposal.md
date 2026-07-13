## Why

BRIEFING.md §12 names the hermetic NixOS-VM lifecycle test as the **Definition-of-Done
gate** for the PoC, and §8.2 specifies it. It proves the full Terraform ↔ provider protocol
path — `init → plan → apply → destroy`, ForceNew replace, in-place update — hermetically,
under both `terraform` and `tofu`, at zero cloud cost, using the fake uploader
(`HCLOUDIMAGE_FAKE=1`). Wiring it into `nix flake check` makes the gate enforceable in CI
and locally with one command. This is the milestone that makes the alpha base real.

## What Changes

- Add `checks.hermetic-e2e` to `flake.nix` (Linux systems only) via
  `pkgs.testers.runNixOSTest`.
- The VM boots with `terraform`, `opentofu`, and a filesystem mirror containing the
  Nix-built provider installed under `registry.terraform.io/nivis-project/hcloudimage` and
  `registry.opentofu.org/...`, plus a CLI config selecting the mirror.
- A test Terraform config (in the store) exercised in the VM: create, change
  `image_sha256` (expect replace), change `labels` (expect in-place), destroy.
- The Python test script runs the whole sequence under **both** `terraform` and `tofu`
  with `HCLOUDIMAGE_FAKE=1`, asserting the plan actions and final empty state.
- `nix flake check` runs it as part of the gate.

## Capabilities

### New Capabilities
- `hermetic-e2e`: the in-VM, cloud-free protocol lifecycle test gating `nix flake check`.

### Modified Capabilities
- `provider-scaffold`: `flake.nix` now exposes `checks.hermetic-e2e` on Linux.

## Impact

- New: `test/e2e/` (the NixOS test module + test HCL + provider-config helper), referenced
  from `flake.nix`.
- Modified: `flake.nix` (`checks` output, Linux-gated).
- Gate: `nix flake check` builds and runs the VM test; it is the PoC DoD gate.
