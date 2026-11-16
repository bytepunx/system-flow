---
id: T-3396
type: task
nature: remediation
title: flai accept's preview passes the queue's hold and the hand-landing flags, and its dry run says the batch runs the stale tiers
status: done
parent: S-0796
owner: alex
created: 2026-11-16T00:51:05Z
updated: 2026-11-16T01:28:10Z
transitions:
  - to: ready
    at: 2026-11-16T01:02:09Z
    by: agent-S-0796
  - to: in-progress
    at: 2026-11-16T01:02:09Z
    by: agent-S-0796
  - to: done
    at: 2026-11-16T01:28:10Z
    by: agent-S-0796
stream: S-0796
tags: [flai, acceptance, merge-queue]
touches: [flai/cmd/accept.go, flai/cmd/accept_queue_test.go]
after: [T-3395]
usage:
  source: log
  seconds: 1561
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 55
      output: 303
      cache_read: 2043186
      cache_write: 91579
      cost: 0.8087
---
# T-3396 flai accept's preview passes the queue's hold and the hand-landing flags, and its dry run says the batch runs the stale tiers

## Work

In `flai/cmd/accept.go`, where `acceptance()` builds the preview's options (line 1304, `QueueBatches` from `a.queueBatches`), set `QueueHeld` from the merge queue's state, the hold `enqueueAccepted` reads at line 460, and pass `WithoutKluster` and `Keep` where they are not passed yet, so the preview's predicate sees what `queuesInstead` sees. A held queue must keep the gate: `queuesInstead` is true then, but `enqueueAccepted` returns nothing and the story lands by hand through the same preview.

Where the dry run prints `QueueTiersWhy`, print `VerifyBehindMainWhy` the same way, so `flai accept --dry-run` says the batch runs the stale tiers.

Add a test to `flai/cmd/accept_queue_test.go`, after `TestTheOperatorsAcceptanceEnqueuesWhileTheQueueBatches` (line 391), with main moved under the story's passing verify: while the queue batches, `flai accept --dry-run --json` lists no `verify_behind_main` and gives `stale_tiers` and the reason, and `flai accept` enqueues the story; with the queue held, or with `--without-queue-tiers`, the dry run lists the blocker and the acceptance is refused before anything is merged.

It waits for T-3395, whose `QueueHeld`, predicate, and `VerifyBehindMainWhy` it sets and prints.

## Done when

- The new test passes, and so do `TestTheOperatorsAcceptanceLandsWhileTheQueueIsHeldOrWithoutTheQueueTiers` and `TestAnAcceptanceIsRefusedWhileMainOutranTheStorysVerify`.
- `flai test flai/cmd/accept.go flai/cmd/accept_queue_test.go` passes.

## Notes
