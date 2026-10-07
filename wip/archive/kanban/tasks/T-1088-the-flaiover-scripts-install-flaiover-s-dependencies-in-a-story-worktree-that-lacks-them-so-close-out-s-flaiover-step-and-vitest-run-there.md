---
id: T-1088
type: task
nature: improvement
title: The flaiover scripts install flaiover's dependencies in a story worktree that lacks them, so close-out's flaiover step and vitest run there
status: done
parent: S-0281
owner: alex
created: 2026-10-06T22:52:53Z
updated: 2026-10-07T00:55:22Z
transitions:
  - to: ready
    at: 2026-10-07T00:50:34Z
    by: agent-S-0281
  - to: in-progress
    at: 2026-10-07T00:50:35Z
    by: agent-S-0281
  - to: done
    at: 2026-10-07T00:55:22Z
    by: agent-S-0281
stream: S-0281
tags: [devex]
touches: [scripts/flaiover-install.sh, scripts/flaiover-test.sh, scripts/flaiover-unit.sh]
---
# T-1088 The flaiover scripts install flaiover's dependencies in a story worktree that lacks them, so close-out's flaiover step and vitest run there

## Work

`flai stream open` makes a story worktree with `git worktree add` (`flai/cmd/branch.go`), which leaves out the git-ignored `flaiover/node_modules`. Then `scripts/close-out.sh` runs `scripts/flaiover-test.sh` for a change under `flaiover/`, and `pnpm lint` stops with `prettier: not found` (I-0080, S-0217 and S-0295). `scripts/flaiover-unit.sh`, which `flai-test.sh` runs for a change under `flai/`, exits 0 without `node_modules`, so in a worktree flaiover's vitest is skipped without a word. This is layer 1 and waits for no task.

Fix it in the project's scripts, not in flai. `flai stream open` is generic, and only this project has a `flaiover/`.

- Give `scripts/flaiover-install.sh` a mode that installs only when it has to: when `flaiover/node_modules` is missing, or when `flaiover/pnpm-lock.yaml` is newer than `flaiover/node_modules/.modules.yaml`. It runs `pnpm install --frozen-lockfile` against the shared store `env.sh` sets in the main checkout's `.flai-cache`, which takes seconds. Each run that installs says so on stderr.
- `scripts/flaiover-test.sh` calls that mode before `pnpm lint`. It cannot run without the dependencies, wherever it runs.
- `scripts/flaiover-unit.sh` calls that mode in a story worktree, where `CACHE_ROOT` differs from `ROOT`. In the main checkout it keeps skipping when `flaiover/node_modules` is missing.
- If the install fails, the script stops with pnpm's error, not a later `not found`.

## Done when

- Reproduced and fixed. In a scratch worktree with no `flaiover/node_modules` (`git worktree add` under `.flai-cache/`, removed afterwards), `scripts/flaiover-test.sh` installs and then passes, and `scripts/flaiover-unit.sh` installs and runs vitest. Run in the same tree a second time, neither one installs again. Record the commands and their output in the narrative: no automated test fits a shell script that needs pnpm and the network.
- In the main checkout, `scripts/flaiover-unit.sh` still exits 0 without installing when `flaiover/node_modules` is missing.
- `scripts/lint-md.sh` and `scripts/check.sh` pass.

## Notes

Written by the planner. If the fix needs a change in `scripts/env.sh` or `scripts/close-out.sh`, widen the task's touches to that file and name it on the plan thread.
