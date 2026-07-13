---
# hcloudimage-mw1v
title: ci.yml workflow
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T15:29:30Z
parent: hcloudimage-lofb
blocked_by:
    - hcloudimage-r4d6
---

PR+push: golangci-lint, go build, go test with coverage→Codecov, terraform+tofu validate matrix over examples/, nix flake check, tfplugindocs diff check. (BRIEFING.md §9)

## Summary of Changes

Added .github/workflows/ci.yml (on pull_request and push to main) with four jobs, all Nix-based (DeterminateSystems installer + magic cache) so CI uses the same toolchain as local: lint-build-test (golangci-lint, go build, go test with TF_ACC lifecycle tests + coverage to Codecov), validate-examples (nix build + scripts/validate-examples.sh under terraform and tofu), flake-check (nix flake check incl. the hermetic VM test, with KVM enablement + documented runner fallback), and docs (tfplugindocs generate then git diff --exit-code). Added .golangci.yml (v2, errcheck/govet/staticcheck/unused/misspell/unconvert/revive with the doc-comment nags disabled). Generated and committed docs/ (index, resources/image, data-sources/snapshot); regeneration is deterministic (no drift).
