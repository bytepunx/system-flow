---
id: S-0796
type: story
nature: remediation
title: The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue
status: done
owner: alex
created: 2026-11-16T00:41:47Z
updated: 2026-11-16T02:24:38Z
transitions:
  - to: ready
    at: 2026-11-16T00:51:16Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-16T00:51:28Z
    by: agent-S-0796
  - to: review
    at: 2026-11-16T02:07:11Z
    by: agent-S-0796
  - to: done
    at: 2026-11-16T02:24:38Z
    by: orchestrator
tags: [flai, acceptance, merge-queue]
topics: [acceptance, merge-queue]
touches: [flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_queue_test.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, design/adrs/0308-while-the-merge-queue-batches-and-is-not-held-the-acceptance-preview-lists.md, design/adrs/README.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0381-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4557
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 2
      work: 146
  models:
    - model: claude-opus-5-5
      input: 492
      output: 83079
      cache_read: 45898617
      cache_write: 885773
      cost: 16.192
  strategic:
    - kind: planner
      seconds: 362
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 38
          output: 24601
          cache_read: 1160717
          cache_write: 102402
          cost: 0.0671
        - model: claude-opus-5-5
          input: 42
          output: 20287
          cache_read: 2632495
          cache_write: 169866
          cost: 2.2913
    - kind: orchestrator
      seconds: 180
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 49
          output: 469
          cache_read: 1717019
          cache_write: 289448
          cost: 0.61
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 13566
          cache_write: 29666
          cost: 0.0343
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h
    by: flai
    at: 2026-11-16T00:41:47Z
  value: 300
  by: planner-S-0796
  at: 2026-11-16T00:50:28Z
forecast:
  duration: 45m
  delivery: 2026-11-16T03:45:00Z
  basis: flai forecast's 27m (116 s per unit over 73 large-band remediation stories, size 14) raised to 45m because the preview must newly read the merge queue's hold and the close-out of a flai/ story runs integration and smoke, about twelve minutes
  by: planner-S-0796
  at: 2026-11-16T00:50:28Z
finalized:
  by: orchestrator
  at: 2026-11-16T00:51:09Z
---
# S-0796 The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue

## Goal

This story remediates [I-0381](../../../design/issues/I-0381-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md), "The acceptance preview lists verify_behind_main while the merge queue batches, so the review page disables Accept and offers Rebase and resolve for a story the acceptance would enqueue". The issue recommends this solution:

While the merge queue batches and is not held, the operator's acceptance preview does not list `verify_behind_main` as a blocker: `verifyBehindMainGate` in `flai/internal/preview/accept.go` skips when `opts.QueueBatches` is set and the acceptance would enqueue the story, as `queuesInstead` decides, and still reports the stale tiers so the preview can say the batch will run them. The review page then offers Accept, which enqueues, in place of Rebase and resolve. An acceptance that lands by hand (`--without-queue-tiers`, `--without-kluster`, `--keep-removal`, or a held queue) keeps the gate. Refine ADR-0291 with an ADR that says so.

## Acceptance criteria
- [x] The cause I-0381 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0381 is closed with `flai issue close I-0381 --reason` saying what fixed it

## Tasks
- T-3394 Write ADR-0308 refining ADR-0291: while the merge queue batches and is not held, the preview lists verify_behind_main only for an acceptance that lands by hand
- T-3395 verifyBehindMainGate leaves verify_behind_main out of the blockers for an acceptance that would enqueue, and keeps the stale tiers and a reason
- T-3396 flai accept's preview passes the queue's hold and the hand-landing flags, and its dry run says the batch runs the stale tiers
- T-3397 The host's accept.preview passes the queue's hold, so the review page's preview lists no verify_behind_main for a story Accept would enqueue
- T-3398 The design and the guides say verify_behind_main holds only an acceptance that lands by hand, and I-0381 is closed

## Notes

Cost of delay inputs set by flai from I-0381. time_lost_per_cycle 2h carried over from I-0381's Impact section. Its Impact time_lost_per_cycle was taken rather than the 1h derived from its cost and count.

### Planning

Touches, by where each came from (the story declared none):

- Design (the goal and I-0381's instance): `flai/internal/preview/accept.go` (`verifyBehindMainGate`, line 544), `flai/cmd/accept.go` (`queuesInstead`, line 435, and the preview options `acceptance()` builds, line 1304), `design/issues/I-0381-…md` (closed by criterion 2).
- Layout (the explorer's read of the code): `flai/internal/preview/orchestrator.go` holds `AcceptOptions`, which gains the queue's hold; `flai/internal/hostapi/reads.go` serves `accept.preview` to the review page with `QueueBatches` and no hold; the tests beside each, `flai/internal/preview/accept_test.go` (`TestAnAcceptanceWaitsForAVerifyOverMainsNewPaths`), `flai/cmd/accept_queue_test.go`, and `flai/internal/hostapi/reads_test.go`.
- Design (the ADR the goal asks for and the documents that state ADR-0291's gate): `design/adrs/0308-…md`, the next free number, its slug from the title the task proposes; `design/system/workflow.md` § What waits on an acceptance; `design/system/flai-cli.md` § The acceptance gate; `docs/users/flai.md` § Accept and release.
- Co-change (`flai touches suggest`): `design/adrs/README.md` (14%), the index `flai adr new` writes; `docs/users/flaiover.md`, whose Rebase and resolve paragraph names "a verify that no longer stands" as a hold. Left out: the other co-change paths, which this change does not reach.
- No flaiover file: `Review.svelte` disables Accept on any blocker in the preview, so the page offers Accept once flai lists none.

No `after`: the story reads nothing another open story writes.

Figures:

- Forecast 45m, raised from `flai forecast`'s 27m: the preview must read the merge queue's hold, which neither it nor `accept.preview` reads today, and a story under `flai/` closes out with the integration and smoke tiers, about twelve minutes (ADR-0291 Consequences). Delivery moved by the same 18 minutes.
- Cost of delay 300 USD a week, as `flai cod` gives it from the inputs flai carried from I-0381 (2h a cycle at 150 USD an hour, one cycle a week). It stands: the inputs are the operator's, and I-0381 records three stories held on one day, so the figure is not high.

### Accepted by the orchestrator

- Verified: 9d2deecbbad7c46c04668c1fc96ed9623117032f
- Batch: 20261116-021146
- Enqueued verified: 9d2deecbbad7c46c04668c1fc96ed9623117032f
- At: 2026-11-16T02:24:38Z

Verdict: meets both criteria; verified at 9d2deecbbad7c46c04668c1fc96ed9623117032f, close-out passed every step
- 1: flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/cmd/accept.go, flai/internal/hostapi/reads.go, flai/internal/preview/accept_test.go, flai/cmd/accept_queue_test.go, flai/internal/hostapi/reads_test.go, design/adrs/0308-while-the-merge-queue-batches-and-is-not-held-the-acceptance-preview-lists.md
- 2: design/issues/I-0381-the-acceptance-preview-lists-verify-behind-main-while-the-merge-queue-batches.md, design/issues/summary.md
