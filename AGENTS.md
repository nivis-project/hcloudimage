# AGENTS.md — Autonomous build charter

This repository is built **autonomously by Claude Code**. This file is the operating
contract. Read it fully before doing anything. The authoritative product spec is
[`BRIEFING.md`](./BRIEFING.md); this file governs *how* you work, `BRIEFING.md` governs
*what* you build.

> **Goal of this phase:** deliver the PoC / alpha base described in `BRIEFING.md`.
> The completion gate for the PoC is **BRIEFING.md §12 (Definition of Done)**, built in
> the order of **BRIEFING.md §13 (Milestones)**.

---

## 0. The three tools you must use

| Concern | Tool | Never use instead |
|---|---|---|
| Milestones & epics (the roadmap) | **beans** | TodoWrite, ad-hoc todo lists |
| Proposals & task tracking within an epic | **OpenSpec** | inventing your own task files |
| Version control | **jj** (Jujutsu, colocated with git) | `git commit` directly |

The remote is `git@github.com:nivis-project/hcloudimage.git` (tracked bookmark `main`).

---

## 1. The build loop (do this for every epic)

Work one **epic** at a time, in milestone order. For each epic:

1. **Pick the work.** `beans list --json --ready` shows unblocked milestones/epics.
   Milestones are chained `01 → 09`; only the current milestone's epics are ready.
   Read the epic: `beans show --json <epic-id>`.

2. **Open the epic in beans.** `beans update <epic-id> -s in-progress`.

3. **Create an OpenSpec change** for the epic's work:
   `/opsx:propose "<short description>"` (or `openspec new change "<kebab-name>"`).
   Generate `proposal.md` (what & why), `design.md` (how), and **`tasks.md`** (the
   checklist). `tasks.md` is the single source of truth for tasks within the epic —
   **do not** track tasks anywhere else.

4. **Record the link** back on the bean so the two systems stay in sync:
   append the change name to the epic body, e.g.
   `beans update <epic-id> --body-append "OpenSpec change: <change-name>"`.

5. **Implement** the change via `/opsx:apply`, checking off `tasks.md` items as you go.
   Keep the code hermetic and reproducible (see §3).

6. **Test** — nothing is "done" until it is proven (see §2). Run the relevant layer.

7. **Archive the change** when all its tasks are checked and tests pass:
   `/opsx:archive <change-name>`. This syncs deltas into `openspec/specs/` and moves the
   change to `openspec/changes/archive/`.

8. **Close the bean.** When the epic's every task is checked, add a
   `## Summary of Changes` section and `beans update <epic-id> -s completed`.
   When all epics under a milestone are completed, mark the milestone completed too.

9. **Commit with jj** — see §4. **Commit immediately after every OpenSpec archive.**

Then loop to the next ready epic.

---

## 2. Testing is not optional (BRIEFING.md §8)

The PoC must be *demonstrably* working. Three layers, in priority order for the PoC:

- **Unit** (Go, no network) — schema, validators, plan-modifiers, config→request
  mapping, lifecycle against the fake uploader. Runs everywhere, always.
- **Hermetic NixOS-VM lifecycle test** (`checks.hermetic-e2e`, gated by
  `nix flake check`) — **this is the PoC Definition-of-Done gate** (BRIEFING.md §8.2).
  It proves the full Terraform ↔ provider protocol path (init/plan/apply/destroy,
  ForceNew vs in-place) hermetically, under both `terraform` and `tofu`, at zero cloud
  cost, using the fake uploader.
- **Acceptance** (real, billable Hetzner) — gated CI only, with mandatory cost controls
  and guaranteed cleanup (BRIEFING.md §8.3). Needed for full DoD; a stubbed+documented
  version is acceptable to *reach* the alpha base as long as the workflow files exist.

Write e2e test cases that assert real behaviour, not just "it ran". For the hermetic
layer that means: apply yields synthetic `id`; changing `image_sha256` forces replace;
changing `labels` does **not**; destroy clears state.

