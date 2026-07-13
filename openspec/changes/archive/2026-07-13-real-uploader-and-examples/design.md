## Context

The `Uploader` seam already exists and is proven against the fake. This change adds the
production implementation and the examples, without disturbing the tested behaviour.

## Decisions

### Real uploader (uploader_hcloud.go), verified against v2.0.1
- Construct: `hcloudimages.NewClient(hcloud.NewClient(hcloud.WithToken(token), opts...))`.
  Honour `endpoint` via `hcloud.WithEndpoint` when set; `poll_interval` via
  `hcloud.WithPollBackoffFunc`/opts if provided (best-effort passthrough).
- `Upload`: build `hcloudimages.UploadOptions{ WriteOptions{ImageURL|ImageReader,
  ImageCompression, ImageFormat, ImageSize}, Architecture, ServerType, Description, Labels,
  Location }`. For `image_path`, open the file and pass it as `ImageReader` (close after).
  For `image_url`, parse to `*url.URL`. Returns `*hcloud.Image`; map `.ID` and
  `.Labels` into the result.
- Mapping tables (single source of truth in uploader.go already):
  `none→CompressionNone, bz2→CompressionBZ2, xz→CompressionXZ, zstd→CompressionZSTD`;
  `raw→FormatRaw, qcow2→FormatQCOW2`; `x86→hcloud.ArchitectureX86, arm→ArchitectureARM`.
- `ServerType`/`Location`: only set the pointer when the user overrode them; otherwise let
  the library apply its per-arch defaults (do not hardcode here to avoid drift).
- `Delete`: `client.Image.Delete(ctx, &hcloud.Image{ID: id})`.
- `Get`: `client.Image.GetByID`; return nil when not found (drives state removal).
- `UpdateMetadata`: `client.Image.Update` with description + labels; return merged labels.
- `Find`: by id → `GetByID`; by selector → `Image.AllWithOpts(ListOpts{LabelSelector})`,
  enforce exactly-one unless most_recent (then newest by `Created`).
- **No public skip-cleanup**: never set `DebugSkipResourceCleanup` from schema. Gate it
  behind `HCLOUDIMAGE_DEBUG_SKIP_CLEANUP=1` env only, for local debugging.

### Provider uploader selection
- `newUploader(cfg)`: if `HCLOUDIMAGE_FAKE=1` → fake (hermetic VM test switch). Else if a
  token is present → real uploader. Else → fake (so `terraform validate`/plan with no
  token, and unit tests, never need real credentials). Document this precedence.

### Examples
- Three files under `examples/` matching BRIEFING.md §5. Pin `hcloud` to `~> 1.48` and this
  provider to `~> 0.1`. The resource example uses `image_path` + `filesha256(...)` and
  composes an `hcloud_server` booting from the snapshot id.
- Validation: a dev/CI step runs, for each example dir, `terraform init -backend=false` +
  `terraform validate`, then the same with `tofu`. Uses a local dev override / the built
  provider so `init` resolves `nivis-project/hcloudimage` without a registry.

## Risks / Trade-offs

- `terraform validate` needs the provider installed; use a `dev_overrides` block pointing
  at the Nix-built binary (documented) so validation is offline and hermetic.
- Field/const names confirmed against v2.0.1; pin the version in go.mod so a v2 minor can't
  silently change the mapping.
