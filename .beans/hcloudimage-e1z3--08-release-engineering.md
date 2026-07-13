---
# hcloudimage-e1z3
title: 08 Release engineering
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:32Z
updated_at: 2026-07-13T15:59:14Z
blocked_by:
    - hcloudimage-r4gl
---

goreleaser + GPG signing + terraform-registry-manifest.json + tfplugindocs docs; CHANGELOG; cut v0.1.0. (BRIEFING.md §13.8, §10, §11)

## Summary of Changes

Milestone 08 complete. goreleaser config produces signed, registry-shaped artifacts (validated by goreleaser check + a full snapshot build); terraform-registry-manifest.json declares protocol 6.0; release.yml cuts the signed release on v* tags; CHANGELOG documents 0.1.0. nix flake check green. The signed v0.1.0 release and registry onboarding are the documented human steps (§14). OpenSpec change release-engineering archived.
