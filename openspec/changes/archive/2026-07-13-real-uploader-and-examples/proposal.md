## Why

Milestones 02–03 built and proved the whole provider against a fake. Milestone 04
(BRIEFING.md §13.4) wires the **real** upload engine behind the same `Uploader` seam using
`github.com/apricote/hcloud-upload-image/hcloudimages/v2`, and adds the runnable HCL
examples that must pass `terraform validate` and `tofu validate`. After this, the provider
does real work when given a token, while all tests still run hermetically against the fake.

## What Changes

- `uploader_hcloud.go`: a real `Uploader` wrapping `hcloudimages.Client` + `hcloud.Client`.
  Maps `UploadRequest` → `hcloudimages.UploadOptions`/`WriteOptions` (both sources, all four
  compressions, `raw`/`qcow2`, arch → `hcloud.ArchitectureX86/ARM`, server_type/location/
  description/labels). `Delete`/`Get`/`Find`/`UpdateMetadata` via the `hcloud` client's
  `Image` API. Never exposes a public skip-cleanup knob (gate any debug skip behind an env
  var only).
- Provider wiring: when a token is present, `newUploader` builds the real uploader;
  otherwise (and in tests) it falls back to the fake. An env override
  (`HCLOUDIMAGE_FAKE=1`) forces the fake for the hermetic VM test (milestone 05).
- Examples: `examples/provider/provider.tf`, `examples/resources/hcloudimage_image/resource.tf`
  (the §5 example verbatim in behaviour), `examples/data-sources/hcloudimage_snapshot/data-source.tf`.
- A GNUmakefile/dev target and a check that runs `terraform validate` and `tofu validate`
  over `examples/`.

## Capabilities

### New Capabilities
- `real-uploader`: the hcloudimages/v2-backed `Uploader` and its config→library mapping.
- `hcl-examples`: runnable examples that validate under both Terraform and OpenTofu.

### Modified Capabilities
- `provider-scaffold`: `newUploader` now selects real-vs-fake based on token / env.

## Impact

- New: `internal/provider/uploader_hcloud.go` (+ mapping unit tests), `examples/**`.
- Modified: `provider.go` (uploader selection), `go.mod`/`go.sum`, `flake.nix` vendorHash.
- Gate: `terraform validate` **and** `tofu validate` pass on every example; unit tests for
  the compression/format/arch mapping green; `nix flake check` stays green.
