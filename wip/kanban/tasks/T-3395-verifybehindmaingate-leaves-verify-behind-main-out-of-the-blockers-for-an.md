---
id: T-3395
type: task
nature: remediation
title: verifyBehindMainGate leaves verify_behind_main out of the blockers for an acceptance that would enqueue, and keeps the stale tiers and a reason
status: done
parent: S-0796
owner: alex
created: 2026-11-16T00:50:54Z
updated: 2026-11-16T01:01:55Z
transitions:
  - to: ready
    at: 2026-11-16T00:53:16Z
    by: agent-S-0796
  - to: in-progress
    at: 2026-11-16T00:53:16Z
    by: agent-S-0796
  - to: done
    at: 2026-11-16T01:01:55Z
    by: agent-S-0796
stream: S-0796
tags: [flai, acceptance, merge-queue]
touches: [flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go]
usage:
  source: log
  seconds: 519
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 38
      output: 250
      cache_read: 1294026
      cache_write: 138891
      cost: 0.5428
---
# T-3395 verifyBehindMainGate leaves verify_behind_main out of the blockers for an acceptance that would enqueue, and keeps the stale tiers and a reason

## Work

In `flai/internal/preview`:

- Give `AcceptOptions` (`orchestrator.go`, line 61) a `QueueHeld bool`: the merge queue is held, so an acceptance lands by hand though the queue batches. `mergequeue.Batches` ignores the hold, so the callers must say it.
- Add one predicate beside `verifyBehindMainGate` (`accept.go`, line 544), such as `enqueues(it, opts)`, true when the item is a story whose state accepts, `opts.QueueBatches`, not `opts.Orchestrator`, `opts.Batch` empty, not `opts.WithoutQueueTiers`, not `opts.WithoutKluster`, `opts.Keep` empty, and not `opts.QueueHeld`. It mirrors `queuesInstead` in `flai/cmd/accept.go` (line 435) but for the dry run, which the preview is; say so in its comment.
- When it is true, `verifyBehindMainGate` still sets `res.StaleTiers`, sets a new `VerifyBehindMainWhy` (JSON `verify_behind_main_why`, after `QueueTiersWhy`) saying the merge queue's batch runs those tiers over the main it lands on, and adds no blocker and no `VerifyBehindMainBlocker`. Otherwise it is unchanged.
- Extend `TestAnAcceptanceWaitsForAVerifyOverMainsNewPaths` (`accept_test.go`, line 644): with `QueueBatches` alone the preview lists no `verify_behind_main`, keeps `stale_tiers`, and gives the reason; with `QueueHeld`, `WithoutQueueTiers`, `WithoutKluster`, `Keep`, or `Orchestrator` added, it lists the blocker as before. The first case fails before the change: it reproduces I-0381.

This task waits for none: it changes only the preview package, which the two caller tasks then wire.

## Done when

- The new cases pass and the existing cases of the test still pass.
- `flai test flai/internal/preview` passes.

## Notes
