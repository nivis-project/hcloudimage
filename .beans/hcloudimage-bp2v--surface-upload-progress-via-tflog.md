---
# hcloudimage-bp2v
title: Surface upload progress via tflog
status: todo
type: feature
created_at: 2026-07-13T18:04:17Z
updated_at: 2026-07-13T18:04:17Z
parent: hcloudimage-2dsy
---

Surface upload progress (e.g. percentage) from the real uploader via tflog, so long multi-minute uploads report progress under TF_LOG_PROVIDER.

## Context / findings

- A Terraform/OpenTofu provider is a headless gRPC plugin: go-plugin captures the provider process stdout/stderr and routes it into the CLI logging system. Writing to os.Stderr does NOT stream a live progress bar to the user; it only appears as log lines when TF_LOG / TF_LOG_PROVIDER is set. There is no provider->terminal live-progress API by design (this is why terraform itself shows no upload bar).
- The sanctioned output channel is tflog (github.com/hashicorp/terraform-plugin-log/tflog), already an indirect dependency via terraform-plugin-framework.
- hcloud-upload-image/v2 emits progress through a context/slog logger. The clean approach: bridge that library logger into tflog inside uploader_hcloud.go so upload progress lands in the provider log.
- Terraform also prints its own "Still creating... [Nm elapsed]" heartbeat during a slow Create for free; the tflog percentage is the richer signal on top.

## Scope

- [ ] In newHcloudUploader / Upload, pass a slog.Logger (or the libs contextlogger) into the hcloudimages client that forwards records to tflog on the ctx.
- [ ] Emit at INFO (or DEBUG) with structured fields (e.g. percent, bytes).
- [ ] Document that users see progress with TF_LOG_PROVIDER=INFO terraform apply (add to README testing/troubleshooting).
- [ ] Verify the bridge emits during a real upload (acceptance) and/or a unit test asserting the adapter forwards records.

## Notes

- Do not attempt a live terminal progress bar — not possible within the plugin protocol.
- Keep the fake uploader unaffected (no real upload, no progress).
