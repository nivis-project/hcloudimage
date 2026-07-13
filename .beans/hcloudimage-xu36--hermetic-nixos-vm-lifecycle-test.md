---
# hcloudimage-xu36
title: Hermetic NixOS-VM lifecycle test
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T15:25:39Z
parent: hcloudimage-r4d6
blocked_by:
    - hcloudimage-1nw4
---

checks.hermetic-e2e via pkgs.testers.runNixOSTest: provider installed into in-VM filesystem mirror; init/plan/apply/destroy with fake uploader (HCLOUDIMAGE_FAKE=1 or build tag); asserts ForceNew vs in-place; runs under both terraform and tofu; gated by nix flake check. PoC DoD gate. (BRIEFING.md §8.2)

## Summary of Changes

Added checks.hermetic-e2e (Linux) via pkgs.testers.runNixOSTest: boots a NixOS VM with terraform + opentofu and the built provider installed through a filesystem mirror, then drives init/plan/apply/destroy under BOTH binaries with HCLOUDIMAGE_FAKE=1. Asserts create (id set), ForceNew replace on image_sha256 change (plan actions [delete,create], id changes), in-place update on labels change (plan action [update], id unchanged), and destroy (empty state). Wired into nix flake check as the PoC DoD gate. Refactored flake.nix so provider + mirror are shared helpers (packages.default, packages.provider-mirror, checks). Added nix/provider-mirror.nix (reused by milestone 09) and test/e2e/ (main.tf + hermetic.nix).

The VM test surfaced and fixed three real provider bugs: (1) labelValuesValidator decoded the labels map into map[string]string, which panics when an element value is unknown (labels = { env = var.x }) during the plan walk — now iterates elements and skips unknowns; (2) sha256RequiredWithPathValidator treated an unknown image_sha256 (from filesha256/var) as absent and wrongly errored — now defers when either attribute is unknown; (3) the fake uploader is per-process, so cross-invocation lifecycle was invisible — added optional file-backed persistence (HCLOUDIMAGE_FAKE_STATE) so the VM lifecycle is observable.
