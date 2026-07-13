## Why

Milestone 02 (BRIEFING.md §13.2) turns the empty scaffold into a fully-described
provider surface. It defines the complete `hcloudimage_image` resource schema with its
config validators and ForceNew/in-place plan behaviour (§3.2), the `hcloudimage_snapshot`
data source (§3.3), and the `Uploader` seam (§4.1) that makes everything testable without
touching real HTTP/SSH. No real cloud calls yet — the resource is driven by an in-memory
fake — but the schema, validation, and CRUD wiring become real.

## What Changes

- Full `hcloudimage_image` schema (§3.2): `image_url`/`image_path`, `image_sha256`,
  `architecture`, `compression`, `format`, `server_type`, `location`, `image_size`,
  `description`, `labels`, computed `id` and `effective_labels`, plus a `timeouts` block.
- ForceNew semantics via `RequiresReplace` plan modifiers on all ForceNew attributes;
  `description` and `labels` update in place.
- Config validators (§3.2): exactly one of `image_url`/`image_path`; `image_sha256`
  required iff `image_path`; Hetzner label rules (no `/` in values); enum validation for
  `architecture`, `compression`, `format`.
- `hcloudimage_snapshot` data source (§3.3): lookup by `id` or `with_selector`
  (+`most_recent`); computed fields.
- `Uploader` interface (`Upload`/`Delete`/`Get`) + `uploader_fake.go` in-memory fake, plus
  the `UploadRequest`/`SnapshotInfo` types. Resource and data source depend only on the
  interface; the provider wires the fake for now (real impl is milestone 04).
- Unit tests: schema correctness, every validator, plan-modifier behaviour, config→request
  mapping across all compressions/both sources/both architectures.

## Capabilities

### New Capabilities
- `image-resource`: the `hcloudimage_image` schema, validators, and ForceNew/in-place plan
  behaviour.
- `snapshot-data-source`: the `hcloudimage_snapshot` lookup-by-id/selector data source.
- `uploader-seam`: the `Uploader` interface and in-memory fake that decouple the resource
  from the real hcloud library.

### Modified Capabilities
- `provider-scaffold`: the provider now wires an `Uploader` into resource/data-source
  configure data and registers the data source.

## Impact

- New: `internal/provider/uploader.go`, `uploader_fake.go`, `validators.go`,
  `image_data_source.go`, and `*_test.go`.
- Modified: `internal/provider/provider.go` (wire fake uploader, register data source),
  `image_resource.go` (full schema + CRUD against the interface).
- Gate: `go test ./...` green with high coverage of schema and validators.
