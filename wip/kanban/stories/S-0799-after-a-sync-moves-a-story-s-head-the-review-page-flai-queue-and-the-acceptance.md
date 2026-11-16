---
id: S-0799
type: story
nature: remediation
title: After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head
status: ready
owner: alex
created: 2026-11-16T00:41:49Z
updated: 2026-11-16T01:53:30Z
transitions:
  - to: ready
    at: 2026-11-16T01:51:03Z
    by: orchestrator
tags: [flai, flaiover, merge-queue, verify]
topics: [acceptance, merge-queue, verify]
touches: [flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/cmd/queue.go, flai/cmd/queue_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json, design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0384-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 80
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 4911
          cache_read: 1436086
          cache_write: 34900
          cost: 0.6647
    - kind: orchestrator
      seconds: 120
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 31
          output: 344
          cache_read: 1311447
          cache_write: 134631
          cost: 0.4402
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-16T00:41:49Z
  value: 50
  by: planner-S-0799
  at: 2026-11-16T01:47:47Z
forecast:
  duration: 1h15m
  delivery: 2026-11-16T03:30:00Z
  basis: "Its own forecast of 1h15m; 2nd in the pull order with an in-progress limit of 3, behind S-0750, S-0796, S-0800 and S-0798."
  by: flai
  at: 2026-11-16T01:53:30Z
finalized:
  by: orchestrator
  at: 2026-11-16T01:50:51Z
---
# S-0799 After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head

## Goal

This story remediates [I-0384](../../../design/issues/I-0384-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md), "After a sync moves a story's head, the review page, flai queue, and the acceptance preview report the verify failure from before the sync while a new close-out runs at the new head". The issue recommends this solution:

While a close-out job runs for a story at a head other than its last report's, `flai queue`, the acceptance preview, and the review page say so, naming the job and its step, in place of or beside the last finished report's failure: the job's state is in `.flai-cache/verify/jobs/<S-nnnn>.job.json`. A `verify_failed` whose report's commit is not the branch's head says the report is from before the head moved.

## Acceptance criteria
- [ ] The cause I-0384 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0384 is closed with `flai issue close I-0384 --reason` saying what fixed it

## Tasks
- T-3409 A running close-out job records the head it runs at, and closeout reads a story's running job with that head and the step it is on
- T-3414 verify reads a story's last report against its branch head, saying when the report is from before the head moved and which close-out job runs at the new head
- T-3415 flai queue names the close-out job running at a story's new head, and says a verify_failed report is from before the head moved
- T-3416 The acceptance preview says a close-out runs at the story's new head, naming the job and its step, beside a last report from before the head moved
- T-3417 flai verify --last and the host's verify.status say a report is from before the head moved and name the close-out job running at the head
- T-3418 The review page and the merge queue panel show the close-out running at a story's new head, its job and step, beside a last report from before the head moved
- T-3419 The design and the guides describe the running close-out's head and step, the verify_running reason, and a report from before the head moved, and I-0384 is closed

## Notes

Cost of delay inputs set by flai from I-0384. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-16T00:41:30Z, 0 days before this story; under one cycle counts as one).

### Planning

Touches, by where each came from (the story declared none):

- Design (the goal and I-0384): `flai/internal/closeout/closeout.go`, the `Job` the job file holds, which records no head and no step while it runs; `flai/internal/serve/verifyjob.go`, which writes it; `design/issues/I-0384-…md` and `design/issues/summary.md`, closed by criterion 2.
- Layout (the explorer's read of the code): `flai/internal/verify/story.go` (`Report`, `LastReport`, `QueueTiersNotRun`); `flai/internal/mergequeue/eligible.go` (`readLastVerify`, the `verify_failed` reason) and `flai/cmd/queue.go`; `flai/internal/preview/accept.go` (`queue_tiers`, `verify_behind_main`); `flai/cmd/verify.go` (`--last`) and `flai/internal/hostapi/writes_host.go` (`verify.status`, which the review page's Verification row reads); `flaiover/src/lib/review.ts`, `flaiover/src/lib/components/Review.svelte`, and `flaiover/src/lib/components/MergeQueue.svelte`; and the test file beside each.
- Layout (the component examples rule): `design/components/flaiover/src/lib/components/MergeQueue.examples.json`, since `QueueReason` gains a field.
- Design (the documents that state what changes): `design/system/flai-cli.md` § The close-out job; `docs/users/flai.md` § Accept and release › The merge queue and § Verify a story before review; `docs/users/flaiover.md` § Board › Reviewing a story.
- Co-change (`flai touches suggest`): the paths it lists are each under 13% and none is reached by this change. `design/system/flaiover-dashboard.md` (12%) has no section on the review page's Verification row, and `docs/users/flai-reference.md` (9%) is generated from flags, which this story does not change. Left out.

No folder touch: every task names its files.

No `after`: the story reads nothing another open story writes. S-0796, in progress, changes `flai/internal/preview/accept.go`, `design/system/flai-cli.md`, and `docs/users/flai.md`, and S-0798 changes the close-out. A trial merge judges those overlaps (ADR-0230), and T-3416 and T-3419 say to sync first.

Figures:

- Forecast 1h15m, raised from `flai forecast`'s 58m: the job must newly record the head it runs at and the step it is on, which nothing reads today, and the change spans flai and flaiover, so the close-out runs the flaiover tiers, svelte-check among them, beside integration and smoke, about twelve minutes. Delivery moved by the same 17 minutes.
- Cost of delay 50 USD a week, as `flai cod` gives it from the inputs flai carried from I-0384 (20m a cycle at 150 USD an hour, one cycle a week). It stands: the inputs are the operator's, and I-0384 records one occurrence.
