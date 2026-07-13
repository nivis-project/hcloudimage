---
# hcloudimage-xu36
title: Hermetic NixOS-VM lifecycle test
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-r4d6
blocked_by:
    - hcloudimage-1nw4
---

checks.hermetic-e2e via pkgs.testers.runNixOSTest: provider installed into in-VM filesystem mirror; init/plan/apply/destroy with fake uploader (HCLOUDIMAGE_FAKE=1 or build tag); asserts ForceNew vs in-place; runs under both terraform and tofu; gated by nix flake check. PoC DoD gate. (BRIEFING.md §8.2)
