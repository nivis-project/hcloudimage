---
# hcloudimage-2dsy
title: Real uploader via hcloudimages/v2
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-1nw4
blocked_by:
    - hcloudimage-541n
---

uploader_hcloud.go wrapping hcloudimages.Client + hcloud.Client; config→UploadRequest mapping (all compressions, both sources, arch mapping); robust cleanup, no public skip-cleanup knob. (BRIEFING.md §4.2, Appendix)
