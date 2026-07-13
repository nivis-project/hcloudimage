## ADDED Requirements

### Requirement: CI runs lint, unit tests, and coverage
The `ci.yml` workflow SHALL run on `pull_request` and `push`, and SHALL run
`golangci-lint`, `go build`, and `go test ./...` with coverage reported to Codecov.

#### Scenario: PR triggers the checks
- **WHEN** a pull request is opened
- **THEN** lint, build, unit tests, and coverage upload run

### Requirement: CI validates examples under both tools
The workflow SHALL run `terraform validate` and `tofu validate` over the examples.

#### Scenario: Examples validated in CI
- **WHEN** CI runs
- **THEN** every example validates under both terraform and tofu

### Requirement: CI runs the hermetic lifecycle test
The workflow SHALL run `nix flake check`, which includes the hermetic NixOS-VM lifecycle
test (the DoD gate).

#### Scenario: Flake check in CI
- **WHEN** CI runs on a runner with VM support
- **THEN** `nix flake check` runs the hermetic-e2e test and must pass

### Requirement: CI enforces docs are up to date
The workflow SHALL fail if `tfplugindocs generate` produces a diff against committed
`docs/`.

#### Scenario: Stale docs fail CI
- **WHEN** the committed docs differ from freshly generated docs
- **THEN** the docs job fails
