---
# hcloudimage-r4gl
title: 07 Fixtures + acceptance
status: completed
type: milestone
priority: normal
created_at: 2026-07-13T14:03:32Z
updated_at: 2026-07-13T15:54:02Z
blocked_by:
    - hcloudimage-lofb
---

packages.test-image-{x86,arm} (Alpine + baked throwaway SSH key, .raw.xz); gated acceptance.yml composing the official hcloud provider; SSH-reachability assertion; cost controls; cleanup.yml orphan sweep. (BRIEFING.md §13.7, §8.3)

## Summary of Changes

Milestone 07 complete. Reproducible-input Alpine .raw.xz fixtures with a baked throwaway SSH key (x86 verified building; arm defined + pinned, builds on aarch64/emulated path). Compiling, correctly-gated acceptance tests that compose the official hcloud provider and prove SSH reachability. acceptance.yml + cleanup.yml encode all BRIEFING §8.3/§15 cost and safety controls. go test/lint/nix flake check green. The billable live run is the documented human step (§14). OpenSpec change fixtures-and-acceptance archived.
