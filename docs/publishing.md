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

> **Status: not currently published to the Terraform Registry.** See the key-type
> caveat below.
>
> **⚠️ Key type:** the Terraform Registry accepts only **RSA or DSA** GPG keys,
> **not ECC/EdDSA (ed25519)**. Our current signing key
> (`keys/hcloudimage-signing.pub.asc`) is ed25519 — accepted by OpenTofu but
> **rejected by the Terraform Registry**. To publish here you must generate an
> **RSA-4096** signing key, add it as the `GPG_PRIVATE_KEY`/`GPG_FINGERPRINT`
> secrets, and cut a new release signed with it (you can sign with both keys, or
> switch to a single RSA key which both registries accept).

Once you have an RSA key and an RSA-signed release:

1. Sign in at <https://registry.terraform.io> **with GitHub** (not HCP — HCP is a
   separate paid product and is not involved). Authorize the Terraform Registry
   OAuth app for the `nivis-project` org.
2. Click **Publish → Provider**, select `nivis-project/terraform-provider-hcloudimage`,
   and add the **RSA GPG public key** to the namespace's signing keys.
3. The registry then picks up tagged releases automatically — each new
   `vX.Y.Z` tag that `release.yml` publishes appears as a provider version.

The provider address is `nivis-project/hcloudimage`, protocol `6.0` (declared in
[`terraform-registry-manifest.json`](../terraform-registry-manifest.json)).

## OpenTofu Registry

OpenTofu reuses the same signed release artifacts. Onboarding is done entirely
through **GitHub issue forms** on
[`github.com/opentofu/registry`](https://github.com/opentofu/registry) — **not**
pull requests. Their README is explicit: submissions must go through the issue
form UI in a browser; the `gh` CLI, the GitHub API, or PRs are rejected and
closed, because the automated validation depends on the structured issue-form
data.

There are two separate submissions, and the **signing key must be registered
first** because the provider validation verifies the release's
`SHA256SUMS.sig` against it:

1. **Submit new Provider Signing Key** —
   <https://github.com/opentofu/registry/issues/new?template=provider_key.yml>
   - Namespace: `nivis-project`
   - GPG public key: paste [`keys/hcloudimage-signing.pub.asc`](../keys/hcloudimage-signing.pub.asc)
     (fingerprint `74F05F879B947F24006761E3FC80F1F128669C1B`).
2. **Submit new Provider** —
   <https://github.com/opentofu/registry/issues/new?template=provider.yml>
   - For `nivis-project/terraform-provider-hcloudimage`.

(If the provider issue was opened before the key issue, that's fine — it
validates once the key submission is accepted.)

Once accepted, tagged releases flow to OpenTofu the same way as Terraform.
