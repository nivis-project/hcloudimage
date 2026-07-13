## 1. Mirror consumption

- [x] 1.1 `scripts/verify-mirror.sh`: build the mirror, write a filesystem_mirror CLI config, run a pinned-version consumer through init + plan under terraform and tofu (HCLOUDIMAGE_FAKE=1, no token)
- [x] 1.2 `make consume-mirror` target wrapping the script
- [x] 1.3 Run it; confirm both binaries resolve the provider from the mirror and plan succeeds

## 2. Docs

- [x] 2.1 `docs/consuming-from-nix-mirror.md`: filesystem_mirror recipe (portable) + dev_overrides alternative, with real paths and the tofu/terraform init difference
- [x] 2.2 `docs/publishing.md`: Terraform + OpenTofu registry onboarding (human §14 steps) and prerequisites

## 3. Close out

- [x] 3.1 `nix flake check` still green
- [x] 3.2 Archive OpenSpec change; complete beans epic + milestone 09; commit with jj (Pim Snel, no self-promotion)
