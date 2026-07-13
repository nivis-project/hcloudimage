# hermetic-e2e Specification

## Purpose
TBD - created by archiving change hermetic-nixos-vm-test. Update Purpose after archive.
## Requirements
### Requirement: Hermetic lifecycle test gates nix flake check
The project SHALL provide a `checks.hermetic-e2e` (on Linux) that boots a NixOS VM,
installs the built provider into an in-VM filesystem mirror, and runs the full Terraform
protocol lifecycle with the fake uploader — with no cloud access — as part of
`nix flake check`.

#### Scenario: Flake check runs the hermetic test
- **WHEN** `nix flake check` runs on a Linux system with a VM builder available
- **THEN** the hermetic-e2e check builds and passes

### Requirement: Lifecycle proven under both terraform and tofu
The hermetic test SHALL run `init → plan → apply → destroy` under both `terraform` and
`tofu`, and SHALL assert create, ForceNew replacement, in-place update, and destroy.

#### Scenario: Create then destroy under each binary
- **WHEN** the lifecycle runs under terraform and again under tofu
- **THEN** apply creates exactly one resource with a synthetic id, and destroy leaves an
  empty state, for both binaries

#### Scenario: ForceNew vs in-place asserted from plan JSON
- **WHEN** the ForceNew trigger changes
- **THEN** the plan's resource change actions are `["delete","create"]` (replace)
- **WHEN** only labels change
- **THEN** the plan's resource change action is `["update"]` (in-place)

### Requirement: Fake uploader keeps the test cloud-free
The hermetic test SHALL run with `HCLOUDIMAGE_FAKE=1` so no Hetzner API or SSH access is
attempted, keeping the check hermetic and free.

#### Scenario: No network required
- **WHEN** the VM runs the lifecycle with no network access
- **THEN** every step succeeds using the fake uploader

