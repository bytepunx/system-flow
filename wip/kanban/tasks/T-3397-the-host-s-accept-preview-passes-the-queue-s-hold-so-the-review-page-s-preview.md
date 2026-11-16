---
id: T-3397
type: task
nature: remediation
title: The host's accept.preview passes the queue's hold, so the review page's preview lists no verify_behind_main for a story Accept would enqueue
status: done
parent: S-0796
owner: alex
created: 2026-11-16T00:51:10Z
updated: 2026-11-16T01:28:18Z
transitions:
  - to: ready
    at: 2026-11-16T01:02:09Z
    by: agent-S-0796
  - to: in-progress
    at: 2026-11-16T01:02:09Z
    by: agent-S-0796
  - to: done
    at: 2026-11-16T01:28:18Z
    by: agent-S-0796
stream: S-0796
tags: [flai, hostapi, acceptance, merge-queue]
touches: [flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go]
after: [T-3395]
usage:
  source: log
  seconds: 1569
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 41
      output: 195
      cache_read: 1080608
      cache_write: 270227
      cost: 0.5117
---
# T-3397 The host's accept.preview passes the queue's hold, so the review page's preview lists no verify_behind_main for a story Accept would enqueue

## Work

In `flai/internal/hostapi/reads.go`, the `accept.preview` method (line 218) passes `QueueBatches` from `mergequeue.Batches`, which ignores the hold. Set `QueueHeld` from the merge queue's state too, so a held queue keeps the gate on the review page, where Accept would land by hand.

Add a test to `flai/internal/hostapi/reads_test.go`, after `TestAcceptPreviewListsTheQueueTiersWhileTheQueueBatches` (line 376), with main moved under the story's passing verify: while the queue batches, `accept.preview` lists no `verify_behind_main` in its blockers and gives `stale_tiers` and `verify_behind_main_why`; with the queue held, it lists the blocker. `Review.svelte` disables Accept and offers Rebase and resolve on any blocker, so this is what the review page reads; no flaiover file changes.

It waits for T-3395, whose `QueueHeld` it sets. It runs beside T-3396: the two change different files.

## Done when

- The new test passes, and so does `TestAcceptPreviewListsTheQueueTiersWhileTheQueueBatches`.
- `flai test flai/internal/hostapi/reads.go flai/internal/hostapi/reads_test.go` passes, the `flaiover-flai` tier with it.

## Notes
