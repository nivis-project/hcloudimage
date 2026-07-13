---
# hcloudimage-val8
title: Uploader interface + fake implementation
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-jzim
blocked_by:
    - hcloudimage-c2wx
---

internal/provider/uploader.go interface (Upload/Delete/Get); uploader_fake.go in-memory fake recording calls, synthetic IDs, out-of-band-deletion simulation. Resource depends only on the interface. (BRIEFING.md §4.1)
