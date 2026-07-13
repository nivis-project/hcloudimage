---
# hcloudimage-val8
title: Uploader interface + fake implementation
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:48:30Z
parent: hcloudimage-jzim
blocked_by:
    - hcloudimage-c2wx
---

internal/provider/uploader.go interface (Upload/Delete/Get); uploader_fake.go in-memory fake recording calls, synthetic IDs, out-of-band-deletion simulation. Resource depends only on the interface. (BRIEFING.md §4.1)

## Summary of Changes

Uploader interface (Upload/Delete/Get/UpdateMetadata/Find) with library-agnostic UploadRequest/SnapshotInfo types and a shared apricote.de/created-by label-merge helper. In-memory FakeUploader: synthetic incrementing IDs, records calls, merges created-by label, simulates out-of-band deletion. Resource and data source depend only on the interface; provider injects the fake (real impl in milestone 04). Unit tests cover create/get/delete, out-of-band delete, update, and find-by-id/selector.
