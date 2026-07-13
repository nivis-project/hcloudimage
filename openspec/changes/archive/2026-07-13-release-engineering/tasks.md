## 1. Release artifacts config

- [x] 1.1 `terraform-registry-manifest.json` = { version 1, metadata.protocol_versions ["6.0"] }
- [x] 1.2 `.goreleaser.yml` (v2): build matrix, zip archive name, manifest embedded, SHA256SUMS, GPG signs block, conventional-commit changelog, ldflags version
- [x] 1.3 `goreleaser check` passes
- [x] 1.4 `goreleaser build --snapshot` compiles the matrix (no tag / no key)

## 2. Release workflow

- [x] 2.1 `.github/workflows/release.yml`: on tag v*.*.*, import GPG key, `goreleaser release --clean` with GPG_FINGERPRINT/PASSPHRASE/GPG_PRIVATE_KEY secrets, contents:write
- [x] 2.2 Validate the workflow YAML

## 3. Changelog + v0.1.0

- [x] 3.1 `CHANGELOG.md` (Keep a Changelog) with the 0.1.0 entry
- [x] 3.2 Document that cutting the signed v0.1.0 tag/release is a human step (GPG key + registry onboarding, §14)

## 4. Close out

- [x] 4.1 `nix flake check` still green; goreleaser check + snapshot verified
- [x] 4.2 Archive OpenSpec change; complete beans epics + milestone 08; commit with jj (Pim Snel, no self-promotion)
