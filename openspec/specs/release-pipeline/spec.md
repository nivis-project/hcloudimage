# release-pipeline Specification

## Purpose
TBD - created by archiving change release-engineering. Update Purpose after archive.
## Requirements
### Requirement: goreleaser produces registry-shaped signed artifacts
The repository SHALL contain a `.goreleaser.yml` that builds the standard provider os/arch
matrix, archives each build as `{{.ProjectName}}_{{.Version}}_{{.Os}}_{{.Arch}}.zip` with
`terraform-registry-manifest.json` embedded, produces
`{{.ProjectName}}_{{.Version}}_SHA256SUMS`, and GPG-signs the checksums to a `.sig` using
`GPG_FINGERPRINT`.

#### Scenario: Config is valid
- **WHEN** `goreleaser check` runs
- **THEN** it reports the configuration as valid

#### Scenario: Build matrix compiles
- **WHEN** `goreleaser build --snapshot` runs
- **THEN** it produces binaries for the configured os/arch combinations without a tag or
  signing key

### Requirement: Registry manifest declares protocol v6
The repository SHALL contain `terraform-registry-manifest.json` declaring protocol version
`6.0`.

#### Scenario: Manifest content
- **WHEN** the manifest is read
- **THEN** it is `{ "version": 1, "metadata": { "protocol_versions": ["6.0"] } }`

### Requirement: Tag-triggered release workflow
The repository SHALL contain `release.yml` that, on a `v*.*.*` tag, imports the GPG key and
runs `goreleaser release --clean`, producing the signed GitHub release the registries
ingest.

#### Scenario: Release runs on a version tag
- **WHEN** a `v*.*.*` tag is pushed
- **THEN** the workflow imports the GPG key and runs goreleaser to publish signed artifacts

### Requirement: CHANGELOG follows Keep a Changelog
The repository SHALL maintain `CHANGELOG.md` in Keep a Changelog style, starting with the
`0.1.0` entry.

#### Scenario: Initial version documented
- **WHEN** `CHANGELOG.md` is read
- **THEN** it contains a `0.1.0` section describing the initial release

