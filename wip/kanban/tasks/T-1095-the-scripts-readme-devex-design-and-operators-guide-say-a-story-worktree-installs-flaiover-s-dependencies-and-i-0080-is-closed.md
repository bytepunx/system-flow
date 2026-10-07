---
id: T-1095
type: task
nature: improvement
title: The scripts README, devex design, and operators' guide say a story worktree installs flaiover's dependencies, and I-0080 is closed
status: done
parent: S-0281
owner: alex
created: 2026-10-06T22:53:04Z
updated: 2026-10-07T00:55:51Z
transitions:
  - to: ready
    at: 2026-10-07T00:55:22Z
    by: agent-S-0281
  - to: in-progress
    at: 2026-10-07T00:55:23Z
    by: agent-S-0281
  - to: done
    at: 2026-10-07T00:55:51Z
    by: agent-S-0281
stream: S-0281
tags: [devex, docs]
touches: [scripts/README.md, design/system/devex.md, docs/operators/index.md, design/issues/I-0080-a-story-worktree-has-no-flaiover-node-modules-so-the-close-out-stops-at-the-flaiover-step-until-flaiover-install-runs-there.md, design/issues/summary.md]
after: [T-1088]
---
# T-1095 The scripts README, devex design, and operators' guide say a story worktree installs flaiover's dependencies, and I-0080 is closed

## Work

Layer 2. It waits for T-1088, whose behavior it describes. Say what T-1088 built, in the words it built, and close the issue it remedies.

- `scripts/README.md`: update the rows for `flaiover-install.sh`, `flaiover-test.sh`, and `flaiover-unit.sh` to say when each installs flaiover's dependencies.
- `design/system/devex.md`, the Test tiers row: `make test` and `make flai-test` run flaiover's vitest in a story worktree after installing its dependencies there. In the main checkout they still skip it while `flaiover/node_modules` is missing.
- `docs/operators/index.md`, the paragraph on what `scripts/flai-test.sh` runs at close-out: the same.
- Close I-0080 with `flai issue close I-0080 --reason "..."`, run from the story worktree. Give as the reason that the flaiover scripts now install the dependencies in a story worktree that lacks them, and name S-0281. It writes the issue file and `design/issues/summary.md`.

## Done when

- The three documents describe the new behavior, and none of them still says a story worktree needs `flaiover-install` run by hand.
- I-0080 is closed, with a reason that names S-0281, and `design/issues/summary.md` agrees.
- `scripts/lint-md.sh` and `scripts/check.sh` pass.

## Notes

Written by the planner. If `flai issue close` writes outside these two issue paths, widen the task's touches to what it wrote.
