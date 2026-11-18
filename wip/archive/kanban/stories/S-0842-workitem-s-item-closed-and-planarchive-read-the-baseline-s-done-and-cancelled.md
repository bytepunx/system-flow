---
id: S-0842
type: story
nature: remediation
title: workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it
status: done
owner: alex
created: 2026-11-17T08:32:53Z
updated: 2026-11-18T08:31:06Z
transitions:
  - to: ready
    at: 2026-11-17T12:35:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-18T07:31:16Z
    by: agent-S-0842
  - to: review
    at: 2026-11-18T08:29:42Z
    by: agent-S-0842
  - to: done
    at: 2026-11-18T08:31:06Z
    by: orchestrator
tags: []
topics: [workflow]
touches: [flai/internal/workitem/item.go, flai/internal/workitem/repo.go, flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go, flai/internal/workitem/store_test.go, flai/cmd/accept_workflow_test.go, flai/cmd/archive_test.go, design/system/workflow.md, flai/internal/workitem/view.go, design/issues/I-0377-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3518
  priced_by: harness
  turns:
    - day: 2026-11-18
      ceremony: 1
      test_runs: 2
      work: 62
  models:
    - model: claude-opus-5-5
      input: 232
      output: 45285
      cache_read: 18419735
      cache_write: 552585
      cost: 8.172
  strategic:
    - kind: planner
      seconds: 161
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 18
          output: 60
          cache_read: 1479862
          cache_write: 41705
          cost: 0.104
        - model: claude-opus-5-5
          input: 40
          output: 8280
          cache_read: 3023647
          cache_write: 99552
          cost: 1.409
    - kind: orchestrator
      seconds: 69
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 209
          cache_read: 731165
          cache_write: 75391
          cost: 0.2534
        - model: claude-sonnet-5-5
          input: 4
          output: 40
          cache_read: 14096
          cache_write: 31564
          cost: 0.0369
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-17T08:32:53Z
  value: 37.5
  by: planner-S-0842
  at: 2026-11-17T12:29:22Z
forecast:
  duration: 45m
  delivery: 2026-11-17T15:01:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 3, behind S-0806, S-0834, S-0840, S-0811 and S-0809."
  by: flai
  at: 2026-11-17T13:23:37Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:35:33Z
---
# S-0842 workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it

## Goal

This story remediates [I-0377](../../../design/issues/I-0377-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md), "workitem's Item.Closed and PlanArchive read the baseline's done and cancelled names, so on a board that renames them flai accept moves a story to the accepted state and then cannot archive it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0377 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0377 is closed with `flai issue close I-0377 --reason` saying what fixed it

## Tasks
- T-3572 Item.Closed and PlanArchive read the accepted and cancelled kinds of the workflow the repository runs
- T-3576 A whole flai accept and a flai archive on S-0744's renamed board archive the shipped story
- T-3578 workflow.md says archiving and Closed read the workflow's closed kinds, and I-0377 is closed

## Notes

Cost of delay inputs set by flai from I-0377. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T01:36:52Z, 2.3 days before this story; under one cycle counts as one).

### Planning

Proposed fix, for the story's agent to confirm or change when it pulls the story: `Item.Closed` answers by the workflow its repository runs. `Repo.List` and `Repo.Get` give each item they return the repository's `Flow()`, and `Closed`, and `CompletedAt` through it, read `Workflow.Closed` (`flai/internal/workflow/review.go`), the accepted and cancelled kinds; an item no repository gave a workflow, such as one a test parses, keeps the baseline's names. This fixes the 27 callers of `Closed` outside tests at once, `PlanArchive` among them, without changing each. `PlanArchive`'s refusal names the workflow's accepted and cancelled states in place of "done or cancelled". The narrower alternative, `PlanArchive` reading `r.Flow().Closed` and an `Item.ClosedIn(w)` for the callers that hold a workflow, leaves the rest on the baseline's names.

Touches, all files, no folder:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/workitem/archive.go` | design (I-0377) | `PlanArchive`, the refusal the issue quotes |
| `flai/internal/workitem/item.go` | design (I-0377) | `Item.Closed` and `CompletedAt` |
| `flai/internal/workitem/repo.go` | layout | `List` and `Get`, where items get the repository's workflow |
| `flai/internal/workitem/store.go` | layout | the parse cache that hands out the items `List` and `Get` return |
| `flai/internal/workitem/archive_test.go` | layout | `PlanArchive` on a renamed workflow |
| `flai/internal/workitem/store_test.go` | layout | items from the cache carry the workflow |
| `flai/cmd/accept_workflow_test.go` | layout | the reproduction: a whole `flai accept` on S-0744's renamed board |
| `flai/cmd/archive_test.go` | layout | `flai archive` on the renamed board |
| `design/system/workflow.md` | design | § Transitions and who makes them says what reads the closed kinds |

`flai touches suggest` listed nothing: the story declared no touches. No `after`: no open story writes what this one reads; S-0843 changes the callers of the guard, forecast replay, and share hold, not `Closed` or `PlanArchive`.

Forecast: `flai forecast` gave 27m, the median 147 s per unit over 70 done remediation stories in the medium band times size 11 (2 criteria, 9 touches). Raised to 45m: a change to `workitem`'s `Item` selects the tests of nearly every Go package that imports it in the close-out. Delivery moved by the same 18m, to 2026-11-17T21:10:00Z.

Cost of delay: `flai cod` gives 37.50 USD a week from the 15m time lost per 168h cycle at 150 USD an hour; it stands, since the input is flai's from I-0377's one instance.

### Accepted by the orchestrator

- Verified: fc5389f5d8fefe17fbf3083dd7e804d0ef4ef382
- At: 2026-11-18T08:31:06Z

Verdict: pass, both criteria are met at fc5389f and the close-out passed every step there.
- 1: flai/internal/workitem/item.go, flai/internal/workitem/repo.go, flai/internal/workitem/archive.go, flai/internal/workitem/view.go, flai/internal/workitem/archive_test.go, flai/internal/workitem/store_test.go, flai/cmd/accept_workflow_test.go, flai/cmd/archive_test.go, design/system/workflow.md
- 2: design/issues/I-0377-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md, design/issues/summary.md
