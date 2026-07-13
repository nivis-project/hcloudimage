## Context

Second change. Builds the full provider surface on the scaffold, driven by a fake
uploader so all behaviour is unit-testable with no network (BRIEFING.md §4.1, §8.1).

## Decisions

### The Uploader seam (BRIEFING.md §4.1)
```go
type UploadRequest struct {
    ImageURL     string            // one of URL/Path
    ImagePath    string
    ImageSHA256  string
    Architecture string            // "x86" | "arm"
    Compression  string            // "none" | "bz2" | "xz" | "zstd"
    Format       string            // "raw" | "qcow2"
    ServerType   string
    Location     string
    ImageSize    int64
    Description  string
    Labels       map[string]string
}
type SnapshotInfo struct {
    ID              int64
    Description     string
    Architecture    string
    Labels          map[string]string
    Created         string
    Name            string
}
type Uploader interface {
    Upload(ctx context.Context, req UploadRequest) (imageID int64, effectiveLabels map[string]string, err error)
    Delete(ctx context.Context, imageID int64) error
    Get(ctx context.Context, imageID int64) (*SnapshotInfo, error)
    // Find resolves a data-source lookup (by id or label selector).
    Find(ctx context.Context, byID int64, selector string, mostRecent bool) (*SnapshotInfo, error)
}
```
- `uploader_fake.go`: in-memory map keyed by synthetic incrementing IDs; records calls;
  always merges the library's `apricote.de/created-by=hcloud-upload-image` label into
  `effectiveLabels` so tests see the real merged set; can be seeded/marked to simulate
  "snapshot deleted out of band" (Get returns not-found).
- The resource/data source read the `Uploader` from configure data. The provider decides
  which implementation to inject — the fake for now (real impl in milestone 04).

### Resource schema and plan behaviour (BRIEFING.md §3.2)
- ForceNew attributes get a `stringplanmodifier.RequiresReplace()` (or the int64/map
  equivalent): `image_url`, `image_path`, `image_sha256`, `architecture`, `compression`,
  `format`, `server_type`, `location`, `image_size`.
- `description` and `labels` have **no** replace modifier → update in place: `Update`
  calls the fake's metadata update path, not `Upload`.
- `id` and `effective_labels` are `Computed`.
- `timeouts` block via `terraform-plugin-framework-timeouts` (create/read/delete).
- `Read`: `uploader.Get(id)`; if not found, `resp.State.RemoveResource` (no error).
- `Delete`: `uploader.Delete(id)`.

### Validators (BRIEFING.md §3.2, in validators.go)
- Config-level `ConfigValidators`: `resourcevalidator.ExactlyOneOf(image_url, image_path)`
  and a custom validator enforcing `image_sha256` set iff `image_path` set.
- Attribute-level: enum validators for `architecture`/`compression`/`format`; a label
  validator rejecting `/` in label values (Hetzner rule).

### Data source (BRIEFING.md §3.3)
- `id` xor `with_selector` (ExactlyOneOf); `most_recent` optional; on selector resolving to
  many with `most_recent=false`, error on ambiguity. Delegates to `uploader.Find`.
- Computed: `name`, `description`, `labels`, `architecture`, `created`.

### Compression/format/arch constants
- Keep string constants + validation sets in one place so the real uploader (milestone 04)
  maps them to the library's typed constants without re-deriving the allowed set.

## Risks / Trade-offs

- The fake's label-merge must mirror the real library exactly, or milestone 04 will show a
  diff. Encode the `apricote.de/created-by` default in a shared helper used by both.
- `Find` on the interface is slightly ahead of the real impl, but keeping it on the seam
  now means the data source is testable in this milestone.
