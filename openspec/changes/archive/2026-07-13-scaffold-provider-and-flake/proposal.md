## Why

Milestone 01 (BRIEFING.md §13.1) needs a compiling foundation before any provider
behaviour can be built: a Go module wired to `terraform-plugin-framework` (protocol v6),
a provider server entrypoint, an empty registered resource, and a hermetic Nix flake so
every later milestone builds and tests reproducibly. Without this scaffold there is
nothing to iterate on.

## What Changes

- Add `go.mod` / `go.sum` pinning the latest stable Go and `terraform-plugin-framework`.
- Add `main.go` running the provider via `providerserver.Serve` (protocol v6).
- Add `internal/provider/provider.go` implementing `provider.Provider` with the
  provider-config schema from BRIEFING.md §3.1 (`token` sensitive + `HCLOUD_TOKEN`
  fallback, `endpoint`, `poll_interval`) and registering a placeholder resource.
- Add an empty `internal/provider/image_resource.go` (`hcloudimage_image`) that compiles
  and registers but implements no behaviour yet (fleshed out in milestone 02).
- Extend `flake.nix` with `packages.default` via `buildGoModule` (pinned `vendorHash`),
  alongside the existing `devShells.default`.
- Add a `GNUmakefile` mirroring the flake dev entrypoints (BRIEFING.md §6).

## Capabilities

### New Capabilities
- `provider-scaffold`: the provider server, module structure, provider-config schema, and
  the hermetic Nix build (devShell + buildGoModule) that all later work depends on.

### Modified Capabilities
<!-- none: this is the first change -->

## Impact

- New: `main.go`, `go.mod`, `go.sum`, `internal/provider/`, `GNUmakefile`.
- Modified: `flake.nix` (adds `packages.default`), `flake.lock`.
- Gate: `go build ./...` succeeds and `nix develop` / `nix build` work.
