---
id: S-0481
type: story
nature: improvement
title: pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY
status: done
owner: alex
created: 2026-10-21T12:00:00Z
updated: 2026-10-23T07:23:39Z
transitions:
  - to: ready
    at: 2026-10-22T20:55:39Z
    by: alex
  - to: in-progress
    at: 2026-10-23T07:15:51Z
    by: agent-S-0481
  - to: review
    at: 2026-10-23T07:21:02Z
    by: agent-S-0481
  - to: done
    at: 2026-10-23T07:23:39Z
    by: orchestrator
tags: []
touches: [design/tech/node.md, scripts/README.md, scripts/flaiover-install-test.sh, scripts/flaiover-install.sh, scripts/flaiover-pnpm.sh, scripts/flaiover-test.sh, design/issues/I-0177-pnpm-add-run-bare-in-a-story-worktree-installs-from-the-default-store-so-the-next-flaiover-install-sh-aborts-with-err-pnpm-aborted-remove-modules-dir-no-tty.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 324
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 1
      test_runs: 2
      hand_edits: 3
      work: 32
  models:
    - model: claude-opus-5-5
      input: 78
      output: 22888
      cache_read: 5098919
      cache_write: 159776
      cost: 2.7561
  strategic:
    - kind: orchestrator
      seconds: 494
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 162
          output: 2758
          cache_read: 62371476
          cache_write: 30344
          cost: 13.7786
        - model: claude-sonnet-5-5
          input: 8
          output: 52
          cache_read: 78632
          cache_write: 41873
          cost: 0.0834
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-21T12:00:00Z
finalized:
  by: alex
  at: 2026-10-22T20:52:11Z
---
# S-0481 pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY

## Goal

This story remediates [I-0177](../../../design/issues/I-0177-pnpm-add-run-bare-in-a-story-worktree-installs-from-the-default-store-so-the-next-flaiover-install-sh-aborts-with-err-pnpm-aborted-remove-modules-dir-no-tty.md), "pnpm add run bare in a story worktree installs from the default store, so the next flaiover-install.sh aborts with ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0177 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0177 is closed with `flai issue close I-0177 --reason` saying what fixed it

## Tasks
- T-2070 flaiover-install.sh purges an install from another store instead of aborting, and flaiover-pnpm.sh adds a dependency with the project's pnpm and store
- T-2071 I-0177 is closed as fixed by the purge and flaiover-pnpm.sh
- T-2072 design/system/devex.md is left as main has it, so S-0481 does not conflict with S-0466's change to its test tiers row

## Notes

No test tier selects `scripts/`, so the close-out ran only the markdown lint. `scripts/flaiover-install-test.sh` was run by hand: it failed against main's `scripts/flaiover-install.sh` with `ERR_PNPM_ABORTED_REMOVE_MODULES_DIR_NO_TTY` and passed with this branch's. It now runs first in `scripts/flaiover-test.sh`, the `flaiover` tier. A real `scripts/flaiover-install.sh --if-needed` in this worktree passed in 3.6s on the project's store.

Cost of delay inputs set by flai from I-0177. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T08:25:41Z, 0.1 days before this story; under one cycle counts as one).

### Accepted by the orchestrator

- Verified: 351d6f1b85cd8c9cd977aeb010204be76773cc5e
- Batch: 20261023-072330
- Enqueued verified: cbd98d1da3172307346c1fee8b4efd8fa65149c6
- At: 2026-10-23T07:23:39Z

Verdict: pass at cbd98d1da3172307346c1fee8b4efd8fa65149c6
- 1: scripts/flaiover-install.sh, scripts/flaiover-pnpm.sh, scripts/flaiover-install-test.sh, scripts/flaiover-test.sh, scripts/README.md, design/tech/node.md
- 2: design/issues/I-0177-pnpm-add-run-bare-in-a-story-worktree-installs-from-the-default-store-so-the-next-flaiover-install-sh-aborts-with-err-pnpm-aborted-remove-modules-dir-no-tty.md, design/issues/summary.md
