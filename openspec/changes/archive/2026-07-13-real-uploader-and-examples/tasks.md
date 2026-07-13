## 1. Real uploader

- [x] 1.1 Add `github.com/apricote/hcloud-upload-image/hcloudimages/v2` (pin v2.0.1); update vendorHash
- [x] 1.2 `uploader_hcloud.go`: construct client from token/endpoint; implement Upload (both sources, all compressions/formats, arch, optional server_type/location, labels/description)
- [x] 1.3 Implement Delete/Get/UpdateMetadata/Find via the hcloud Image API; Get returns nil when not found
- [x] 1.4 Gate any debug skip-cleanup behind `HCLOUDIMAGE_DEBUG_SKIP_CLEANUP` env only (no schema attribute)
- [x] 1.5 Unit tests for the pure mapping helpers (compression/format/architecture → library constants)

## 2. Provider wiring

- [x] 2.1 `newUploader`: `HCLOUDIMAGE_FAKE=1` → fake; token present → real; else fake
- [x] 2.2 Unit test the selection precedence

## 3. Examples

- [x] 3.1 `examples/provider/provider.tf`
- [x] 3.2 `examples/resources/hcloudimage_image/resource.tf` (§5 behaviour: image_path + filesha256, composes hcloud_server from the snapshot id)
- [x] 3.3 `examples/data-sources/hcloudimage_snapshot/data-source.tf`
- [x] 3.4 Dev/CI validation: for each example, `terraform validate` and `tofu validate` (offline via dev_overrides to the Nix-built provider)

## 4. Close out

- [x] 4.1 `go build ./...`, `go test ./...`, examples validate under both binaries, `nix flake check`
- [x] 4.2 Archive OpenSpec change; complete beans epics + milestone 04; commit with jj (Pim Snel, no self-promotion)
