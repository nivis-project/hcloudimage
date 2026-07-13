## Context

First change in the project. Establishes the module layout of BRIEFING.md §6 and the Nix
build contract of §7. Everything downstream (schema, validators, fake uploader, hermetic
test) attaches to this skeleton.

## Decisions

### Module path and framework
- Module path: `github.com/nivis-project/terraform-provider-hcloudimage`.
- Use `terraform-plugin-framework` (protocol v6). **Not** SDKv2 (locked, BRIEFING.md §1).
- Pin the latest stable Go in both `go.mod` and the flake so CI, local, and hermetic
  builds agree.

### Provider server
- `main.go` calls `providerserver.Serve(ctx, provider.New(version), providerserver.ServeOpts{Address: "registry.terraform.io/nivis-project/hcloudimage"})`.
- `provider.New` returns a `func() provider.Provider` capturing the build version.

### Provider configuration schema (BRIEFING.md §3.1)
- `token` — string, sensitive, optional; falls back to `HCLOUD_TOKEN` env var when unset.
- `endpoint` — string, optional; overrides the hcloud API endpoint (testing/mock).
- `poll_interval` — string (duration), optional; passthrough for action polling.
- Resolve config into a small struct stored on `resp.ResourceData` / `resp.DataSourceData`
  so resources/data-sources can read it. The real client is wired in milestone 04 behind
  the `Uploader` interface — the scaffold just plumbs config through.

### Placeholder resource
- `hcloudimage_image` registered so `provider.Resources` is non-empty and the server
  starts, but with a minimal/empty schema and no CRUD behaviour. Milestone 02 replaces it
  with the full §3.2 schema. This keeps milestone 01's gate purely "compiles + builds".

### Nix build
- `packages.default = pkgs.buildGoModule { ... vendorHash = "..."; }` per system via the
  existing `forAllSystems` helper (plain nix, **no flake-utils** — BRIEFING.md §7).
- Obtain `vendorHash` by first setting `vendorHash = lib.fakeHash`, running
  `nix build .#default`, and copying the "got:" hash from the error. Commit the real hash.
- `nix flake check` must stay green; `flake.lock` stays committed.

### GNUmakefile
- Thin targets (`build`, `test`, `lint`, `fmt`) that mirror what the flake devShell
  provides, so contributors not using Nix still have an entrypoint (BRIEFING.md §6).

## Risks / Trade-offs

- `vendorHash` drift: any dependency change requires recomputing it. Documented in the
  Makefile/README so it is not a surprise in later milestones.
- Keeping the resource empty now means milestone 02 does a larger edit — acceptable, it
  keeps each milestone's gate crisp and independently verifiable.
