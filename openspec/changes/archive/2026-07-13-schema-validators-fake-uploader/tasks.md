## 1. Uploader seam

- [x] 1.1 `uploader.go`: `Uploader` interface (Upload/Delete/Get/Find), `UploadRequest`, `SnapshotInfo`, shared constants + the `apricote.de/created-by` default-label helper
- [x] 1.2 `uploader_fake.go`: in-memory fake keyed by synthetic IDs; records calls; merges created-by label; supports marking snapshots deleted out of band
- [x] 1.3 Unit tests for the fake (create→get→delete, out-of-band delete, label merge)

## 2. Validators

- [x] 2.1 `validators.go`: `image_sha256`-required-iff-`image_path` config validator; label-value `/` validator
- [x] 2.2 Wire `ExactlyOneOf(image_url, image_path)` and enum validators (architecture/compression/format)
- [x] 2.3 Unit tests: both-sources, neither-source, path-without-sha, url-without-sha, bad enum, bad label

## 3. hcloudimage_image resource

- [x] 3.1 Full §3.2 schema with MarkdownDescription on every attribute; computed `id`/`effective_labels`; `timeouts` block
- [x] 3.2 RequiresReplace plan modifiers on all ForceNew attributes; none on `description`/`labels`
- [x] 3.3 Create/Read/Update/Delete against the `Uploader` (Update does in-place metadata for description/labels; Read removes missing from state)
- [x] 3.4 Config→UploadRequest mapping helper
- [x] 3.5 Unit tests: schema shape, plan-modifier ForceNew vs in-place, config→request mapping across all compressions/both sources/both arches, effective_labels merge, read-removes-missing

## 4. hcloudimage_snapshot data source

- [x] 4.1 `image_data_source.go`: §3.3 schema (`id` xor `with_selector`, `most_recent`, computed fields)
- [x] 4.2 Resolve via `uploader.Find`; ambiguity error unless `most_recent`
- [x] 4.3 Unit tests: by-id, both-set-rejected, ambiguous-without/with-most_recent

## 5. Provider wiring

- [x] 5.1 Provider `Configure` constructs the fake `Uploader` and sets it on ResourceData/DataSourceData
- [x] 5.2 Register `hcloudimage_snapshot` in `DataSources`; resource/data source read the uploader from configure data
- [x] 5.3 Update provider unit tests (data source registered)

## 6. Close out

- [x] 6.1 `go test ./... -cover` green with high `internal/provider` coverage; `go build ./...`; `nix flake check`
- [x] 6.2 Archive OpenSpec change; complete beans epics + milestone 02; commit with jj (Pim Snel, no self-promotion)
