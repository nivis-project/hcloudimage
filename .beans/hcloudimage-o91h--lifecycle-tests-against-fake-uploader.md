---
# hcloudimage-o91h
title: Lifecycle tests against fake uploader
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-541n
blocked_by:
    - hcloudimage-jzim
---

terraform-plugin-testing suite with fake injected: apply creates state with synthetic id; ForceNew attrs fire RequiresReplace; description/labels update in place; Read removes missing snapshot from state; destroy deletes. Coverage reported. (BRIEFING.md §8.1)
