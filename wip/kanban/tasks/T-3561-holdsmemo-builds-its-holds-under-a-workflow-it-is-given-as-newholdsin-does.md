---
id: T-3561
type: task
nature: remediation
title: HoldsMemo builds its Holds under a workflow it is given, as NewHoldsIn does
status: done
parent: S-0843
owner: alex
created: 2026-11-17T12:29:00Z
updated: 2026-11-18T10:21:11Z
transitions:
  - to: ready
    at: 2026-11-18T10:04:11Z
    by: agent-S-0843
  - to: in-progress
    at: 2026-11-18T10:04:11Z
    by: agent-S-0843
  - to: done
    at: 2026-11-18T10:21:11Z
    by: agent-S-0843
stream: S-0843
tags: [flai, workflow, claims]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go]
usage:
  source: log
  seconds: 1020
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2852
      cache_read: 832566
      cache_write: 37698
      cost: 0.449
---
# T-3561 HoldsMemo builds its Holds under a workflow it is given, as NewHoldsIn does

## Work

`HoldsMemo.NewHolds` in `flai/internal/workitem/hold.go` always calls `newHolds(baselineFlow, ...)`, so a memo's Holds read a task's closed state, and a story's active state, by the baseline's names. Under a project's own workflow file that renames done or cancelled, a closed task's touches stay in its story's claim (I-0379).

Add `HoldsMemo.NewHoldsIn(w workflow.Workflow, items, projects)`, which calls `newHolds(w, items, projects, m)`, and keep `NewHolds` as `NewHoldsIn` under the baseline, as the package's other wrappers are. A memo's Holds must be given the same workflow each time, as they are given the same projects: say so in the memo's comment.

Waits for nothing: it is the first layer. The metrics task waits for it, since `flai/internal/metrics/claims.go` calls the new method.

## Done when

- `HoldsMemo.NewHoldsIn` exists and `HoldsMemo.NewHolds` is its baseline wrapper.
- A test in `flai/internal/workitem/hold_test.go` builds a workflow whose closed task state is not named `done`, gives a story a task in that state, and shows the claim from `memo.NewHoldsIn` leaves the task's touches out while `memo.NewHolds` keeps them.
- `flai test flai/internal/workitem/hold.go flai/internal/workitem/hold_test.go` passes.

## Notes
