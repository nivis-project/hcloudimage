## ADDED Requirements

### Requirement: Flake exposes the test-image fixtures
The flake SHALL expose `packages.test-image-x86` and `packages.test-image-arm` (plain nix,
no flake-utils) so the acceptance job can build the fixtures it uploads.

#### Scenario: Fixture packages present
- **WHEN** the flake packages are listed on x86_64-linux
- **THEN** `test-image-x86` and `test-image-arm` are present
