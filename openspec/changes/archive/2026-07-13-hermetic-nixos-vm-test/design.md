## Context

`pkgs.testers.runNixOSTest` boots a real NixOS VM under QEMU, runs a Python driver script,
and fails the build if any assertion fails. It is Linux-only, so the check is added only
for `x86_64-linux`/`aarch64-linux` via the `forAllSystems`/`lib.optionalAttrs` pattern
(still no flake-utils).

## Decisions

### Provider mirror inside the VM
- Build a filesystem-mirror derivation from `packages.default`:
  `<mirror>/registry.terraform.io/nivis-project/hcloudimage/0.1.0/<os>_<arch>/terraform-provider-hcloudimage_v0.1.0`
  and the same under `registry.opentofu.org`. This is the same layout milestone 09 ships as
  `packages.provider-mirror`; factor it into a shared helper so 05 and 09 don't diverge.
- A CLI config file (also in the store) sets `provider_installation { filesystem_mirror {
  path = <mirror>; include = [...hcloudimage] } direct { exclude = [...hcloudimage] } }`.
- The VM has no network; `direct{}` is never exercised because the test HCL uses only our
  provider (no hcloud), so `init` resolves entirely from the mirror. This keeps the check
  hermetic. (The hcloud-composing example is validated separately in milestone 04.)

### Test HCL
- Minimal config using only `hcloudimage_image` with `image_url` (no local file needed).
  Variables drive `image_sha256`-like ForceNew and `labels` changes across steps. Because
  the config uses `image_url`, `image_sha256` isn't the trigger — use `architecture` or
  `location` change, or better: use `image_path` pointing at a tiny file created in the VM
  so `image_sha256` is the realistic ForceNew trigger (matches §8.2 wording).
- Decision: use `image_path` with a small file written in the VM and a var for
  `image_sha256`, so the test mirrors the documented ForceNew-on-sha behaviour.

### Driver script (Python)
For each binary in [terraform, tofu]:
1. Fresh temp dir, copy the test HCL, `export HCLOUDIMAGE_FAKE=1`.
2. `init` (offline, from mirror).
3. `apply -auto-approve` with sha=A, labels={a}; assert state has one resource with an id.
4. `plan` with sha=B; assert the plan reports a replace (grep `must be replaced` / `-/+`).
   `apply`; assert still one resource (new id).
5. `plan` with labels={b} (sha=B); assert in-place (`~ update` / `1 to change`, `0 to
   destroy`, `0 to add`). `apply`.
6. `destroy -auto-approve`; assert empty state.

Assertions use `terraform show -json` / `plan -json` where practical, else string checks on
human output; keep them robust across both binaries.

### Wiring
- `checks.<linux-system>.hermetic-e2e = runNixOSTest { ... }`.
- `nix flake check` picks it up automatically.

## Risks / Trade-offs

- VM tests are slow (minutes) and need KVM; that's inherent to the DoD gate and acceptable.
- `terraform`/`tofu` plan JSON differs subtly; prefer `-json` + parse `resource_changes[].
  change.actions` (`["delete","create"]` = replace, `["update"]` = in-place) which both
  tools emit identically.
- aarch64-linux VM needs an aarch64 builder; the check is defined for both Linux systems but
  will only run where a builder exists (documented).
