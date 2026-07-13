## Context

This is the one milestone whose *verification* needs paid cloud + secrets. The strategy:
build everything reproducibly and make it compile and gate correctly, so a human with a
budget-limited Hetzner project can run it unchanged (BRIEFING.md §14, §15).

## Decisions

### Alpine fixture (nix/test-image.nix)
- Input: the Alpine `nocloud`/generic-cloud raw image for amd64 and aarch64, pinned by URL
  + sha256 as a `fetchurl` (fixed-output → hermetic). Alpine publishes
  `alpine-virt`/`nocloud_alpine-*-x86_64.raw` style artifacts; pin exact version + hash.
- Bake step: mount the raw image (offline, via `libguestfs`/`guestfish` in the derivation,
  or loop-mount in a fixed-output builder), write a **throwaway** SSH public key to
  `/root/.ssh/authorized_keys`, enable `sshd` (rc-update add sshd) and DHCP
  (`/etc/network/interfaces`), then recompress to `.raw.xz` (`xz -T0`).
- The throwaway keypair lives in `test/fixtures/` (public key baked in; private key used by
  the acceptance test to SSH in). It is explicitly disposable and documented as such.
- Determinism: `SOURCE_DATE_EPOCH`, sorted tar, fixed xz preset; assert two builds match.

### aarch64 build path
- Support all three per BRIEFING.md §7: a native aarch64 runner, a configured remote
  builder, or `boot.binfmt`/QEMU emulation. Document that CI uses QEMU/binfmt by default
  (works on the hosted x86 runner) and that a native/remote builder is faster. The arm
  acceptance job must not silently fall back to an x86 image — assert the fixture arch.

### Acceptance tests (internal/provider/*_acc_test.go)
- `resource.Test` gated by `TF_ACC=1`; skip with a clear message when `HCLOUD_TOKEN` is
  unset so `go test ./...` stays green locally/CI without secrets.
- Config composes `hcloudimage_image.test` (uploads the fixture) + `hcloud_server.test`
  (`hetznercloud/hcloud`, pinned `~> 1.48`) booting from `.id`.
- Reachability assertion: from the runner, `ssh -i <throwaway-key> root@<ip> cat
  /etc/os-release` and check it says Alpine — outbound only, no inbound-to-runner infra.
- Matrix: x86 (cx22 + amd64 fixture), arm (cax11 + aarch64 fixture). arm behind a toggle.
- `CheckDestroy` + a `defer`/`t.Cleanup` that deletes leftover snapshots even on failure.

### Workflows
- `acceptance.yml`: `on: workflow_dispatch` (with `run_arm` boolean input), `push` to main,
  `schedule` nightly. Guard `if: github.event.pull_request.head.repo.fork == false` — never
  fork PRs. `concurrency: { group: acceptance, cancel-in-progress: false }` (one at a
  time). Sets `HCLOUD_TOKEN` from secrets. Always-run cleanup step (`if: always()`).
- `cleanup.yml`: nightly `schedule`, runs `hcloud-upload-image cleanup` scoped to the
  library's label so only orphans from crashed runs are swept.

### Cost/safety (mandatory, §8.3/§15)
- Cheapest server types (cx22/cax11), smallest fixture, short timeouts, pinned hcloud
  version, guaranteed cleanup, nightly sweep, never fork PRs, concurrency-limited, arm
  toggle default-off nightly but on before release.

## Risks / Trade-offs

- Loop-mounting/guestfish in a Nix derivation may need extra privileges; if the offline
  bake proves too fiddly, BRIEFING.md §8.3 permits a minimal NixOS fixture fallback — keep
  the derivation interface identical so the acceptance test is unaffected.
- The live acceptance run is intentionally *not* executed here (no token, costs money). The
  deliverable is: reproducible fixtures, compiling+gated tests, and correct workflows.