---

## 3. Nix from the start (BRIEFING.md §7) — hard constraints

- The project **must** work with nix and nix flakes from day one.
- **Do not use `flake-utils`.** Use plain nix (`nixpkgs.lib.genAttrs` + a `forAllSystems`
  helper) to enumerate supported systems. This is already set up in `flake.nix`.
- Builds are hermetic and reproducible; no network in the build beyond fixed-output
  vendoring (`buildGoModule` with a pinned `vendorHash`).
- `flake.nix` will grow to expose: `devShells.default`, `packages.default`,
  `checks.hermetic-e2e`, `packages.test-image-{x86,arm}`, `packages.provider-mirror`.
  Add each as its milestone reaches it; the current file documents the roadmap inline.
- Keep `flake.lock` committed. After editing `flake.nix`, run `nix flake check`.

---

## 4. Version control with jj (Jujutsu)

The repo is a **colocated** jj/git repo (`.jj/` + `.git/`). Use jj, not raw git.

**Commit as Pim Snel. No self-promotion** — no "Generated with Claude", no
`Co-Authored-By` trailer, no tool advertising in messages. Author is already configured
(`Pim Snel <post@pimsnel.com>`).

Commit cadence: **commit after every OpenSpec change archival** (step 9 above). A commit
should contain the code changes **and** the beans/OpenSpec file changes together, so the
tracker state and the code never drift.

Typical flow with jj:

```bash
# Describe the current change (the working-copy commit @) and start a fresh one:
jj commit -m "<type>: <summary>

<body — what changed and why, in Pim's voice>"

# Advance the main bookmark to the commit you just finished, then push:
jj bookmark set main -r @-
jj git push
```

Write conventional-commit style subjects (`feat:`, `test:`, `chore:`, `docs:`,
`ci:`) — BRIEFING.md §10 derives the changelog from them. Keep messages factual and
first-person-neutral.

> **Remote push:** Pim provides the remote URL (already set to
> `git@github.com:nivis-project/hcloudimage.git`). If a push is rejected for auth, stop
> and report — do not invent credentials.

---

## 5. Milestone map (beans ⇄ BRIEFING.md §13)

Milestones are titled with an incremental two-digit prefix starting at `01`. Each has
child epics already created in beans. See `beans roadmap` for the live tree.

| # | Milestone | Gate |
|---|---|---|
| 01 | Scaffold | `go build`, `nix develop` |
| 02 | Schema, validators, fake uploader | unit tests green |
| 03 | Lifecycle against fake | lifecycle tests green |
| 04 | Real uploader + examples | `terraform`/`tofu validate` |
| 05 | Hermetic NixOS-VM test | **`nix flake check` — PoC DoD gate** |
| 06 | CI pipeline | `ci.yml` lint/unit/hermetic/validate/docs |
| 07 | Fixtures + acceptance | Alpine fixtures + gated acceptance + cleanup |
| 08 | Release engineering | goreleaser + signing + docs + `v0.1.0` |
| 09 | Registry + Nix mirror | registry publish + Nivis mirror consumption |

**Human-only prerequisites** (BRIEFING.md §14 — GPG key, GitHub secrets, registry
onboarding) are **out of scope**: stub and document them, then move on. Do not block on
them.

---

## 6. Scope discipline

- Follow `BRIEFING.md` exactly where it gives concrete config (goreleaser, manifest,
  schema attributes, workflow triggers) — those are requirements, not suggestions.
- Where it describes behaviour in prose, satisfy the behaviour; you choose the code.
- Do not add multi-cloud, image building, or server management — those are explicit
  non-goals (BRIEFING.md §2).
- If a decision is genuinely ambiguous and not covered here or in `BRIEFING.md`, prefer
  the choice that keeps the PoC hermetic, reproducible, and testable — and note it in the
  OpenSpec `design.md`.
