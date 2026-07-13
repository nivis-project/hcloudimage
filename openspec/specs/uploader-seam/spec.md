# uploader-seam Specification

## Purpose
TBD - created by archiving change schema-validators-fake-uploader. Update Purpose after archive.
## Requirements
### Requirement: Resource depends only on the Uploader interface
The `hcloudimage_image` resource and `hcloudimage_snapshot` data source SHALL depend only
on the `Uploader` interface, never on the concrete hcloud library, so implementations can
be swapped for tests.

#### Scenario: Fake uploader drives the resource
- **WHEN** the provider is wired with the in-memory fake uploader
- **THEN** create/read/update/delete operate against the fake with no network access

### Requirement: Config maps to UploadRequest
The resource SHALL map its configuration to an `UploadRequest` covering both sources
(url/path), all compressions (`none`/`bz2`/`xz`/`zstd`), and both architectures
(`x86`/`arm`).

#### Scenario: Every compression maps
- **WHEN** a resource is configured with each supported compression value in turn
- **THEN** the resulting `UploadRequest.Compression` matches the configured value

### Requirement: Fake simulates out-of-band deletion
The fake uploader SHALL be able to report a snapshot as no longer existing, so the
resource's state-removal-on-read behaviour is testable.

#### Scenario: Get after out-of-band delete returns not found
- **WHEN** a snapshot is marked deleted out of band in the fake
- **THEN** `Get` for that id returns a not-found result

