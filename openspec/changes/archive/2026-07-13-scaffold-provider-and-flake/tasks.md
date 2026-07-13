## 1. Go module

- [x] 1.1 `go mod init github.com/nivis-project/terraform-provider-hcloudimage`; pin latest stable Go
- [x] 1.2 Add `github.com/hashicorp/terraform-plugin-framework` dependency
- [x] 1.3 `go mod tidy` and commit `go.sum`

## 2. Provider server

- [x] 2.1 `main.go`: run `providerserver.Serve` with address `registry.terraform.io/nivis-project/hcloudimage` (protocol v6)
- [x] 2.2 `internal/provider/provider.go`: implement `provider.Provider` with `Metadata`, `Schema`, `Configure`
- [x] 2.3 Provider config schema per §3.1: `token` (sensitive, `HCLOUD_TOKEN` fallback), `endpoint`, `poll_interval`
- [x] 2.4 `internal/provider/image_resource.go`: empty `hcloudimage_image` resource that compiles and is registered in `provider.Resources`

## 3. Nix build

- [x] 3.1 Add `packages.default` to `flake.nix` via `buildGoModule` (reuse `forAllSystems`, no flake-utils)
- [x] 3.2 Compute and pin the real `vendorHash` (start from `lib.fakeHash`, read the "got:" hash)
- [x] 3.3 Verify `nix develop`, `nix build .#default`, and `nix flake check` all succeed

## 4. Dev entrypoints

- [x] 4.1 `GNUmakefile` with `build`, `test`, `lint`, `fmt` targets mirroring the flake devShell
- [x] 4.2 Confirm gate: `go build ./...` passes and `nix develop` works

## 5. Close out

- [x] 5.1 Run unit build/test; archive this OpenSpec change (`/opsx:archive`)
- [x] 5.2 Mark the beans epic "Go module & provider server skeleton" and "Nix flake foundation" progress; commit with jj (author Pim Snel, no self-promotion)
