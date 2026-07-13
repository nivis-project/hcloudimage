## ADDED Requirements

### Requirement: Image source is exactly one of URL or path
The `hcloudimage_image` resource SHALL accept exactly one of `image_url` or `image_path`,
enforced by a config validator (not documentation alone).

#### Scenario: Both sources set is rejected
- **WHEN** a config sets both `image_url` and `image_path`
- **THEN** validation fails with an error naming the conflict

#### Scenario: Neither source set is rejected
- **WHEN** a config sets neither `image_url` nor `image_path`
- **THEN** validation fails

### Requirement: image_sha256 required with local path
The resource SHALL require `image_sha256` if and only if `image_path` is set, enforced by
a config validator.

#### Scenario: Path without sha is rejected
- **WHEN** `image_path` is set but `image_sha256` is not
- **THEN** validation fails

#### Scenario: URL without sha is allowed
- **WHEN** `image_url` is set and `image_sha256` is not
- **THEN** validation passes

### Requirement: Enumerated attributes are validated
The resource SHALL validate that `architecture` is `x86` or `arm`, `compression` is one of
`none`/`bz2`/`xz`/`zstd`, and `format` is one of `raw`/`qcow2`.

#### Scenario: Invalid architecture rejected
- **WHEN** `architecture` is set to `sparc`
- **THEN** validation fails

### Requirement: Label values obey Hetzner rules
The resource SHALL reject label values containing `/` (a Hetzner label constraint).

#### Scenario: Slash in label value rejected
- **WHEN** a label value contains `/`
- **THEN** validation fails

### Requirement: ForceNew attributes trigger replacement
The resource SHALL require replacement (a new upload producing a new snapshot) when any
ForceNew attribute changes: `image_url`, `image_path`, `image_sha256`, `architecture`,
`compression`, `format`, `server_type`, `location`, or `image_size`.

#### Scenario: Changing sha forces replacement
- **WHEN** `image_sha256` changes in a subsequent plan
- **THEN** the plan shows the resource must be replaced

### Requirement: Description and labels update in place
Changing `description` or `labels` SHALL update the existing snapshot in place without a
new upload.

#### Scenario: Changing labels does not replace
- **WHEN** `labels` changes in a subsequent plan
- **THEN** the plan shows an in-place update, not a replacement

### Requirement: effective_labels includes library defaults
The computed `effective_labels` SHALL contain the user labels merged with the library
default `apricote.de/created-by=hcloud-upload-image`.

#### Scenario: Created-by label surfaced
- **WHEN** the resource is created with `labels = { os = "nixos" }`
- **THEN** `effective_labels` contains both `os=nixos` and `apricote.de/created-by=hcloud-upload-image`

### Requirement: Read removes a missing snapshot from state
On `Read`, if the snapshot no longer exists, the resource SHALL be removed from state
without raising an error.

#### Scenario: Out-of-band deletion clears state
- **WHEN** the snapshot backing a resource is deleted out of band
- **AND** a refresh/read runs
- **THEN** the resource is removed from state and no error is raised
