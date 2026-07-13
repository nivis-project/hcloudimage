## ADDED Requirements

### Requirement: Apply creates state with a synthetic id
The resource SHALL, when applied against the fake uploader, produce state whose `id` is
set and whose `effective_labels` includes the library's created-by label.

#### Scenario: Create populates id and effective labels
- **WHEN** a config with `image_url` and `architecture = "x86"` is applied
- **THEN** the resulting state has a non-empty `id`
- **AND** `effective_labels` contains `apricote.de/created-by = hcloud-upload-image`

### Requirement: ForceNew attribute change replaces the resource
The resource SHALL be replaced (destroy + recreate) when a ForceNew attribute changes
between plans.

#### Scenario: Changing image_sha256 replaces
- **WHEN** an applied `image_path` resource has its `image_sha256` changed
- **THEN** the plan reports a replace action
- **AND** after apply a new upload occurred and the previous snapshot was deleted

### Requirement: Metadata change updates in place
The resource SHALL update `labels` and `description` in place, with no new upload.

#### Scenario: Changing labels updates without re-upload
- **WHEN** an applied resource has its `labels` changed
- **THEN** the plan reports an update action, not a replace
- **AND** after apply no additional upload occurred

### Requirement: Out-of-band deletion is reconciled
The resource SHALL be removed from state on refresh when its snapshot no longer exists.

#### Scenario: Refresh after out-of-band delete
- **WHEN** the snapshot backing an applied resource is deleted out of band
- **AND** Terraform refreshes
- **THEN** the resource is no longer present in state (it will be recreated on the next apply)

### Requirement: Destroy deletes the snapshot
On destroy, the resource SHALL delete its snapshot.

#### Scenario: Destroy calls delete
- **WHEN** an applied resource is destroyed
- **THEN** the snapshot's delete path is invoked
