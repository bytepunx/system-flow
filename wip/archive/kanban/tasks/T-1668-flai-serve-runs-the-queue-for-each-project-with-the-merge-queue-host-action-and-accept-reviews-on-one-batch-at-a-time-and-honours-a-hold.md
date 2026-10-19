---
id: T-1668
type: task
nature: feature
title: flai serve runs the queue for each project with the merge-queue host action and accept_reviews on, one batch at a time, and honours a hold
status: done
parent: S-0399
owner: alex
created: 2026-10-16T05:02:49Z
updated: 2026-10-19T19:55:43Z
transitions:
  - to: ready
    at: 2026-10-19T19:49:48Z
    by: agent-S-0399
  - to: in-progress
    at: 2026-10-19T19:49:49Z
    by: agent-S-0399
  - to: done
    at: 2026-10-19T19:55:43Z
    by: agent-S-0399
stream: S-0399
tags: [cli]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/serve/serve.go]
after: [T-1640, T-1657]
usage:
  source: log
  seconds: 354
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 13923
      cache_read: 2102527
      cache_write: 77116
      cost: 1.1127
---
# T-1668 flai serve runs the queue for each project with the merge-queue host action and accept_reviews on, one batch at a time, and honours a hold

## Work

- New `serve/mergequeue.go`, started from `serve.go` beside the orchestrator's loop. For each served project with the `merge-queue` host action and `accept_reviews` on, `acceptance.mode` local, and the queue not held, it builds and verifies one batch at a time with T-1657's package when eligible stories are enqueued.
- Under pull-request acceptance it starts no batch and logs why once, not on every tick.
- It runs off the agents' turns: no agent's call waits on it. It logs each step.
- It stops cleanly with `flai serve`, leaving a batch's worktree to be built again at the next start.
- Test with a fake clock and fake verify: the action off, the permission off, pull-request acceptance, held, and one batch run.
- Waits for T-1640, which adds the host action, and T-1657, which builds and verifies the batch.

## Done when

`flai serve` runs batches as above, and the tests pass.

## Notes
