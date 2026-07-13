---
# hcloudimage-pdyi
title: Nix flake foundation
status: todo
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:11:14Z
parent: hcloudimage-c2wx
---

flake.nix with plain-nix forAllSystems (NO flake-utils): devShells.default (go, golangci-lint, terraform, opentofu, tfplugindocs, goreleaser, gnumake, hcloud-upload-image), packages.default via buildGoModule with pinned vendorHash. Gate: nix develop + nix build. (BRIEFING.md §7)

OpenSpec change: scaffold-provider-and-flake
