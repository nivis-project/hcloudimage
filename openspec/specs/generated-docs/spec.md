# generated-docs Specification

## Purpose
TBD - created by archiving change ci-pipeline. Update Purpose after archive.
## Requirements
### Requirement: Provider docs are generated and committed
The repository SHALL contain `tfplugindocs`-generated documentation for the provider,
resource, and data source, committed under `docs/`, derived from schema
`MarkdownDescription`s and the `examples/`.

#### Scenario: Docs exist for each surface
- **WHEN** `docs/` is inspected
- **THEN** it contains `index.md`, `resources/image.md`, and `data-sources/snapshot.md`

#### Scenario: Docs regenerate deterministically
- **WHEN** `tfplugindocs generate` is run against the committed sources
- **THEN** it produces no diff against the committed `docs/`

