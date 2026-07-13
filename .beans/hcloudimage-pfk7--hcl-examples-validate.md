---
# hcloudimage-pfk7
title: HCL examples + validate
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T15:09:18Z
parent: hcloudimage-1nw4
blocked_by:
    - hcloudimage-541n
---

examples/resources/hcloudimage_image/resource.tf (§5 verbatim behaviour), examples/provider/provider.tf, examples/data-sources/hcloudimage_snapshot/data-source.tf. terraform validate AND tofu validate pass. (BRIEFING.md §5)

## Summary of Changes

Added examples/provider, examples/resources/hcloudimage_image (BRIEFING §5 verbatim: image_path + filesha256, composes hcloud_server from the snapshot id), examples/data-sources/hcloudimage_snapshot. scripts/validate-examples.sh builds a filesystem mirror so both terraform and tofu resolve nivis-project/hcloudimage offline while hcloud installs from its registry; all three examples validate under both binaries. Added a validate-examples GNUmakefile target.
