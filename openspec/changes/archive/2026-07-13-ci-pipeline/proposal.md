## Why

Milestones 01–05 built and proved the provider locally. Milestone 06 (BRIEFING.md §13.6,
§9) makes the quality bar enforceable on every change via GitHub Actions `ci.yml`: lint,
build, unit tests with coverage to Codecov, `terraform`/`tofu validate` over the examples,
the hermetic `nix flake check`, and a `tfplugindocs` diff check. To make the docs-diff
check meaningful, this change also generates and commits the initial `docs/`.

## What Changes

- Generate `docs/` with `tfplugindocs` (from schema `MarkdownDescription` + examples) and
  commit it; add a `templates/`-free default layout.
- Add `.golangci.yml` pinning the linters CI runs (matching what the devShell already runs).
- Add `.github/workflows/ci.yml` running on `pull_request` and `push`:
  - lint (`golangci-lint`), `go build`, `go test ./...` with coverage → Codecov,
  - a matrix `terraform validate` + `tofu validate` over `examples/` (via the mirror script),
  - `nix flake check` (runs the hermetic lifecycle test),
  - a docs-diff job asserting `tfplugindocs generate` produces no diff.
- Use the Nix devShell / `DeterminateSystems` nix installer so CI uses the same toolchain
  versions as local.

## Capabilities

### New Capabilities
- `ci-pipeline`: the `ci.yml` workflow enforcing lint/unit/validate/hermetic/docs on PRs.
- `generated-docs`: committed `tfplugindocs` output, diff-checked in CI.

### Modified Capabilities
<!-- none: adds CI + docs, no provider behaviour change -->

## Impact

- New: `.github/workflows/ci.yml`, `.golangci.yml`, `docs/**`.
- Gate: the workflow file is valid and its steps mirror the locally-verified commands;
  docs are committed and `tfplugindocs generate` is diff-clean.
