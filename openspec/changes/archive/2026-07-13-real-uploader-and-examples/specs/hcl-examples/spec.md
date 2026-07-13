## ADDED Requirements

### Requirement: Runnable examples validate under Terraform and OpenTofu
The repository SHALL contain runnable HCL examples for the provider, the resource, and the
data source, and each SHALL pass both `terraform validate` and `tofu validate`.

#### Scenario: Resource example validates
- **WHEN** `terraform validate` and `tofu validate` run on the resource example
- **THEN** both succeed

#### Scenario: Example composes the official hcloud provider
- **WHEN** the resource example is read
- **THEN** it boots an `hcloud_server` from `hcloudimage_image.<name>.id`, proving the int64
  id flows directly into `hcloud_server.image`
