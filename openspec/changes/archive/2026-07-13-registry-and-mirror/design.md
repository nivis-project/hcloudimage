## Context

The mirror layout and the `packages.provider-mirror` output already exist and are exercised
inside the hermetic VM test. This change surfaces them for human/Nivis use and verifies the
consumer path directly (outside the VM), then documents the registry onboarding that only a
human with the GPG key and org access can complete.

## Decisions

### Mirror consumption verification
- A `scripts/verify-mirror.sh` that: `nix build .#provider-mirror`, writes a temp CLI
  config with a `filesystem_mirror { path = <result>; include = [...hcloudimage] }` + a
  `direct { exclude = [...] }`, drops a minimal consumer `main.tf` using
  `hcloudimage_image` with a pinned `version`, and runs `init -backend=false` + `plan`
  (with `HCLOUDIMAGE_FAKE=1`, no token) under both `terraform` and `tofu`. Success proves
  registry-less consumption end to end.
- This differs from `validate-examples.sh` (which validates the shipped examples): this one
  specifically proves the **mirror install path** with a pinned version, the way Nivis will.
- A `make consume-mirror` target wraps it.

### Docs
- `docs/consuming-from-nix-mirror.md`: the filesystem_mirror recipe (preferred, works for
  `init`) and the `dev_overrides` recipe (tight iteration, skips `init`), with the exact
  paths from `nix build .#provider-mirror` / `.#default`. Note the tofu vs terraform
  `init`/`dev_overrides` difference discovered in milestone 04.
- `docs/publishing.md`: the §14 human steps — Terraform Registry (connect the GitHub repo
  and the GPG public key to the `nivis-project` namespace; it then picks up tagged releases
  from `release.yml`); OpenTofu Registry (open a submission PR to `opentofu/registry`
  referencing the repo + GPG key; confirm the current procedure at submission time). List
  the prerequisite secrets/keys.

## Risks / Trade-offs

- `tofu init` does not honour `dev_overrides` the way terraform does (it still tries to
  resolve the provider), so the **filesystem_mirror** recipe is the portable one for both
  tools — the docs lead with it and mark `dev_overrides` as terraform-friendly iteration.
- Registry submission itself can't be automated here (org access + GPG key + external PR);
  it is documented as the human step, consistent with BRIEFING §14.
