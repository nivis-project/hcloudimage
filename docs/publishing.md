# Publishing to the Terraform and OpenTofu registries

Both registries ingest the identically-signed GitHub release assets produced by
`release.yml` / goreleaser (see [`.goreleaser.yml`](../.goreleaser.yml) and
BRIEFING §10). The steps below are **human-only prerequisites** (BRIEFING §14) —
the agent produces everything up to the signed release; a maintainer with the
GPG key and org access performs the onboarding and cuts the release.

## Signing key

The provider uses a **dedicated** release-signing GPG key (not a personal
identity), so it can be rotated or revoked independently:

- **Fingerprint:** `74F05F879B947F24006761E3FC80F1F128669C1B`
- **UID:** `Nivis hcloudimage provider <post@pimsnel.com>`
- **Public key (committed):** [`keys/hcloudimage-signing.pub.asc`](../keys/hcloudimage-signing.pub.asc)
  — this is what you register with each registry namespace.
- **Private key:** kept in the maintainer's password manager (Bitwarden), and
  mirrored into the repo's GitHub Actions secrets for `release.yml`. It is
  **never** committed.

## Prerequisites (human, §14)

- The `github.com/nivis-project/terraform-provider-hcloudimage` repo and org
  settings.
- Repository secrets used by the workflows:
  - `GPG_PRIVATE_KEY` — the private key from Bitwarden (ASCII-armored).
  - `GPG_FINGERPRINT` — `74F05F879B947F24006761E3FC80F1F128669C1B`.
  - `PASSPHRASE` — the key's passphrase (empty/omit if the key has none).
  - `HCLOUD_TOKEN` — a dedicated, budget-limited Hetzner project for acceptance.
  - `CODECOV_TOKEN` — coverage upload (optional for public repos).

## Cut a signed release

```sh
git tag v0.1.0
git push origin v0.1.0    # triggers .github/workflows/release.yml
```

`release.yml` imports the GPG key and runs `goreleaser release --clean`,
producing the zip archives, `..._SHA256SUMS`, and the `.sig` the registries
verify.

## Terraform Registry

1. Sign in at <https://registry.terraform.io> with the `nivis-project` org.
2. Publish the provider: connect the GitHub repo and add the **GPG public key**
   to the namespace's signing keys.
3. The registry then picks up tagged releases automatically — each new
   `vX.Y.Z` tag that `release.yml` publishes appears as a provider version.

The provider address is `nivis-project/hcloudimage`, protocol `6.0` (declared in
[`terraform-registry-manifest.json`](../terraform-registry-manifest.json)).

## OpenTofu Registry

OpenTofu reuses the same signed release artifacts. Onboarding is a submission to
[`github.com/opentofu/registry`](https://github.com/opentofu/registry)
referencing the repo and the GPG key.

1. Open a submission PR to `opentofu/registry` for `nivis-project/hcloudimage`.
2. Provide the repo URL and the GPG public key so OpenTofu can verify the
   signed `SHA256SUMS`.
3. Confirm the current submission procedure at submission time — it evolves;
   follow the registry's `CONTRIBUTING`/issue templates.

Once accepted, tagged releases flow to OpenTofu the same way as Terraform.
