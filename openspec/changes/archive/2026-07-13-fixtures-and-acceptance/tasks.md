## 1. Throwaway SSH key + Alpine fixture

- [x] 1.1 `test/fixtures/`: throwaway SSH keypair (public baked in; private used by the acceptance test), documented as disposable
- [x] 1.2 `nix/test-image.nix`: reproducible derivation — pinned Alpine raw (amd64/aarch64) via fetchurl, bake authorized_keys, enable sshd + DHCP, recompress to .raw.xz (SOURCE_DATE_EPOCH, deterministic)
- [x] 1.3 flake.nix: `packages.test-image-x86` and `packages.test-image-arm`; document the aarch64 build path (native / remote / binfmt+QEMU)
- [x] 1.4 Build `test-image-x86` and verify it produces a valid .raw.xz (hermetic, input-pinned; disk images aren't byte-reproducible)

## 2. Acceptance tests

- [x] 2.1 `internal/provider/image_acc_test.go`: TF_ACC-gated, skips without HCLOUD_TOKEN; composes hcloudimage_image + hcloud_server (pinned ~> 1.48)
- [x] 2.2 SSH-reachability assertion from the runner with the baked key (read /etc/os-release), not just `running`
- [x] 2.3 Both x86 (cx22) and arm (cax11, toggle-gated); guaranteed cleanup on failure
- [x] 2.4 Confirm `go test ./...` still passes with acceptance skipped (no token)

## 3. Workflows

- [x] 3.1 `.github/workflows/acceptance.yml`: dispatch + push-main + nightly; never fork PRs; concurrency-limited; arm behind an input; always-run cleanup; HCLOUD_TOKEN secret
- [x] 3.2 `.github/workflows/cleanup.yml`: nightly label-scoped `hcloud-upload-image cleanup`
- [x] 3.3 Validate workflow YAML; confirm all §8.3/§15 cost/safety controls are encoded

## 4. Docs + close out

- [x] 4.1 README/docs: running acceptance (needs HCLOUD_TOKEN, budget-limited project), aarch64 path, cost controls
- [x] 4.2 `nix flake check` still green; `go build`/`go test` green
- [x] 4.3 Archive OpenSpec change; complete beans epics + milestone 07; commit with jj (Pim Snel, no self-promotion)
