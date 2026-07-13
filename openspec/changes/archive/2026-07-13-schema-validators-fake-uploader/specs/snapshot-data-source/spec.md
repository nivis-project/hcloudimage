## ADDED Requirements

### Requirement: Snapshot lookup by id or selector
The `hcloudimage_snapshot` data source SHALL accept exactly one of `id` or
`with_selector`, and SHALL populate computed fields (`name`, `description`, `labels`,
`architecture`, `created`) from the resolved snapshot.

#### Scenario: Lookup by id
- **WHEN** the data source is configured with a known `id`
- **THEN** its computed fields reflect that snapshot

#### Scenario: Both id and selector rejected
- **WHEN** both `id` and `with_selector` are set
- **THEN** validation fails

### Requirement: Selector ambiguity handling
When `with_selector` matches more than one snapshot, the data source SHALL error unless
`most_recent = true`, in which case it SHALL pick the newest.

#### Scenario: Ambiguous selector without most_recent errors
- **WHEN** a selector matches multiple snapshots and `most_recent` is not `true`
- **THEN** the lookup fails with an ambiguity error

#### Scenario: Ambiguous selector with most_recent picks newest
- **WHEN** a selector matches multiple snapshots and `most_recent = true`
- **THEN** the newest snapshot is returned
