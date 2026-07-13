## ADDED Requirements

### Requirement: Flake exposes the hermetic check on Linux
The flake SHALL expose `checks.hermetic-e2e` on Linux systems, defined with plain nix (no
flake-utils), so `nix flake check` includes it where a VM builder is available.

#### Scenario: Check present on Linux
- **WHEN** the flake outputs are inspected on x86_64-linux
- **THEN** `checks.x86_64-linux.hermetic-e2e` exists
