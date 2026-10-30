---
id: S-0529
type: story
nature: remediation
title: A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration
status: done
owner: alex
created: 2026-10-28T12:54:24Z
updated: 2026-10-30T19:19:52Z
transitions:
  - to: ready
    at: 2026-10-30T00:10:24Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-30T00:15:24Z
    by: agent-S-0529
  - to: review
    at: 2026-10-30T01:06:06Z
    by: agent-S-0529
  - to: in-progress
    at: 2026-10-30T18:30:48Z
    by: flai
  - to: review
    at: 2026-10-30T19:04:05Z
    by: agent-S-0529
  - to: done
    at: 2026-10-30T19:19:52Z
    by: orchestrator
tags: [cli]
topics: [orchestration, testing]
touches: [design/adrs, design/adrs/README.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/accept.go, flai/cmd/verify.go, flai/cmd/accept_queue_tiers_test.go, flai/cmd/move_refused_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, design/system/workflow.md, design/system/flai-cli.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md, docs/operators/settings.md, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flai/internal/preview/orchestrator.go, design/issues/I-0221-the-close-out-s-go-test-tier-runs-only-the-changed-packages-so-a-change-that-breaks-a-caller-package-s-tests-passes-before-review.md, flai/internal/mcpserver/verify_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5043
  priced_by: harness
  turns:
    - day: 2026-10-30
      ceremony: 5
      hand_edits: 7
      work: 71
  models:
    - model: claude-opus-5-5
      input: 414
      output: 125588
      cache_read: 33529801
      cache_write: 1384074
      cost: 18.2027
  strategic:
    - kind: orchestrator
      seconds: 370
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 97
          output: 1037
          cache_read: 5477924
          cache_write: 356163
          cost: 1.2895
        - model: claude-sonnet-5-5
          input: 18
          output: 93
          cache_read: 228943
          cache_write: 106889
          cost: 0.2332
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:54:24Z
  value: 37.5
  by: planner-S-0529
  at: 2026-10-28T18:16:03Z
forecast:
  duration: 1h30m
  delivery: 2026-10-30T02:58:00Z
  basis: "Its own forecast of 1h30m; 5th in the pull order with an in-progress limit of 5, behind S-0419, S-0511, S-0527, S-0424, S-0523, S-0524 and S-0425."
  by: flai
  at: 2026-10-30T00:13:09Z
finalized:
  by: alex
  at: 2026-10-28T12:57:38Z
---
# S-0529 A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Goal

This story remediates [I-0203](../../../design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md), "A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0203 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0203 is closed with `flai issue close I-0203 --reason` saying what fixed it

## Tasks
- T-2261 An ADR refining ADR-0174 records that a hand acceptance while the queue can batch is refused until the queue tiers ran at the story's head
- T-2262 flai verify --all reuses the tiers a run that left the queue tiers passed at the same patch, and names the queue tiers a report did not run
- T-2263 flai accept by hand while the queue can batch refuses a story whose queue tiers did not run at its head, lists it as a blocker in --dry-run, and lands it with an override flag
- T-2264 The docs say a hand acceptance while the queue can batch waits for the queue tiers at the story's head, and I-0203 is closed
- T-2345 TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch expects verify --all after a run that left the queue tiers to reuse the tiers that run passed, as ADR-0221 makes it
- T-2346 ADR-0221 is accepted, as TH-0701's confirmed recommendation settles the fix

## Notes

Cost of delay inputs set by flai from I-0203. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-24T00:15:33Z, 4.5 days before this story; under one cycle counts as one).

### Planning

