---
id: T-1644
type: task
nature: feature
title: A red batch is bisected to the first story with which it fails, with each step kept in the queue's state
status: done
parent: S-0400
owner: alex
created: 2026-10-16T05:01:03Z
updated: 2026-10-19T22:32:28Z
transitions:
  - to: ready
    at: 2026-10-19T22:24:54Z
    by: agent-S-0400
  - to: in-progress
    at: 2026-10-19T22:24:54Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T22:32:28Z
    by: agent-S-0400
stream: S-0400
tags: [cli]
touches: [flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/progress.go]
after: [T-1909]
usage:
  source: log
  seconds: 454
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 46
      output: 22887
      cache_read: 3313379
      cache_write: 94310
      cost: 1.6658
---
# T-1644 A red batch is bisected to the first story with which it fails, with each step kept in the queue's state

## Work

- `bisect.go`: `Bisect` takes a batch whose result failed and finds, by halving, the first story with which it fails. It verifies main with the first half of the batch, then halves again. Each prefix is checked out at its last story's rebased head (`Batch.Heads`) in the batch's worktree and verified with `Deps.Verify`. Main itself, the batch's base, counts as passing. A batch of one is its own answer.
- Each step, its prefix, commit, and result, is written to `Batch.Bisect` as it begins and ends. A bisection stopped by ctx resumes from the steps kept.
- Once it finds the story, it sets `Batch.Landing` to the stories before it. The stories after it stay enqueued.
- `SendBack` moves the failing story from review to in-progress as `flai`, with a reason naming the failing tier and its first findings, which the move appends to its notes. It opens a thread on the story to its agent saying what failed and at which commit, takes the story out of the queue, and records `OutcomeFailing`.
- Test with a fake verify: a failure at each position, a failure that needs two stories together, and a resumed bisection.
- Waits for T-1909, whose `Heads` and `Bisect` it reads and writes.

## Done when

Bisection finds the failing story in every tested case and records its steps. The send-back moves the story, writes its notes, and opens the thread. The tests pass.

## Notes
