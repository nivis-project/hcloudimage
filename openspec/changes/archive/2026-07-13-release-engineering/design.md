## Context

The Terraform and OpenTofu registries both ingest GitHub release assets produced by
goreleaser and verify a GPG signature over the SHA256SUMS. The config must match the
registry's expected asset shape exactly (BRIEFING.md §10) — that shape is the error-prone
bit, so it is treated as a requirement.

## Decisions

### .goreleaser.yml (schema v2)
- `builds`: single binary `terraform-provider-hcloudimage`, `CGO_ENABLED=0`,
  `ldflags = -s -w -X main.version={{.Version}}`, `goos: [linux, darwin, windows, freebsd]`,
  `goarch: [amd64, arm64, arm, '386']`, with the standard provider `ignore` combinations
  (darwin/386, darwin/arm, freebsd/arm64 as upstream templates do — keep the documented
  matrix).
- `archives`: `formats: [zip]`, `name_template:
  '{{.ProjectName}}_{{.Version}}_{{.Os}}_{{.Arch}}'`, and `files: [terraform-registry-manifest.json]`
  so the manifest ships inside every archive.
- `checksum`: `name_template: '{{.ProjectName}}_{{.Version}}_SHA256SUMS'`, algorithm sha256.
- `signs`: one entry signing the checksums with GPG detached-sign
  (`artifacts: checksum`, `signature: '${artifact}.sig'`,
  `args: [--batch, --local-user, '{{.Env.GPG_FINGERPRINT}}', --output, '${signature}',
  --detach-sign, '${artifact}']`).
- `changelog`: `use: github`, group by conventional-commit prefixes.
- `release`: draft optional; `disable` false.

### terraform-registry-manifest.json
- Exactly `{ "version": 1, "metadata": { "protocol_versions": ["6.0"] } }`.

### release.yml
- `on: push: tags: ['v*.*.*']`. Steps: checkout (fetch-depth 0 for changelog), nix
  installer, import GPG key (`crazy-max/ghaction-import-gpg` with `GPG_PRIVATE_KEY` +
  `PASSPHRASE`), then `nix develop --command goreleaser release --clean` with env
  `GPG_FINGERPRINT` and `GITHUB_TOKEN`. `permissions: contents: write` for the release.

### Verification without secrets
- `goreleaser check` validates the config.
- `goreleaser build --snapshot --clean --single-target` (and a full `--snapshot` where
  feasible) proves the build compiles across the matrix without needing a tag or GPG key.
- Signing is exercised only in a real tagged release by a human holding the key (§14).

### CHANGELOG.md
- Keep a Changelog format; `## [0.1.0]` documenting the initial provider surface, tests,
  Nix build, CI, and release tooling.

## Risks / Trade-offs

- goreleaser v2 renamed some fields (`format` → `formats`, `archives.builds` semantics);
  `goreleaser check` catches these — run it as the gate.
- Cutting the actual `v0.1.0` tag + signed release needs the GPG key and registry
  onboarding (human, §14). The agent produces everything up to, but not including, the
  key-holding release step, and documents it.
