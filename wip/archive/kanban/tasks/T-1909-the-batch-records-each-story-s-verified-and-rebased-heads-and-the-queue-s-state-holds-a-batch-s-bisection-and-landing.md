---
id: T-1909
type: task
nature: feature
title: The batch records each story's verified and rebased heads, and the queue's state holds a batch's bisection and landing
status: done
parent: S-0400
owner: alex
created: 2026-10-19T22:02:02Z
updated: 2026-10-19T22:24:13Z
transitions:
  - to: ready
    at: 2026-10-19T22:02:57Z
    by: agent-S-0400
  - to: in-progress
    at: 2026-10-19T22:04:39Z
    by: agent-S-0400
  - to: done
    at: 2026-10-19T22:24:13Z
    by: agent-S-0400
stream: S-0400
tags: [cli]
touches: [flai/internal/mergequeue/progress.go, flai/internal/mergequeue/progress_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/state.go]
usage:
  source: log
  seconds: 1174
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 15848
      cache_read: 2294352
      cache_write: 65305
      cost: 1.1535
---
# T-1909 The batch records each story's verified and rebased heads, and the queue's state holds a batch's bisection and landing

## Work

- A new `progress.go` defines the types the landing and the bisection keep in the state: `Head` (a batched story, the verified head the batch rebased a copy of, and that copy's head in the batch), `Bisect` with its `BisectStep`s, and `Landing`. It also defines the outcomes `OutcomeLanded` and `OutcomeFailing`.
- `state.go` gets three fields at the end of `Batch` and nothing else: `Heads`, `Bisect`, and `Landing`. They wait on S-0401's answer on MS-0143, and the two files join this task's touches once it agrees.
- `batch.go`'s `run` records each story's `Head` where its rebase succeeds, and nothing else changes.
- It waits for no task. It is first because every other code task reads these types.

## Done when

A built batch's state names each kept story's verified and rebased heads, the new fields round-trip through `Update`, and the tests pass.

## Notes
