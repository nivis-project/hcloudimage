---
# hcloudimage-q9c2
title: goreleaser + signing + manifest
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-e1z3
blocked_by:
    - hcloudimage-r4gl
---

.goreleaser.yml (standard provider os/arch matrix, SHA256SUMS + GPG .sig, manifest in archives, conventional-commit changelog); terraform-registry-manifest.json; release.yml on v* tags. Secrets are human-provided (§14) — stub and document. (BRIEFING.md §10)
