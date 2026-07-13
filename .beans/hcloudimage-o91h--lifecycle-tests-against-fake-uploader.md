---
# hcloudimage-o91h
title: Lifecycle tests against fake uploader
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:58:55Z
parent: hcloudimage-541n
blocked_by:
    - hcloudimage-jzim
---

terraform-plugin-testing suite with fake injected: apply creates state with synthetic id; ForceNew attrs fire RequiresReplace; description/labels update in place; Read removes missing snapshot from state; destroy deletes. Coverage reported. (BRIEFING.md §8.1)

## Summary of Changes

Added terraform-plugin-testing lifecycle suite driving the provider (fake uploader injected via NewWithUploader) through the real Terraform protocol: create populates id + effective_labels; ForceNew (image_sha256) plancheck expects Replace with re-upload + old-snapshot delete; labels/description change plancheck expects in-place Update with no re-upload; out-of-band deletion removes the resource from state on refresh; destroy deletes.

The lifecycle tests surfaced and fixed three real provider bugs: (1) Computed+ForceNew attributes (compression/format/location) needed UseStateForUnknown or they went "known after apply" and spuriously forced replacement on every refresh; (2) a custom effective_labels plan modifier now computes the exact merged label set at plan time, fixing the in-place update inconsistent-result error and giving precise plans; (3) Read preserves a null description instead of coercing to "". Coverage of internal/provider rose from ~41% to ~74%. Upgraded terraform-plugin-framework to v1.19.0 for terraform-plugin-go v0.31.0 compatibility. buildGoModule uses doCheck=false (protocol tests need a terraform binary + TF_ACC, provided in devShell/CI and the hermetic VM test).
