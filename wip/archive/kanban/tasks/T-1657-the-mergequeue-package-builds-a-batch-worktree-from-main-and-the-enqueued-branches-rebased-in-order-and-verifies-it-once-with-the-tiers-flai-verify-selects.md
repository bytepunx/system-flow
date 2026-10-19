---
id: T-1657
type: task
nature: feature
title: The mergequeue package builds a batch worktree from main and the enqueued branches rebased in order, and verifies it once with the tiers flai verify selects
status: done
parent: S-0399
owner: alex
created: 2026-10-16T05:01:58Z
updated: 2026-10-19T19:49:24Z
transitions:
  - to: ready
    at: 2026-10-19T19:21:16Z
    by: agent-S-0399
  - to: in-progress
    at: 2026-10-19T19:21:16Z
    by: agent-S-0399
  - to: done
    at: 2026-10-19T19:49:24Z
    by: agent-S-0399
stream: S-0399
tags: [cli]
touches: [flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/mergequeue/state.go]
after: [T-1641]
usage:
  source: log
  seconds: 1688
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 102
      output: 47225
      cache_read: 7131529
      cache_write: 261568
      cost: 3.7742
---
# T-1657 The mergequeue package builds a batch worktree from main and the enqueued branches rebased in order, and verifies it once with the tiers flai verify selects

## Work

- `batch.go`: build a worktree under `.flai-cache/worktrees/` on a batch branch from main. Rebase each enqueued story's branch onto the one before, in order, up to `batch_max`.
- A story whose rebase stops is left out, with the stop's paths as its outcome. The rebase is aborted and the next story goes on.
- Verify the batch once: add to `verify` an entry point that runs the tiers selected by the batch's paths against main, in a given worktree, without writing any story's report. Record the step and result in the state as each tier runs.
- Test a batch of two clean stories, one that stops, and a failing tier.
- Waits for T-1641, whose state writer it records in.

## Done when

A batch builds and verifies as above, no story's own verify report changes, and the tests pass.

## Notes
