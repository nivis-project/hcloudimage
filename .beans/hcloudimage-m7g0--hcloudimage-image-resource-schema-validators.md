---
# hcloudimage-m7g0
title: hcloudimage_image resource schema & validators
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:48:29Z
parent: hcloudimage-jzim
blocked_by:
    - hcloudimage-c2wx
---

Full §3.2 schema incl. timeouts block; config validators: image_url/image_path mutual exclusion, image_sha256 required iff image_path, label rules; plan modifiers for ForceNew set. Unit tests. (BRIEFING.md §3.2)

## Summary of Changes

Full hcloudimage_image schema per BRIEFING §3.2: all attributes with MarkdownDescription, computed id/effective_labels, timeouts block. RequiresReplace plan modifiers on all ForceNew attributes; description/labels have none (in-place). Config validators: ExactlyOneOf(image_url,image_path), image_sha256-iff-image_path, enum validation (architecture/compression/format), label-value / rule. CRUD wired against the Uploader interface. Unit tests cover schema shape, ForceNew-vs-in-place plan modifiers, config to UploadRequest mapping across all compressions/both sources/both arches, and the label validator. Full plan/apply/destroy lifecycle assertions are milestone 03 (plugin-testing framework).
