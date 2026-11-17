---
id: S-0843
type: story
nature: remediation
title: Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them
status: backlog
owner: alex
created: 2026-11-17T08:32:54Z
updated: 2026-11-17T12:35:58Z
transitions: []
tags: []
topics: [workflow]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/itemedit/covers.go, flai/internal/itemedit/covers_test.go, flai/internal/itemedit/claim.go, flai/internal/itemedit/claim_test.go, flai/internal/metrics/costorder.go, flai/internal/metrics/costorder_test.go, flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go, design/issues/I-0379-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md, design/issues/summary.md]
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
      seconds: 14
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 24
          output: 89
          cache_read: 1509392
          cache_write: 40816
          cost: 0.1059
        - model: claude-opus-5-5
          input: 31
          output: 2391
          cache_read: 2599588
          cache_write: 40766
          cost: 0.8704
    - kind: orchestrator
      seconds: 34
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 44
          cache_read: 423129
          cache_write: 15722
          cost: 0.1373
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-17T08:32:54Z
  value: 37.5
  by: planner-S-0843
  at: 2026-11-17T12:29:46Z
forecast:
  duration: 1h
  delivery: 2026-11-17T21:53:00Z
  basis: "Its own forecast of 1h; 21st in the pull order with an in-progress limit of 3, behind S-0806, S-0834, S-0811, S-0838, S-0829, S-0840, S-0809, S-0842, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815, S-0822, S-0824, S-0835, S-0837, S-0839 and S-0841."
  by: flai
  at: 2026-11-17T12:35:58Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:34Z
---
# S-0843 Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them

## Goal

This story remediates [I-0379](../../../design/issues/I-0379-callers-of-the-workflow-aware-guard-forecast-replay-and-share-hold-still-pass.md), "Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0379 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0379 is closed with `flai issue close I-0379 --reason` saying what fixed it

## Tasks
- T-3561 HoldsMemo builds its Holds under a workflow it is given, as NewHoldsIn does
- T-3563 itemedit's covering stories and grown-claim watch read claims under the project's workflow
- T-3565 flai stats' cost order and hold replay play out under the report's workflow
- T-3567 Confirm no caller with a workflow at hand runs a baseline wrapper, and close I-0379

## Notes

Cost of delay inputs set by flai from I-0379. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T20:04:03Z, 1.5 days before this story; under one cycle counts as one).

### Planning

Proposed solution: pass the project's workflow, `repo.Flow()` or `Options.flow()`, to every caller of a baseline wrapper that has one at hand and whose result depends on a state name.

What main holds already, read on 2026-11-17: three of the four callers the instance names were fixed by later stories. `flai/cmd/guard.go` gives the guard `repo.Flow()` for the orchestrator and for every role a role file holds, the planner among them. `flai/internal/serve/replan.go` calls `planning.ReplayIn`. `flai/internal/workitem/hold.go` names a holding story's state by the workflow's labels, with no `InProgress` left. What still runs the baseline:

| Caller | Baseline call | Task |
|--------|---------------|------|
| `flai/internal/metrics/costorder.go` | `planning.NewPuller` | T-3565 |
| `flai/internal/metrics/claims.go` | `HoldsMemo.NewHolds` | T-3561, T-3565 |
| `flai/internal/itemedit/covers.go` | `workitem.NewHolds` inside `CoveringIn(flow, ...)` | T-3563 |
| `flai/internal/itemedit/claim.go` | `workitem.NewHolds` | T-3563 |

The instance's last sentence, `digest.Decisions` filtering on `Item.CompletedAt`, follows from `Item.Closed` reading the baseline's names, which is I-0377 and S-0842's to fix; this story does not touch `flai/internal/digest/sources.go`. No `after` on S-0842: this story reads nothing S-0842 writes.

Callers of `NewHolds(nil, ...).Claim(&workitem.Item{Touches: ...})` in `taskdone`, `widen`, `context`, and `metrics/claims.go` read no task, so their claim is the same under every workflow; they stay.

Touches, all files, no folder:

- design: the issue's instance names `hold.go`, and `flai/internal/workitem/hold_test.go` is its test.
- layout: the other code files are where the baseline calls are, found by search on main; each `_test.go` beside them is where the reproducing test goes.
- layout: `design/issues/I-0379-...md` and `design/issues/summary.md`, which `flai issue close` writes.
- co-change: `flai touches suggest` listed nothing above 4% of the commits that changed these files (`design/system/flai-cli.md`, `docs/users/flai.md`, other issues). None is taken: no command or user-facing behaviour changes.

Tasks and layers: T-3561 and T-3563 first, together, no path in common; T-3565 after T-3561, whose `HoldsMemo.NewHoldsIn` it calls; T-3567 last, after T-3563 and T-3565, to search with their changes in and close the issue.

Forecast: flai gave 28m, 119 s per unit of size over 82 done remediation stories, size 14. Set to 1h: four tasks in three packages, each needing a test that builds a renamed workflow, and a close-out that took 9m43s on S-0828 today. Delivery moved by the same 32m, to 2026-11-17T21:13:00Z.

Cost of delay: 37.50 USD a week, as `flai cod` worked it out from flai's input, 15m lost per 168h cycle at 150 USD an hour. It stands: no renamed board runs yet, so nothing is lost until one does.
