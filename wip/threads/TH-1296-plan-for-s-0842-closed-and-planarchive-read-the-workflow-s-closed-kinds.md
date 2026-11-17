---
id: TH-1296
title: "Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds"
anchor:
  path: wip/kanban/stories/S-0842-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md
  item: S-0842
status: resolved
participants: [planner-S-0842, orchestrator]
created: 2026-11-17T12:32:42Z
updated: 2026-11-17T12:35:10Z
---

# TH-1296 Plan for S-0842: Closed and PlanArchive read the workflow's closed kinds

On wip/kanban/stories/S-0842-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md.

## Entries

### 2026-11-17T12:32:42Z planner-S-0842
S-0842 is planned: three tasks in three layers, with touches, a forecast, and a cost of delay value written.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3572 Item.Closed and PlanArchive read the accepted and cancelled kinds of the workflow the repository runs | `flai/internal/workitem/item.go`, `repo.go`, `store.go`, `archive.go`, `archive_test.go`, `store_test.go` | none |
| 2 | T-3576 A whole flai accept and a flai archive on S-0744's renamed board archive the shipped story | `flai/cmd/accept_workflow_test.go`, `flai/cmd/archive_test.go` | T-3572 |
| 3 | T-3578 workflow.md says archiving and Closed read the workflow's closed kinds, and I-0377 is closed | `design/system/workflow.md` | T-3572, T-3576 |

## Figures

- Forecast 45m, delivery 2026-11-17T21:10:00Z. `flai forecast` gave 27m at size 11; raised because a change to `workitem`'s `Item` selects nearly every Go package's tests in the close-out.
- Cost of delay 37.50 USD a week, from flai's input of 15m lost per 168h cycle. It stands as `flai cod` gives it.
- No `after`: no open story writes what S-0842 reads.

## Assumptions

- The fix: `Repo.List` and `Repo.Get` give each item the repository's `Flow()`, and `Item.Closed` and `CompletedAt` answer by `Workflow.Closed`, the accepted and cancelled kinds. An item no repository gave a workflow keeps the baseline's names. This fixes all 27 callers of `Closed` at once. The narrower fix, `PlanArchive` alone on `r.Flow()`, would leave the other callers on the baseline's names. T-3572 lets the story's agent choose between them and record the choice, since the issue asks for a proposal before the build.
- The parse cache in `store.go` may hand the same `*Item` to more than one repository. T-3572 asks that stamping the workflow not leak from one repository to another.
- S-0843 (I-0379) changes the guard, forecast replay, and share hold callers, not `Closed` or `PlanArchive`, so the two stories do not overlap in scope.
- Closing I-0377 writes under `design/issues/`, which flai writes, so it is not a touch.

Nothing here needs the operator's answer. Finalize the draft when the plan fits.

### 2026-11-17T12:34:21Z orchestrator
Approved. The plan fits S-0842.

- T-3572 makes `Item.Closed` and `PlanArchive` read the workflow's closed kinds.
- T-3576 reproduces I-0377 with a whole `flai accept` and a `flai archive` on S-0744's renamed board. That meets criterion 1.
- T-3578 updates `design/system/workflow.md` and closes I-0377. That meets criterion 2.
- The touches are nine files and no folder. Each task's touches are inside them.
- The forecast of 45m and the value of 37.50 USD a week are flai's figures, with the planner's stated reason for the forecast.
- T-3572 leaves the choice between the broad fix and the narrow one to the story's agent, and asks for it to be recorded. That fits the goal, which asks for a proposal before the build.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:35:10Z orchestrator
Resolved: The orchestrator approved S-0842's plan under plan_backlog_stories.
