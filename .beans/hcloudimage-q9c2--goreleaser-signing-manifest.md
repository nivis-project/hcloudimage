---
# hcloudimage-q9c2
title: goreleaser + signing + manifest
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T18:35:02Z
parent: hcloudimage-e1z3
blocked_by:
    - hcloudimage-r4gl
---

.goreleaser.yml (standard provider os/arch matrix, SHA256SUMS + GPG .sig, manifest in archives, conventional-commit changelog); terraform-registry-manifest.json; release.yml on v* tags. Secrets are human-provided (§14) — stub and document. (BRIEFING.md §10)

## Summary of Changes

Added .goreleaser.yml (v2): builds linux/darwin/windows/freebsd x amd64/arm64/arm/386 (standard provider ignores), zip archives named {ProjectName}_{Version}_{Os}_{Arch} with terraform-registry-manifest.json embedded, {ProjectName}_{Version}_SHA256SUMS, a GPG detached-sign of the checksums via GPG_FINGERPRINT, conventional-commit changelog, and -trimpath/ldflags version stamping. Verified: goreleaser check passes and goreleaser build --snapshot compiles the full matrix (no tag/key). Added terraform-registry-manifest.json (version 1, protocol_versions [6.0], matches spec byte-for-byte). Added .github/workflows/release.yml (on v*.*.* tags: import GPG key, goreleaser release --clean; contents:write). Secrets/keys are human-provided (§14).


## Release v0.1.0 cut (2026-07-13)

Signed release published via release.yml (run 29274385810, success). Verified: 14-target archive matrix + SHA256SUMS + .sig; signature Good against the committed public key (fpr 74F05F879B947F24006761E3FC80F1F128669C1B); linux_amd64 checksum matches; archive embeds the binary + terraform-registry-manifest.json. Dedicated signing key generated; GitHub secrets GPG_PRIVATE_KEY + GPG_FINGERPRINT set (no PASSPHRASE, key has none). Repo renamed to terraform-provider-hcloudimage.
