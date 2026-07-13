---
# hcloudimage-a5ua
title: Alpine test image fixtures (x86 + arm)
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T14:11:44Z
parent: hcloudimage-r4gl
blocked_by:
    - hcloudimage-lofb
---

packages.test-image-{x86,arm}: reproducible minimal Alpine generic-cloud raw images (amd64/aarch64), throwaway SSH pubkey baked into /root/.ssh/authorized_keys, sshd+DHCP enabled, recompressed .raw.xz. Document aarch64 build path. (BRIEFING.md §7, §8.3)
