---
id: TH-1295
title: "Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left"
anchor:
  path: wip/kanban/stories/S-0843-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md
  item: S-0843
status: resolved
participants: [planner-S-0843, orchestrator]
created: 2026-11-17T12:30:01Z
updated: 2026-11-17T12:35:09Z
---

# TH-1295 Plan for S-0843: four tasks in three layers pass the project's workflow to the baseline callers left

On wip/kanban/stories/S-0843-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md.

## Entries

### 2026-11-17T12:30:01Z planner-S-0843
S-0843 is planned: four tasks in three layers, 12 file touches, forecast 1h, cost of delay 37.50 USD a week. No question is open; this thread summarises the plan.

## Proposed solution

Every caller that has the project's workflow at hand and whose result depends on a state name calls the `In` form of the wrapper, with `repo.Flow()` or `Options.flow()`.

## Tasks and layers

| Layer | Task | Changes | After |
|-------|------|---------|-------|
| 1 | T-3561 | `HoldsMemo.NewHoldsIn` in `flai/internal/workitem/hold.go` | none |
| 1 | T-3563 | `CoveringIn` and the grown-claim watch in `flai/internal/itemedit` use `NewHoldsIn` | none |
| 2 | T-3565 | `flai/internal/metrics` cost order uses `NewPullerIn`, hold replay uses `memo.NewHoldsIn` | T-3561 |
| 3 | T-3567 | search again for baseline callers, close I-0379 | T-3563, T-3565 |

Each code task has a test that builds a renamed workflow and reproduces the baseline reading.

## Assumptions

- Three of the four callers the instance names are fixed on main already: `flai/cmd/guard.go` gives the guard `repo.Flow()`, `flai/internal/serve/replan.go` calls `ReplayIn`, and `hold.go` names no `InProgress`. T-3567 checks them again and names them in the close reason.
- The digest part, `digest.Decisions` on `Item.CompletedAt`, belongs to S-0842, whose `Item.Closed` fix it rides on. This story takes no `after` on S-0842, since it reads nothing S-0842 writes.
- A `NewHolds(nil, ...).Claim` on bare touches reads no task and stays as it is.
- Forecast raised from flai's 28m to 1h, for three packages' tests and a close-out of about 10m. Cost of delay value 37.50 USD a week, from flai's input, unchanged.

Nothing for the operator to do but read this. Finalizing the draft is the operator's call.

### 2026-11-17T12:34:10Z orchestrator
Approved. The plan fits S-0843.

- T-3561, T-3563, and T-3565 move the four callers still on the baseline to the `In` forms. Each has a test under a renamed workflow. That meets criterion 1.
- T-3567 searches again and closes I-0379. That meets criterion 2.
- The touches are 12 files and no folder. Each task's touches are inside them.
- The layers follow the reads: T-3565 calls T-3561's `NewHoldsIn`.
- The forecast of 1h and the value of 37.50 USD a week are flai's figures, with the planner's stated reason for the forecast.
- Leaving `digest.Decisions` to S-0842's `Item.Closed` fix keeps the two stories apart. Neither needs `after` on the other.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:35:09Z orchestrator
Resolved: The orchestrator approved S-0843's plan under plan_backlog_stories.
