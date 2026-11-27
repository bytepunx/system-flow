---
id: T-3610
type: task
nature: remediation
title: The task plan, changes, epic follow, cancel cascade, plan and promotion candidates read states under the repository's workflow
status: done
parent: S-0843
owner: alex
created: 2026-11-18T10:35:33Z
updated: 2026-11-18T10:50:57Z
transitions:
  - to: ready
    at: 2026-11-18T10:35:52Z
    by: agent-S-0843
  - to: in-progress
    at: 2026-11-18T10:35:52Z
    by: agent-S-0843
  - to: done
    at: 2026-11-18T10:50:57Z
    by: agent-S-0843
stream: S-0843
tags: [flai, workflow]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/mcpserver/server.go, flai/cmd/items.go, flai/internal/serve/replan.go, flai/internal/inbox/inbox.go, flai/cmd/plan.go, flai/internal/inbox/starved.go, flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/kluster.go, flai/internal/serve/replan_test.go]
usage:
  source: log
  seconds: 905
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 5761
      cache_read: 1681670
      cache_write: 76145
      cost: 0.9068
---
# T-3610 The task plan, changes, epic follow, cancel cascade, plan and promotion candidates read states under the repository's workflow

## Work

T-3567's search found eight more callers of a baseline wrapper that have the repository at hand and whose result depends on a state name (I-0379):

| Caller | Baseline call | In form |
|--------|---------------|---------|
| `flai/internal/hostapi/hostapi.go` | `workitem.PlanOf` | `PlanOfIn(repo.Flow(), ...)` |
| `flai/internal/mcpserver/server.go` | `workitem.PlanOf` | `PlanOfIn` |
| `flai/cmd/items.go` | `workitem.PlanOf` | `PlanOfIn` |
| `flai/internal/serve/replan.go` | `workitem.Changes` | `ChangesIn` |
| `flai/internal/inbox/inbox.go` | `workitem.Changes` | `ChangesIn` with the `flow` it reads already |
| `flai/cmd/plan.go` | `workitem.PlanCandidatesOf` | `PlanCandidatesIn` |
| `flai/internal/inbox/starved.go` | `workitem.PromotionCandidates`, and `StarvedFrom` where it reads states | `PromotionCandidatesIn` |
| `flai/internal/check/check_items.go` | `workitem.CancelledWith`, `workitem.EpicFollows` | `CancelledWithIn`, `EpicFollowsIn` with `c.repo.Flow()` |

Call the `In` form with the repository's workflow in each. Add a test in `flai/internal/check/items_test.go` that reproduces the baseline reading for the epic follow or the cancel cascade under a renamed board, as T-3609's `TestOverlapReadsTasksUnderTheWorkflow` does; the other callers are one-line changes to functions whose `In` forms are tested in `workitem`.

Waits for no task of the story: no path in common with an open task. T-3567 waits for it.

## Done when

- No caller above calls the baseline wrapper.
- The new check test fails on main and passes on the branch.
- `flai test` on the changed paths passes.

## Notes
