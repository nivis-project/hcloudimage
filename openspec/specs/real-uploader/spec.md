# real-uploader Specification

## Purpose
TBD - created by archiving change real-uploader-and-examples. Update Purpose after archive.
## Requirements
### Requirement: Real uploader wraps hcloudimages/v2
The provider SHALL provide an `Uploader` implementation backed by
`github.com/apricote/hcloud-upload-image/hcloudimages/v2`, wired only through the
`Uploader` interface so the resource never references the library directly.

#### Scenario: Real uploader constructed with a token
- **WHEN** the provider is configured with a non-empty token and `HCLOUDIMAGE_FAKE` is unset
- **THEN** the real hcloudimages/v2-backed uploader is used

### Requirement: Config maps to library options
The real uploader SHALL map every supported compression (`none`/`bz2`/`xz`/`zstd`), format
(`raw`/`qcow2`), and architecture (`x86`/`arm`) to the corresponding hcloudimages/v2
constants, and SHALL use `image_url` as `ImageURL` and `image_path` as `ImageReader`.

#### Scenario: Compression constants map
- **WHEN** each compression value is mapped
- **THEN** it maps to CompressionNone/BZ2/XZ/ZSTD respectively

#### Scenario: Architecture constants map
- **WHEN** `x86` or `arm` is mapped
- **THEN** it maps to hcloud.ArchitectureX86 or hcloud.ArchitectureARM respectively

### Requirement: No public skip-cleanup escape hatch
The provider schema SHALL NOT expose a way to skip the library's resource cleanup. Any
debug skip SHALL be gated behind an environment variable only.

#### Scenario: Schema has no skip-cleanup attribute
- **WHEN** the resource schema is inspected
- **THEN** it contains no attribute that disables cleanup

