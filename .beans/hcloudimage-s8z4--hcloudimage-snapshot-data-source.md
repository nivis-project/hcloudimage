---
# hcloudimage-s8z4
title: hcloudimage_snapshot data source
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:48:30Z
parent: hcloudimage-jzim
blocked_by:
    - hcloudimage-c2wx
---

Lookup by id or with_selector (+most_recent); exactly-one resolution semantics; computed fields. Unit tests. (BRIEFING.md §3.3)

## Summary of Changes

hcloudimage_snapshot data source per BRIEFING §3.3: id xor with_selector (ExactlyOneOf), most_recent, computed name/description/architecture/created/labels. Resolves via uploader.Find with ambiguity-unless-most_recent semantics. Unit tests cover schema, config validators, and the by-id/selector lookup via the fake.
