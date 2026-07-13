## 1. Mirror + test config

- [x] 1.1 Shared mirror helper: derivation building the filesystem-mirror layout from packages.default (reused by milestone 09)
- [x] 1.2 Test HCL under test/e2e/ using only hcloudimage_image with image_path + image_sha256 var; a CLI config selecting the mirror

## 2. NixOS VM test

- [x] 2.1 test/e2e/hermetic.nix: runNixOSTest node with terraform, opentofu, the mirror, and the test config in the store
- [x] 2.2 Python driver: for terraform and tofu — init; apply(sha=A,labels=a); plan(sha=B) expect replace; apply; plan(labels=b) expect in-place; apply; destroy expect empty state; all with HCLOUDIMAGE_FAKE=1
- [x] 2.3 Assert plan actions via -json resource_changes (["delete","create"]=replace, ["update"]=in-place)

## 3. Flake wiring

- [x] 3.1 flake.nix: checks.hermetic-e2e on Linux systems (plain nix, no flake-utils)
- [x] 3.2 Run `nix flake check` — hermetic test builds and passes

## 4. Close out

- [x] 4.1 Confirm the PoC DoD gate is green (nix flake check runs the VM test)
- [x] 4.2 Archive OpenSpec change; complete beans epic + milestone 05; commit with jj (Pim Snel, no self-promotion)
