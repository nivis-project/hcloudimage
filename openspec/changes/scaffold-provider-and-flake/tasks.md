## 1. Go module

- [ ] 1.1 `go mod init github.com/nivis-project/terraform-provider-hcloudimage`; pin latest stable Go
- [ ] 1.2 Add `github.com/hashicorp/terraform-plugin-framework` dependency
- [ ] 1.3 `go mod tidy` and commit `go.sum`

## 2. Provider server

- [ ] 2.1 `main.go`: run `providerserver.Serve` with address `registry.terraform.io/nivis-project/hcloudimage` (protocol v6)
- [ ] 2.2 `internal/provider/provider.go`: implement `provider.Provider` with `Metadata`, `Schema`, `Configure`
- [ ] 2.3 Provider config schema per §3.1: `token` (sensitive, `HCLOUD_TOKEN` fallback), `endpoint`, `poll_interval`
- [ ] 2.4 `internal/provider/image_resource.go`: empty `hcloudimage_image` resource that compiles and is registered in `provider.Resources`

## 3. Nix build

- [ ] 3.1 Add `packages.default` to `flake.nix` via `buildGoModule` (reuse `forAllSystems`, no flake-utils)
- [ ] 3.2 Compute and pin the real `vendorHash` (start from `lib.fakeHash`, read the "got:" hash)
- [ ] 3.3 Verify `nix develop`, `nix build .#default`, and `nix flake check` all succeed

## 4. Dev entrypoints

- [ ] 4.1 `GNUmakefile` with `build`, `test`, `lint`, `fmt` targets mirroring the flake devShell
- [ ] 4.2 Confirm gate: `go build ./...` passes and `nix develop` works

## 5. Close out

- [ ] 5.1 Run unit build/test; archive this OpenSpec change (`/opsx:archive`)
- [ ] 5.2 Mark the beans epic "Go module & provider server skeleton" and "Nix flake foundation" progress; commit with jj (author Pim Snel, no self-promotion)
