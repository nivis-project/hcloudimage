---
# hcloudimage-yvqy
title: Go module & provider server skeleton
status: completed
type: epic
priority: normal
created_at: 2026-07-13T14:03:59Z
updated_at: 2026-07-13T14:37:17Z
parent: hcloudimage-c2wx
---

main.go with providerserver, go.mod (terraform-plugin-framework, protocol v6), internal/provider/provider.go with provider config schema (token/endpoint/poll_interval), empty resource registered. Gate: go build. Tracked via an OpenSpec change; tasks live in its tasks.md.

OpenSpec change: scaffold-provider-and-flake

## Summary of Changes

Scaffolded the Go module (terraform-plugin-framework, protocol v6), main.go provider server, provider config schema per BRIEFING §3.1 (token/endpoint/poll_interval with HCLOUD_TOKEN fallback), and an empty registered hcloudimage_image resource. Unit tests cover metadata, schema, and resource registration. go build and go test pass.
