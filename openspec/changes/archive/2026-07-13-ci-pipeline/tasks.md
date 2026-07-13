## 1. Docs + lint config

- [x] 1.1 Generate `docs/` with tfplugindocs and commit (index, resources/image, data-sources/snapshot)
- [x] 1.2 Add `.golangci.yml` matching the linters the devShell runs

## 2. ci.yml

- [x] 2.1 lint-build-test job: golangci-lint, go build, go test with coverage → Codecov (TF_ACC + terraform on PATH via nix develop)
- [x] 2.2 validate-examples job: `nix build .#default` + scripts/validate-examples.sh
- [x] 2.3 flake-check job: `nix flake check` (hermetic lifecycle test); document KVM/runner requirement
- [x] 2.4 docs job: `tfplugindocs generate` then `git diff --exit-code docs/`
- [x] 2.5 Use DeterminateSystems nix installer + cache so CI == local toolchain

## 3. Close out

- [x] 3.1 Validate the workflow YAML; confirm each step mirrors a locally-verified command; docs diff-clean
- [x] 3.2 Archive OpenSpec change; complete beans epic + milestone 06; commit with jj (Pim Snel, no self-promotion)
