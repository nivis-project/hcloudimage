---
# hcloudimage-a5ua
title: Alpine test image fixtures (x86 + arm)
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T15:54:02Z
parent: hcloudimage-r4gl
blocked_by:
    - hcloudimage-lofb
---

packages.test-image-{x86,arm}: reproducible minimal Alpine generic-cloud raw images (amd64/aarch64), throwaway SSH pubkey baked into /root/.ssh/authorized_keys, sshd+DHCP enabled, recompressed .raw.xz. Document aarch64 build path. (BRIEFING.md §7, §8.3)

## Summary of Changes

Added nix/test-image.nix: reproducible-input Alpine fixture derivation. Fetches a pinned Alpine nocloud raw image (amd64/aarch64) via fixed-output fetchurl, converts to raw, and uses guestfs-tools virt-customize (driven by libguestfs-with-appliance, offline) to bake the throwaway SSH pubkey into /root/.ssh/authorized_keys, enable sshd + networking, and set DHCP; asserts the guest arch via virt-inspector (no silent x86 fallback); recompresses to .raw.xz. flake.nix exposes packages.test-image-x86 and test-image-arm. Verified: test-image-x86 builds to a valid 93MB .raw.xz. Added test/fixtures/ throwaway ed25519 keypair (documented as disposable). Disk images are hermetic + input-pinned rather than byte-reproducible (fs metadata), documented as such.
