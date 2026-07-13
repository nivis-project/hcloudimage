---
# hcloudimage-2dsy
title: Real uploader via hcloudimages/v2
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T15:09:18Z
parent: hcloudimage-1nw4
blocked_by:
    - hcloudimage-541n
---

uploader_hcloud.go wrapping hcloudimages.Client + hcloud.Client; config→UploadRequest mapping (all compressions, both sources, arch mapping); robust cleanup, no public skip-cleanup knob. (BRIEFING.md §4.2, Appendix)

## Summary of Changes

Added uploader_hcloud.go: a real Uploader backed by hcloudimages/v2 (pinned v2.0.1) + hcloud-go v2. Upload maps UploadRequest to UploadOptions/WriteOptions (image_url as ImageURL, image_path opened as ImageReader; all compressions, raw/qcow2, arch to hcloud.ArchitectureX86/ARM; optional server_type/location so library defaults apply otherwise; labels/description). Delete/Get/UpdateMetadata/Find via the hcloud Image API; Get returns nil when absent. No public skip-cleanup: gated behind HCLOUDIMAGE_DEBUG_SKIP_CLEANUP env only. Unit tests cover the compression/format/architecture mapping and toSnapshotInfo.
