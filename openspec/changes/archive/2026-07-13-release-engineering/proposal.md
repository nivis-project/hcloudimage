## Why

Milestone 08 (BRIEFING.md §13.8, §10, §11) makes the provider releasable to both the
Terraform and OpenTofu registries via identically-signed GitHub release artifacts. It adds
the goreleaser config (standard provider build matrix + GPG-signed SHA256SUMS), the
registry manifest, the tag-triggered release workflow, and the CHANGELOG. The signing
**keys and secrets** are human-provided (BRIEFING.md §14), so `release.yml` is complete and
correctly wired; the actual signed release is cut by a human with the GPG key.

## What Changes

- `.goreleaser.yml` (v2): builds `linux/darwin/windows/freebsd × amd64/arm64/arm/386`,
  archive name `{{.ProjectName}}_{{.Version}}_{{.Os}}_{{.Arch}}.zip`, produces
  `{{.ProjectName}}_{{.Version}}_SHA256SUMS` + a GPG `.sig` (`signs:` using
  `GPG_FINGERPRINT`), embeds `terraform-registry-manifest.json` in each archive, changelog
  from conventional commits, correct `ldflags` version stamping.
- `terraform-registry-manifest.json` = `{ "version": 1, "metadata": { "protocol_versions":
  ["6.0"] } }`.
- `.github/workflows/release.yml`: on tag `v*.*.*`, runs `goreleaser release --clean` with
  the GPG import + `GPG_FINGERPRINT`/`PASSPHRASE`/`GPG_PRIVATE_KEY` secrets.
- `CHANGELOG.md` (Keep a Changelog) with the `v0.1.0` entry.
- Verify the goreleaser config with `goreleaser check` and a `--snapshot` build (no
  signing) so the build matrix is proven locally.

## Capabilities

### New Capabilities
- `release-pipeline`: goreleaser config, registry manifest, and `release.yml` producing
  signed, registry-shaped artifacts.

### Modified Capabilities
<!-- none: adds release tooling; docs were generated in milestone 06 -->

## Impact

- New: `.goreleaser.yml`, `terraform-registry-manifest.json`,
  `.github/workflows/release.yml`, `CHANGELOG.md`.
- Gate: `goreleaser check` passes; a `goreleaser build --snapshot` produces the expected
  os/arch matrix; the manifest matches the spec; signing/secrets are documented human
  steps (§14).
