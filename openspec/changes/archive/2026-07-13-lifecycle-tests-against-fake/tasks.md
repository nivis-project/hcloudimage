## 1. Test harness

- [x] 1.1 Add `github.com/hashicorp/terraform-plugin-testing` as a dependency; update vendorHash
- [x] 1.2 Helper: `protoV6ProviderFactories` wiring `NewWithUploader` so the test shares the fake instance

## 2. Lifecycle cases

- [x] 2.1 Create: apply `image_url` config; assert `id` set, `effective_labels` has created-by, one UploadCall
- [x] 2.2 ForceNew replace: `image_path`+`image_sha256`; change sha; plancheck expects Replace; second UploadCall + old id deleted
- [x] 2.3 In-place update: change labels/description; plancheck expects Update; UpdateCall, no extra UploadCall
- [x] 2.4 Out-of-band deletion: `MarkDeleted` in PreConfig; expect non-empty plan / recreate
- [x] 2.5 Destroy: `CheckDestroy` asserts a DeleteCall

## 3. Close out

- [x] 3.1 `go test ./... -cover` green (coverage of resource CRUD materially up); `nix flake check`
- [x] 3.2 Archive OpenSpec change; complete beans epic + milestone 03; commit with jj (Pim Snel, no self-promotion)
