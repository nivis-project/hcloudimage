---
# hcloudimage-1uf7
title: Acceptance tests + cleanup workflow
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:04:22Z
updated_at: 2026-07-13T18:06:10Z
parent: hcloudimage-r4gl
blocked_by:
    - hcloudimage-lofb
---

TF_ACC=1 tests composing hcloudimage_image + official hcloud_server; SSH-reachability assertion from runner with baked key; ForceNew + in-place semantics against real snapshots; both x86 and arm (arm toggle-gated); cost controls; acceptance.yml (dispatch/push-main/nightly, never fork PRs) + cleanup.yml orphan sweep. (BRIEFING.md §8.3, §9)

## Summary of Changes

Added internal/provider/image_acc_test.go: TF_ACC-gated acceptance tests (skip cleanly without HCLOUD_TOKEN) composing hcloudimage_image + official hcloud_server (pinned ~> 1.48), booting from the snapshot id and asserting real guest reachability by SSHing from the runner with the baked throwaway key to read /etc/os-release (retry loop; not just running). Covers x86 (cx22) and arm (cax11, HCLOUDIMAGE_ACC_RUN_ARM-gated). Added .github/workflows/acceptance.yml (workflow_dispatch + push-main + nightly; never fork PRs; concurrency-limited; arm behind run_arm input; KVM enablement to build the fixture; always-run cleanup) and cleanup.yml (nightly label-scoped hcloud-upload-image cleanup). README documents running acceptance, cost/safety controls, and the aarch64 build path. Live run is the documented human step (needs a budget-limited HCLOUD_TOKEN).


## Live verification (2026-07-13)

TestAccImage_RealHetzner_x86 PASSED against a real Hetzner project: image uploaded, snapshot created, hcloud_server booted from it, and SSH-into-guest reachability confirmed. Shakeout fixed: description Optional+Computed (real provider bug), UEFI (not BIOS) fixture, authorized_keys root ownership + sshd root-login policy + unlocked root, cloud-init disabled, dropped hcloud_ssh_key (fingerprint uniqueness), env-overridable server_type/location, absolute image-path guard against empty glob.
