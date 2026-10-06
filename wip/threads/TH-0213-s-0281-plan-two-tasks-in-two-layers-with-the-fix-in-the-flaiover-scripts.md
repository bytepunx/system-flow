---
id: TH-0213
title: "S-0281 plan: two tasks in two layers, with the fix in the flaiover scripts"
anchor:
  path: wip/kanban/stories/S-0281-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md
  item: S-0281
status: resolved
participants: [planner-S-0281, alex]
created: 2026-10-06T22:54:17Z
updated: 2026-10-06T23:02:40Z
---

# TH-0213 S-0281 plan: two tasks in two layers, with the fix in the flaiover scripts

On wip/kanban/stories/S-0281-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md.

## Entries

### 2026-10-06T22:54:17Z planner-S-0281
The plan for S-0281 is two tasks, one per layer.

**Layer 1: T-1088.** The flaiover scripts install flaiover's dependencies in a story worktree that lacks them.
- `scripts/flaiover-install.sh` gets a mode that installs only when `flaiover/node_modules` is missing, or is older than `flaiover/pnpm-lock.yaml`. It runs `pnpm install --frozen-lockfile` against the shared pnpm store in the main checkout's `.flai-cache`, which took 3s in S-0295.
- `scripts/flaiover-test.sh` (close-out's flaiover step) always calls that mode first.
- `scripts/flaiover-unit.sh` (the vitest run in `make test` and `flai-test.sh`) calls it in a story worktree only.
- Touches: those three scripts.

**Layer 2: T-1095, after T-1088.** It updates the docs and closes I-0080 with a reason that names S-0281. Touches: `scripts/README.md`, `design/system/devex.md` (the Test tiers row), `docs/operators/index.md`, the I-0080 file, and `design/issues/summary.md`.

**Assumptions.** Each is my recommendation. Say if you want otherwise.

1. **The fix goes in the project's scripts, not in `flai stream open`.** `stream open` is generic, and only this repository has a `flaiover/`. The other option is a manifest key for commands to run after a worktree is made, such as `worktree.setup`. That changes flai and the template, needs an ADR, and is a larger story. Propose it separately if you want it.
2. **In a story worktree, `flaiover-unit.sh` installs too, instead of skipping.** Today a worktree that changes `flai/` skips flaiover's vitest without a word, which is the same cause. In the main checkout it keeps skipping while `flaiover/node_modules` is missing, as it does now.
3. **The test is a recorded reproduction, not an automated one.** T-1088 reproduces the failure in a scratch worktree with no `node_modules` and records the before and after in the narrative. The criterion says "where one fits", and a test of a shell script that needs pnpm and the network does not fit `make test`.
4. **The figures.**
   - Forecast: 30m, double flai's 15m, because the reproduction runs flaiover's lint, svelte-check, and vitest more than once. Delivery 2026-10-07T09:52Z.
   - Cost of delay: 12.50 USD a week, from the 5m-per-cycle input flai set from I-0080.
   - I-0080 has had a second occurrence since then (S-0295). If you want the cost to count both, raise `time_lost_per_cycle` to 10m. I left the inputs as they are, since they are yours.

I have planned what needs no answer, so nothing waits on this thread. Reply `ok` to keep the plan, or say what to change.

### 2026-10-06T23:02:40Z alex
Resolved.
