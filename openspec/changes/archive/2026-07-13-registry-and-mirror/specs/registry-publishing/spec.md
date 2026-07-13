## ADDED Requirements

### Requirement: Registry onboarding is documented
The repository SHALL document the human steps to publish to the Terraform Registry and the
OpenTofu Registry, including the `nivis-project` namespace connection, the GPG public key
registration, and the prerequisite secrets.

#### Scenario: Publishing docs exist
- **WHEN** the publishing docs are read
- **THEN** they describe connecting the repo + GPG key to the Terraform Registry namespace
  and submitting to `opentofu/registry`, and list the §14 prerequisites

#### Scenario: Releases feed the registries
- **WHEN** a signed release is cut by `release.yml`
- **THEN** the documented onboarding causes both registries to ingest the tagged release
  assets
