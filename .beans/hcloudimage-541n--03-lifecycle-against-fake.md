---
# hcloudimage-541n
title: 03 Lifecycle against fake
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:07Z
updated_at: 2026-07-13T14:58:55Z
blocked_by:
    - hcloudimage-jzim
---

plan/apply/destroy + ForceNew/in-place behaviour under terraform-plugin-testing with the fake uploader injected. Gate: lifecycle unit tests green. (BRIEFING.md §13.3, §8.1)

## Summary of Changes

Milestone 03 complete. Full plan/apply/update/replace/destroy lifecycle proven against the fake through terraform-plugin-testing. Fixed three provider bugs found by the tests (Computed+ForceNew UseStateForUnknown, effective_labels plan modifier, null-description preservation). go test (unit + lifecycle) green at ~74% coverage, golangci-lint clean, nix build .#default and nix flake check pass. OpenSpec change lifecycle-tests-against-fake archived.
