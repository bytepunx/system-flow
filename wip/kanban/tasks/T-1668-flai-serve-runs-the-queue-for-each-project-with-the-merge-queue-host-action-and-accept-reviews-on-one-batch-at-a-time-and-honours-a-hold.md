---
id: T-1668
type: task
nature: feature
title: flai serve runs the queue for each project with the merge-queue host action and accept_reviews on, one batch at a time, and honours a hold
status: backlog
parent: S-0399
owner: alex
created: 2026-10-16T05:02:49Z
updated: 2026-10-18T02:20:23Z
transitions: []
stream: S-0399
tags: [cli]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/serve/serve.go]
after: [T-1640, T-1657]
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
