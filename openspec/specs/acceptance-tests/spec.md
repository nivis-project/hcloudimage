# acceptance-tests Specification

## Purpose
TBD - created by archiving change fixtures-and-acceptance. Update Purpose after archive.
## Requirements
### Requirement: Acceptance tests compose both providers and prove reachability
The acceptance suite SHALL, when `TF_ACC=1` and `HCLOUD_TOKEN` are set, upload a fixture
via `hcloudimage_image`, boot an `hcloud_server` from the snapshot id, and assert the guest
is reachable by SSHing from the runner with the baked throwaway key and reading
`/etc/os-release` — not merely that the server reports `running`.

#### Scenario: Skipped without credentials
- **WHEN** `TF_ACC` or `HCLOUD_TOKEN` is unset
- **THEN** the acceptance tests are skipped, keeping `go test ./...` green without secrets

#### Scenario: Reachability asserted via SSH
- **WHEN** the acceptance test runs against a real project
- **THEN** it SSHes into the booted server and confirms the guest OS, and cleans up even on
  failure

### Requirement: Acceptance covers both architectures with cost controls
The suite SHALL cover both `x86` (cx22 + amd64 fixture) and `arm` (cax11 + aarch64
fixture), with arm toggle-gated, cheapest server types, short timeouts, a pinned `hcloud`
provider version, and guaranteed cleanup.

#### Scenario: arm run uses the aarch64 fixture
- **WHEN** the arm acceptance path runs
- **THEN** it uses the aarch64 fixture and a cax11 server, not an x86 fallback

### Requirement: Gated acceptance and cleanup workflows
The repository SHALL provide `acceptance.yml` (workflow_dispatch + push-to-main + nightly,
never on fork PRs, concurrency-limited, always-run cleanup) and `cleanup.yml` (nightly
label-scoped orphan sweep).

#### Scenario: Never runs on fork PRs
- **WHEN** a pull request from a fork is opened
- **THEN** the acceptance workflow does not run (secrets are not exposed)

