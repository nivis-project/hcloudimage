## Why

Milestone 02 wired the resource's CRUD to the `Uploader` but only unit-tested the pure
pieces (schema, validators, mapping). Milestone 03 (BRIEFING.md §13.3, §8.1) proves the
**full resource lifecycle** through the real Terraform plugin protocol using
`terraform-plugin-testing`, with the fake uploader compiled in: apply creates state with a
synthetic `id`; changing a ForceNew attribute forces replacement; changing `labels`/
`description` updates in place; a snapshot deleted out of band is removed from state on
refresh; destroy tears down. This is the behavioural proof the PoC needs before a real
uploader or hermetic VM test is worth building.

## What Changes

- Add `terraform-plugin-testing` as a test dependency.
- Add a test harness that stands up the provider with a **shared, inspectable fake
  uploader** (via the `NewWithUploader` seam) so tests can both drive Terraform and assert
  against the fake's recorded calls.
- Add `resource.Test` lifecycle cases (using the framework's built-in Terraform binary
  discovery) covering: create→state; ForceNew replace on `image_sha256` change; in-place
  update on `labels`/`description` change; out-of-band deletion → removed from state;
  destroy.
- Wire coverage reporting for the `internal/provider` package.

## Capabilities

### New Capabilities
- `resource-lifecycle`: verified plan/apply/update/replace/destroy behaviour of
  `hcloudimage_image` through the Terraform protocol against the fake.

### Modified Capabilities
<!-- none: this adds test coverage; no spec-level behaviour changes -->

## Impact

- New: `internal/provider/image_resource_lifecycle_test.go` and a small test helper.
- Modified: `go.mod`/`go.sum` (test dep), `flake.nix` `vendorHash` if the module set
  changes.
- Gate: `TF_ACC` is **not** required (these use the plugin-testing harness, not real
  cloud); `go test ./...` runs them, and `nix flake check` stays green.
