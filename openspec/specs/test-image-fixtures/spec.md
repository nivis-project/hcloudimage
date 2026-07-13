# test-image-fixtures Specification

## Purpose
TBD - created by archiving change fixtures-and-acceptance. Update Purpose after archive.
## Requirements
### Requirement: Reproducible Alpine fixtures with a baked SSH key
The flake SHALL expose `packages.test-image-x86` and `packages.test-image-arm` that
produce reproducible minimal Alpine `.raw.xz` images (amd64 and aarch64) with a throwaway
SSH public key baked into `/root/.ssh/authorized_keys`, `sshd` enabled, and DHCP
configured — cloud-init-free.

#### Scenario: Fixture builds to a compressed raw image
- **WHEN** `nix build .#test-image-x86` runs
- **THEN** it produces a `.raw.xz` file hermetically from a pinned Alpine image

#### Scenario: Build is hermetic and input-pinned
- **WHEN** a fixture is built
- **THEN** its inputs are pinned (fixed-output Alpine fetch + committed key) and the
  customize step runs fully offline, so the same inputs always produce a functionally
  equivalent image (disk images are not byte-reproducible due to filesystem metadata)

#### Scenario: arm fixture is aarch64
- **WHEN** `packages.test-image-arm` is built
- **THEN** the image is aarch64, not a silent x86 fallback