Proposed fix (an assumption until the plan's thread settles it): while the merge queue can batch ([ADR-0174](../../../design/adrs/0174-while-the-merge-queue-can-batch-flai-verify-and-the-close-out-leave-the-test.md)), `flai accept` by hand, the operator's from the CLI, the board, or the dashboard, refuses a story whose last `flai verify` at its branch's head did not run each tier marked `queue` that its paths select, naming the tiers and `flai verify <S-nnnn> --all`; `--dry-run` and the dashboard's preview list it as a blocker; one flag lets the operator land it anyway, recorded in the acceptance. `flai verify --all` reuses the tiers a close-out at the same patch passed, so it runs only the queue tiers. A new ADR refines ADR-0174's consequence that a hand acceptance runs no tier.

Touches, and where each came from (`flai touches suggest` had nothing to start from: the story declared none):

- `flai/internal/verify/story.go`, `story_test.go`: layout. The report's fingerprint holds the names of the tiers left to the queue, so a `--all` run reuses nothing after a close-out; the reuse rule and a helper that names the queue tiers a report did not run go here.
- `flai/cmd/accept.go`, `flai/cmd/verify.go`, `flai/internal/preview/accept.go`, `accept_test.go`, `flai/cmd/accept_queue_tiers_test.go` (new), `flai/cmd/move_refused_test.go`: layout. `acceptItem` and `preview.AcceptWith` compute the blockers; `queueBatches` is in `verify.go`; `flai move <S> done` and the host API's `item.move` and `accept.run` reach `acceptItem`.
- `docs/users/flai-reference.md`: co-change. It is generated from `flai accept`'s help, which gains the flag, and `TestReferenceIsCurrent` fails while it is stale.
- `docs/users/flai.md`, `docs/operators/index.md`, `design/system/workflow.md`, `design/system/flai-cli.md`: design. Each says a hand acceptance runs no tier, or describes the queue tiers.
- `design/adrs` (folder) and `design/adrs/README.md`: design. The new ADR's file name is not known until `flai adr new` numbers it, so the folder stays; it lies inside `claims.shared`, so it holds no ready story.
- The I-0203 file and `design/issues/summary.md`: criterion 2, `flai issue close`.

Figures:

- Forecast 1h30m, adjusted: `flai forecast` gave 32m (99 s a unit of size over 39 done remediation stories in the large band, size 19). Four tasks, two of them Go changes with real-git tests in `verify` and `accept`, and a close-out whose go-test tier runs without `-short`, take longer than the band's median. The delivery, 2026-10-29T08:45:00Z, is flai's: 37th in the pull order.
- Cost of delay value 37.50 USD a week: `flai cod` from the operator's input, 15m lost per 168h cycle at 150 USD an hour. It stands: one occurrence in 4.5 days gives no reason to scale it.
- 2026-10-30T18:30:48Z: moved to in-progress: merge queue batch 20261030-182125 fails with S-0529 at a1fc087d13fe: integration failed: flai/internal/mcpserver/verify_test.go:275 TestVerifyLeavesTheQueueTiersToTheMergeQueueWhileItCanBatch: all: passed true, steps rebase=passed sync=passed narrative=passed check=passed ok=reused late=passed, want rebase=passed sync=passed narrative=passed check=passed ok=passed late=passed

### Accepted by the orchestrator

- Verified: 3ef6d4b45261f2c8499e082ef33f796008c3c1fe
- Batch: 20261030-190545
- Enqueued verified: 3ef6d4b45261f2c8499e082ef33f796008c3c1fe
- At: 2026-10-30T19:19:52Z

Verdict: both criteria met at 3ef6d4b; flai verify passed every step at the head, integration and smoke included.
- 1: flai/cmd/accept.go, flai/internal/preview/accept.go, flai/internal/verify/story.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/reads.go, flai/internal/preview/orchestrator.go, flai/cmd/accept_queue_tiers_test.go, flai/cmd/move_refused_test.go, flai/internal/preview/accept_test.go, flai/internal/verify/story_test.go, flai/internal/hostapi/reads_test.go, flai/internal/mcpserver/verify_test.go, design/adrs/0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md, design/adrs/README.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/index.md, docs/operators/settings.md, design/system/workflow.md, design/system/flai-cli.md
- 2: design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md, design/issues/I-0221-the-close-out-s-go-test-tier-runs-only-the-changed-packages-so-a-change-that-breaks-a-caller-package-s-tests-passes-before-review.md
