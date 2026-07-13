# nix-mirror-consumption Specification

## Purpose
TBD - created by archiving change registry-and-mirror. Update Purpose after archive.
## Requirements
### Requirement: Provider is consumable from the Nix mirror without a registry
`packages.provider-mirror` SHALL produce a filesystem-mirror layout that lets a real
consumer configuration resolve `nivis-project/hcloudimage` at a pinned version without any
public registry, under both `terraform` and `tofu`.

#### Scenario: init and plan resolve from the mirror
- **WHEN** a consumer config with a pinned `hcloudimage` version uses a `filesystem_mirror`
  pointing at `nix build .#provider-mirror` and runs `init` + `plan` (with `HCLOUDIMAGE_FAKE=1`)
- **THEN** both `terraform` and `tofu` resolve the provider from the mirror and succeed with
  no registry access and no token

### Requirement: Mirror consumption is documented
The repository SHALL document how Nivis consumes the provider from the Nix mirror,
including the `filesystem_mirror` CLI-config recipe and the `dev_overrides` alternative.

#### Scenario: Consumption docs exist
- **WHEN** the docs are read
- **THEN** they give the exact CLI-config recipe using the `nix build .#provider-mirror`
  output

