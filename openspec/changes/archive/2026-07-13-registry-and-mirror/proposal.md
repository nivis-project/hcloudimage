## Why

Milestone 09 (BRIEFING.md §13.9, §7, §10) is the last mile: publishing to the Terraform
and OpenTofu registries and, crucially for Nivis, consuming the provider from a Nix-built
filesystem mirror **without** any public registry. The `packages.provider-mirror` output
already exists (added in milestone 05 for the hermetic test); this change documents its
end-to-end consumption, verifies it works with a real config under both `terraform` and
`tofu`, and captures the registry onboarding as documented human steps (§14).

## What Changes

- A `docs/consuming-from-nix-mirror.md` (or README section) with the exact `dev_overrides`
  and `filesystem_mirror` CLI-config recipes Nivis uses to consume the provider from
  `nix build .#provider-mirror`, plus the tight-iteration `dev_overrides` alternative.
- A verification script/target that stands up a real consumer config against the mirror and
  runs `init` + `plan` under both `terraform` and `tofu`, proving registry-less consumption
  end to end (this reuses the mirror the hermetic test already exercises).
- A `docs/publishing.md` capturing the human onboarding steps: Terraform Registry (connect
  repo + GPG public key to the `nivis-project` namespace; tagged releases picked up
  automatically) and OpenTofu Registry (submit to `opentofu/registry` referencing the repo
  + GPG key), with the §14 prerequisites listed.

## Capabilities

### New Capabilities
- `nix-mirror-consumption`: documented + verified registry-less consumption of the provider
  via the Nix filesystem mirror.
- `registry-publishing`: documented Terraform/OpenTofu registry onboarding (human steps).

### Modified Capabilities
<!-- none: provider-mirror package already exists; this documents + verifies it -->

## Impact

- New: `docs/consuming-from-nix-mirror.md`, `docs/publishing.md`, a mirror-consumption
  verification (script + make target).
- Gate: a real consumer config resolves the provider from the mirror and `init` + `plan`
  succeed under both `terraform` and `tofu`; onboarding steps are documented. Actual
  registry submission is a human step (§14).
