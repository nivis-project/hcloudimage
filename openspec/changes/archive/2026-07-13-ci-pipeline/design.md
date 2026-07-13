## Context

Everything CI runs has already been verified locally in earlier milestones; this change
just codifies it in `ci.yml` and adds the committed docs the diff-check compares against.

## Decisions

### Nix-first CI
- Use `DeterminateSystems/nix-installer-action` + `magic-nix-cache-action` so CI uses the
  exact toolchain from the flake (`nix develop -c ...`), avoiding version drift between
  local and CI. This also gives `nix flake check` for free.

### Jobs
1. **lint-build-test** (ubuntu): `nix develop -c golangci-lint run ./...`,
   `nix develop -c go build ./...`, `nix develop -c go test ./... -covermode=atomic
   -coverprofile=coverage.out`, upload to Codecov (`codecov/codecov-action`, token optional
   for public repos). The lifecycle tests need a terraform binary + `TF_ACC=1`; set both.
2. **validate-examples** (ubuntu): `nix build .#default` then
   `nix develop -c bash scripts/validate-examples.sh`. Runs `terraform`+`tofu validate`.
3. **flake-check** (ubuntu, KVM): `nix flake check` — runs the hermetic lifecycle test.
   Uses a KVM-enabled runner label or the default (GitHub's ubuntu runners support nested
   virt for NixOS tests via the standard KVM group workaround; document if it needs a
   self-hosted runner).
4. **docs** (ubuntu): `nix develop -c tfplugindocs generate` then `git diff --exit-code
   docs/` — fails if generated docs drift from committed.

### golangci config
- `.golangci.yml` enabling the same linters the local run used (errcheck, staticcheck,
  govet, etc., defaults are fine) so local `golangci-lint run` and CI agree.

### Docs
- Commit `docs/` (index + resources/image + data-sources/snapshot). tfplugindocs pulls
  `MarkdownDescription` and the `examples/` files, so docs stay in sync as long as the
  diff-check passes.

## Risks / Trade-offs

- NixOS VM tests on hosted GitHub runners need KVM; if unavailable, the flake-check job is
  documented as requiring a self-hosted or nested-virt runner. The check still runs locally
  and is the DoD gate regardless.
- Codecov token: public-repo uploads work tokenless; `CODECOV_TOKEN` is a documented human
  secret (BRIEFING §14) for private/rate-limited cases.
