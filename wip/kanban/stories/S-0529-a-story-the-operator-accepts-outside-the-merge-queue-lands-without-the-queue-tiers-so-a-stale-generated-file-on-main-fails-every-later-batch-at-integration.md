---
id: S-0529
type: story
nature: remediation
title: A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration
status: backlog
owner: alex
created: 2026-10-28T12:54:24Z
updated: 2026-10-28T18:28:28Z
transitions: []
tags: [cli]
topics: [orchestration, testing]
touches: [design/adrs, design/adrs/README.md, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/accept.go, flai/cmd/verify.go, flai/cmd/accept_queue_tiers_test.go, flai/cmd/move_refused_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, design/system/workflow.md, design/system/flai-cli.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 164
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 520
          cache_read: 3117542
          cache_write: 27928
          cost: 0.691
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
  delivery: 2026-10-29T10:30:00Z
  basis: "Its own forecast of 1h30m; 38th in the pull order with an in-progress limit of 5, behind S-0415, S-0502, S-0531, S-0419, S-0423, S-0418, S-0421, S-0416, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517, S-0518, S-0519, S-0520, S-0523, S-0524, S-0525 and S-0527."
  by: flai
  at: 2026-10-28T18:28:28Z
finalized:
  by: alex
  at: 2026-10-28T12:57:38Z
---
# S-0529 A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration

## Goal

This story remediates [I-0203](../../../design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md), "A story the operator accepts outside the merge queue lands without the queue tiers, so a stale generated file on main fails every later batch at integration". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0203 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0203 is closed with `flai issue close I-0203 --reason` saying what fixed it

## Tasks
- T-2261 An ADR refining ADR-0174 records that a hand acceptance while the queue can batch is refused until the queue tiers ran at the story's head
- T-2262 flai verify --all reuses the tiers a run that left the queue tiers passed at the same patch, and names the queue tiers a report did not run
- T-2263 flai accept by hand while the queue can batch refuses a story whose queue tiers did not run at its head, lists it as a blocker in --dry-run, and lands it with an override flag
- T-2264 The docs say a hand acceptance while the queue can batch waits for the queue tiers at the story's head, and I-0203 is closed

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
