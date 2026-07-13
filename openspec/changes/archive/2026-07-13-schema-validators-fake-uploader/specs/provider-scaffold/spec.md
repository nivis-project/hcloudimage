## ADDED Requirements

### Requirement: Provider injects an Uploader and registers the data source
The provider SHALL construct an `Uploader` implementation and pass it to resources and
data sources via configure data, and SHALL register the `hcloudimage_snapshot` data
source alongside the `hcloudimage_image` resource.

#### Scenario: Uploader reaches the resource
- **WHEN** the provider's `Configure` runs
- **THEN** an `Uploader` is available to the resource and data source through configure data

#### Scenario: Data source is registered
- **WHEN** the provider reports its data sources
- **THEN** `hcloudimage_snapshot` is present
