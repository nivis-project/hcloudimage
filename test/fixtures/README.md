# Acceptance test fixtures

## `throwaway_ed25519` / `throwaway_ed25519.pub`

A **deliberately disposable** SSH keypair used only by the billable acceptance
tests (BRIEFING.md §8.3). The public key is baked into the Alpine test images
built by `packages.test-image-{x86,arm}`; the private key is used by the
acceptance test to SSH from the CI runner into the booted Hetzner server and
read `/etc/os-release`, proving real guest reachability (not just that the
server reports `running`).

This key protects nothing of value:

- It only ever authorizes login to **throwaway** servers created and destroyed
  within a single acceptance run, in an isolated, budget-limited Hetzner
  project.
- It is committed on purpose so the fixture image and the test agree on the key
  with zero setup.

Do **not** reuse it for anything else. If you want to rotate it:

```sh
ssh-keygen -t ed25519 -N "" -C "hcloudimage-acceptance-throwaway" \
  -f test/fixtures/throwaway_ed25519
```

then rebuild the fixtures so the new public key is baked in.
