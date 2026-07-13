---
# hcloudimage-yw1z
title: Registry publication + Nivis mirror consumption
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:23Z
updated_at: 2026-07-13T16:02:50Z
parent: hcloudimage-z0sk
blocked_by:
    - hcloudimage-e1z3
---

packages.provider-mirror filesystem layout (registry.terraform.io/nivis-project/hcloudimage/<version>/<os>_<arch>/); document dev_overrides alternative; document Terraform + OpenTofu registry onboarding (human steps §14); Nivis end-to-end consumption doc. (BRIEFING.md §7, §10, §14)

## Summary of Changes

Added scripts/verify-mirror.sh + make consume-mirror: builds packages.provider-mirror, points a filesystem_mirror CLI config at it, and runs a pinned-version (0.1.0) consumer through init + plan under BOTH terraform and tofu with HCLOUDIMAGE_FAKE=1 — proving registry-less consumption end to end (verified: both binaries resolve the provider from the mirror and plan succeeds, no registry, no token). Added docs/consuming-from-nix-mirror.md (filesystem_mirror recipe leading, dev_overrides alternative, tofu/terraform init difference noted) and docs/publishing.md (Terraform + OpenTofu registry onboarding as human §14 steps + prerequisites). The tfplugindocs docs-diff check stays clean (hand-written docs are not tfplugindocs-managed).
