# provider-scaffold Specification

## Purpose
TBD - created by archiving change scaffold-provider-and-flake. Update Purpose after archive.
## Requirements
### Requirement: Provider server starts under protocol v6
The provider SHALL be served via `terraform-plugin-framework` using Terraform Plugin
Protocol v6, addressed as `registry.terraform.io/nivis-project/hcloudimage`.

#### Scenario: Binary builds and serves
- **WHEN** `go build ./...` is run
- **THEN** it succeeds and produces a `terraform-provider-hcloudimage` binary
- **AND** the binary registers at least one resource so the server starts

### Requirement: Provider configuration schema
The provider SHALL expose the configuration schema of BRIEFING.md §3.1: an optional
sensitive `token` that falls back to the `HCLOUD_TOKEN` environment variable when unset,
an optional `endpoint` override, and an optional `poll_interval` duration string.

#### Scenario: Token falls back to environment
- **WHEN** the provider block sets no `token`
- **AND** `HCLOUD_TOKEN` is set in the environment
- **THEN** the provider resolves its token from the environment variable

#### Scenario: Endpoint override accepted
- **WHEN** the provider block sets `endpoint`
- **THEN** configuration succeeds and the value is retained for client construction

### Requirement: Hermetic Nix build
The project SHALL build reproducibly with Nix flakes using plain nix for multi-system
support (no `flake-utils`), exposing a `devShells.default` and a `packages.default` built
via `buildGoModule` with a pinned `vendorHash`.

#### Scenario: Dev shell provides the toolchain
- **WHEN** `nix develop` is entered
- **THEN** `go`, `golangci-lint`, `terraform`, `opentofu`, `terraform-plugin-docs`,
  `goreleaser`, `gnumake`, and `hcloud-upload-image` are on `PATH`

#### Scenario: Provider builds via Nix
- **WHEN** `nix build .#default` is run
- **THEN** it produces the provider binary hermetically with no network beyond
  fixed-output vendoring

#### Scenario: Flake check passes
- **WHEN** `nix flake check` is run
- **THEN** all checks pass

### Requirement: Provider injects an Uploader and registers the data source
The provider SHALL construct an `Uploader` implementation and pass it to resources and
data sources via configure data, and SHALL register the `hcloudimage_snapshot` data
source alongside the `hcloudimage_image` resource.

#### Scenario: Uploader reaches the resource
- **WHEN** the provider's `Configure` runs
- **THEN** an `Uploader` is available to the resource and data source through configure data

#### Scenario: Data source is registered
- **WHEN** the provider reports its data sources
- **THEN** `hcloudimage_snapshot` is present

### Requirement: Uploader selection prefers real, falls back to fake
The provider SHALL select the real uploader when a token is present and `HCLOUDIMAGE_FAKE`
is unset, force the fake when `HCLOUDIMAGE_FAKE=1`, and otherwise fall back to the fake so
that validation, planning without a token, and tests need no real credentials.

#### Scenario: No token falls back to fake
- **WHEN** the provider is configured with no token and `HCLOUDIMAGE_FAKE` unset
- **THEN** the fake uploader is used, so `terraform validate`/plan needs no credentials

#### Scenario: Fake forced by environment
- **WHEN** `HCLOUDIMAGE_FAKE=1`
- **THEN** the fake uploader is used even if a token is present

