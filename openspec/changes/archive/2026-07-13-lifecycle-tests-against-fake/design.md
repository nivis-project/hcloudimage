## Context

`terraform-plugin-testing`'s `resource.Test` drives a real `terraform` binary against an
in-process instance of the provider over the plugin protocol. The devShell already
provides `terraform` and `opentofu`; the harness finds a binary automatically (or via
`TF_ACC_TERRAFORM_PATH`). No `TF_ACC=1` gating is needed — these are protocol tests, not
billable acceptance tests, so they should run in the normal `go test` pass.

## Decisions

### Injecting an inspectable fake
- Use the `NewWithUploader(version, uploader)` seam added in milestone 02 so the test holds
  the **same** `*FakeUploader` the provider uses. This lets a test both apply HCL and then
  assert `fake.UploadCalls` / `fake.DeleteCalls` / `fake.UpdateCalls`.
- `protoV6ProviderFactories` maps `"hcloudimage"` to that provider instance via
  `providerserver.NewProtocol6WithError`.

### Lifecycle cases
1. **Create**: apply a minimal `image_url` config → `TestCheckResourceAttrSet(id)`,
   `effective_labels` contains the created-by label, exactly one `UploadCall`.
2. **ForceNew replace**: `image_path`+`image_sha256` config; step 2 changes `image_sha256`
   → assert `plancheck` reports a replace, and after apply a second `UploadCall` happened
   and the old id was deleted.
3. **In-place update**: change `labels`/`description` between steps → `plancheck` reports
   an update (not replace); assert an `UpdateCall`, and **no** extra `UploadCall`.
4. **Out-of-band deletion**: after create, call `fake.MarkDeleted(id)` in a `PreConfig`,
   then expect a non-empty plan / the resource to be recreated (state removed on refresh).
5. **Destroy**: framework runs destroy at the end; `CheckDestroy` asserts a `DeleteCall`.

### Plan assertions
- Use `plancheck.ExpectResourcePlanAction` (or `ExpectResourceAction` with
  `plancheck.ResourceActionReplace` / `.Update`) in `ConfigPlanChecks.PreApply` to prove
  ForceNew vs in-place without relying on side effects alone.

### Terraform vs tofu
- The lifecycle harness uses whichever binary it discovers; running the whole suite under
  both `terraform` and `tofu` is the hermetic VM test's job (milestone 05). Here we assert
  provider behaviour once through the protocol.

## Risks / Trade-offs

- `image_path` cases need a real local file for `filesha256`-style flows; use a tiny temp
  file created by the test and a fixed sha string (the fake does not verify the hash).
- If the framework can't find a terraform binary in some environments, document
  `TF_ACC_TERRAFORM_PATH`; in the devShell/CI it is always present.
