---
# hcloudimage-1uf7
title: Acceptance tests + cleanup workflow
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-r4gl
blocked_by:
    - hcloudimage-lofb
---

TF_ACC=1 tests composing hcloudimage_image + official hcloud_server; SSH-reachability assertion from runner with baked key; ForceNew + in-place semantics against real snapshots; both x86 and arm (arm toggle-gated); cost controls; acceptance.yml (dispatch/push-main/nightly, never fork PRs) + cleanup.yml orphan sweep. (BRIEFING.md §8.3, §9)
