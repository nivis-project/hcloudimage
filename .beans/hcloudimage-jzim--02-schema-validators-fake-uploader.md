---
# hcloudimage-jzim
title: 02 Schema, validators, fake uploader
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:07Z
updated_at: 2026-07-13T14:48:30Z
blocked_by:
    - hcloudimage-c2wx
---

Full resource/data-source schema, config validators, Uploader interface + in-memory fake, unit tests green. Gate: go test ./... green with high coverage of schema/validators. (BRIEFING.md §13.2, §3, §4.1)

## Summary of Changes

Milestone 02 complete. Full provider surface built and driven by the fake uploader: hcloudimage_image schema+validators+plan-modifiers, hcloudimage_snapshot data source, and the Uploader seam. go build, go test (unit), nix build .#default and nix flake check all green. OpenSpec change schema-validators-fake-uploader archived.
