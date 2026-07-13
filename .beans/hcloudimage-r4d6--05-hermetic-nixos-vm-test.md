---
# hcloudimage-r4d6
title: 05 Hermetic NixOS-VM test
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:07Z
updated_at: 2026-07-13T15:25:39Z
blocked_by:
    - hcloudimage-1nw4
---

checks.hermetic-e2e wired into nix flake check: mirror install, init/plan/apply/destroy with fake uploader, under both terraform and tofu. This is the PoC Definition-of-Done gate. (BRIEFING.md §13.5, §8.2)

## Summary of Changes

Milestone 05 complete — this is the PoC Definition-of-Done gate (BRIEFING §12). The hermetic NixOS-VM lifecycle test runs under nix flake check, proving the full Terraform/OpenTofu protocol path (create/replace/in-place/destroy) with the fake uploader, hermetically and cloud-free, under both terraform and tofu. Fixed three provider bugs found by the VM test. nix flake check passes end to end. OpenSpec change hermetic-nixos-vm-test archived. The alpha base is now demonstrably working.
