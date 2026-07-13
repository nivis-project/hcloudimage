---
# hcloudimage-1nw4
title: 04 Real uploader + examples
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:07Z
updated_at: 2026-07-13T15:09:18Z
blocked_by:
    - hcloudimage-541n
---

Wire hcloudimages/v2 behind the Uploader interface; runnable HCL examples; terraform validate and tofu validate pass. (BRIEFING.md §13.4, §4.2, §5)

## Summary of Changes

Milestone 04 complete. Real hcloudimages/v2 uploader wired behind the Uploader interface with token-vs-fake selection precedence (HCLOUDIMAGE_FAKE env, token, fallback). Runnable HCL examples validate under both terraform and tofu via a filesystem mirror. go build/test green, golangci-lint clean, nix build .#default and nix flake check pass. OpenSpec change real-uploader-and-examples archived.
