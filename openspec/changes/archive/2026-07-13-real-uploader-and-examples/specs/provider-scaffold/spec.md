## ADDED Requirements

### Requirement: Uploader selection prefers real, falls back to fake
The provider SHALL select the real uploader when a token is present and `HCLOUDIMAGE_FAKE`
is unset, force the fake when `HCLOUDIMAGE_FAKE=1`, and otherwise fall back to the fake so
that validation, planning without a token, and tests need no real credentials.

#### Scenario: No token falls back to fake
- **WHEN** the provider is configured with no token and `HCLOUDIMAGE_FAKE` unset
- **THEN** the fake uploader is used, so `terraform validate`/plan needs no credentials

#### Scenario: Fake forced by environment
- **WHEN** `HCLOUDIMAGE_FAKE=1`
- **THEN** the fake uploader is used even if a token is present
