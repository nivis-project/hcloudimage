---
# hcloudimage-lofb
title: 06 CI pipeline
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:32Z
updated_at: 2026-07-13T15:29:30Z
blocked_by:
    - hcloudimage-r4d6
---

ci.yml: golangci-lint, unit tests + coverage, terraform/tofu validate matrix over examples/, nix flake check (hermetic test), docs-diff check. (BRIEFING.md §13.6, §9)

## Summary of Changes

Milestone 06 complete. CI pipeline codifies the quality bar built in milestones 01-05: lint, unit+lifecycle tests with coverage, example validation under both terraform and tofu, the hermetic nix flake check, and a tfplugindocs docs-diff gate. Committed tfplugindocs output. Every CI step mirrors a locally-verified command. OpenSpec change ci-pipeline archived.
