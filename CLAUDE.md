# CLAUDE.md

You are building this repository **autonomously**.

**Start here, in order:**

1. **[`AGENTS.md`](./AGENTS.md)** — the operating contract: the build loop, the three
   tools (beans / OpenSpec / jj), testing requirements, the Nix constraint, and commit
   rules. Follow it exactly.
2. **[`BRIEFING.md`](./BRIEFING.md)** — the authoritative product spec for the provider.
   Completion gate is its **§12 (Definition of Done)**; build order is its **§13**.

**Non-negotiables (full detail in `AGENTS.md`):**

- Track the roadmap in **beans** (milestones + epics), not TodoWrite. `beans prime` for usage.
- Track tasks-within-an-epic in **OpenSpec** (`tasks.md`); one change per epic.
- Version control is **jj** (colocated). Commit as **Pim Snel, no self-promotion**.
  **Commit after every OpenSpec archive.**
- Nix flakes from the start; **no `flake-utils`** — plain nix `forAllSystems`.
- Prove it works: unit + hermetic NixOS-VM (`nix flake check`) + gated acceptance e2e.

**The build loop, in one line:** pick a ready epic in beans → open an OpenSpec change
(`/opsx:propose`) → implement (`/opsx:apply`) → test → archive (`/opsx:archive`) →
close the bean → `jj commit` + push.

Run `beans list --ready` to see what to work on next.
