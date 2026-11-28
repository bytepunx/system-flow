---
id: S-0844
type: story
nature: remediation
title: "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds"
status: done
owner: alex
created: 2026-11-17T08:32:55Z
updated: 2026-11-28T00:43:34Z
transitions:
  - to: ready
    at: 2026-11-18T09:20:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-27T23:37:13Z
    by: agent-S-0844
  - to: review
    at: 2026-11-28T00:41:59Z
    by: agent-S-0844
  - to: done
    at: 2026-11-28T00:43:34Z
    by: orchestrator
tags: [flai]
topics: [workflow, acceptance, orchestration, dashboard]
touches: [flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/cmd/accept_workflow_file_test.go, flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go, flai/internal/hostapi/writes_test.go, flai/internal/check/check_items.go, flai/internal/check/items_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/internal/workitem/board.go, flai/internal/workitem/board_update_test.go, flai/internal/workitem/boardview_test.go, flai/cmd/board.go, flai/cmd/board_limit_test.go, docs/users/flai-reference.md, docs/operators/runbooks/change-the-board.md, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/summary.md, flai/cmd/accept.go, flai/internal/preview/accept_test.go, flai/cmd/check.go, flai/internal/check/itemwrite_test.go, design/issues/I-0420-the-dashboard-s-review-page-shows-operator-only-only-beside-protected-paths-so.md, design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3903
  priced_by: harness
  turns:
    - day: 2026-11-27
      ceremony: 1
      work: 27
    - day: 2026-11-28
      ceremony: 5
      test_runs: 1
      work: 66
  models:
    - model: claude-opus-5-5
      input: 618
      output: 154568
      cache_read: 46838785
      cache_write: 1288808
      cost: 20.108
  strategic:
    - kind: planner
      seconds: 89
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 38
          output: 7959
          cache_read: 3568357
          cache_write: 91660
          cost: 1.4863
    - kind: orchestrator
      seconds: 79
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 208
          cache_read: 974888
          cache_write: 108757
          cost: 0.3399
        - model: claude-sonnet-5-5
          input: 1
          output: 8
          cache_read: 7641
          cache_write: 5121
          cost: 0.0103
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-17T08:32:55Z
  value: 37.5
  by: planner-S-0844
  at: 2026-11-17T12:37:42Z
forecast:
  duration: 1h30m
  delivery: 2026-11-28T01:28:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 3, behind S-0846, S-0847, S-0850 and S-0808."
  by: flai
  at: 2026-11-27T22:47:30Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:41:08Z
---
# S-0844 ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds

## Goal

This story remediates [I-0380](../../../design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md), "ADR-0301 rules for a project's own workflow file that no code enforces yet: operator-only acceptance of a change to it, the dashboard's refusal to save it on main, board.wip-state, and WIP limits on a state the project adds". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0380 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0380 is closed with `flai issue close I-0380 --reason` saying what fixed it

## Tasks
- T-3595 Only the operator accepts a story whose branch changes the workflow file: the acceptance preview, the orchestrator's blockers, and the merge queue's eligibility
- T-3596 flai doc save, and so the dashboard's doc.save, refuses the workflow file on main and takes it only on a story's branch
- T-3597 flai check reports board.wip-state for a limit in wip/kanban/board.md on a state that takes none
- T-3598 A state the project adds with limit true takes a WIP limit through flai board limit, whose help names no baseline states, and the baseline-only LimitedColumns and DefaultWIPLimits are retired
- T-3599 The guides and the design say the four rules are enforced, and I-0380 is closed naming what fixed each

## Notes

Cost of delay inputs set by flai from I-0380. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-15T23:00:06Z, 1.4 days before this story; under one cycle counts as one).

### Planning

The proposed solution, one task per instance of I-0380, then the docs:

1. Operator-only acceptance (T-3595). Use a predicate of its own in `flai/internal/preview`, beside the ADR-0106 check. Do not add the file to `flai/internal/protected`, which also drives `permission_prompt` and the guard. Add a new orchestrator blocker and a new merge queue eligibility kind, so the queue never batches such a story.
2. Save refusal (T-3596). `docedit.Save` refuses the workflow file when no story is named. `doc.save` and `flai doc save` share that refusal. `file.save` already refuses every `.md` file.
3. `board.wip-state` (T-3597). A warning in `(*checker).board`, as `board.wip-limit` is, for a `wip_limits` key that names no state with `limit: true`.
4. WIP limits on an added state (T-3598). `SetLimit` and `Save` already follow the repository's workflow, since `LoadBoard` sets it. What is left is the baseline-only `LimitedColumns` and `DefaultWIPLimits`, which no production code reads, and the help of `flai board limit`, which is also I-0403's second half. A test reproduces the instance.

Touches, by where each came from:

- Design, from ADR-0301 §6 and §7 and I-0380's instances, located in the code by the explorer: every `flai/` path above, which are T-3595 to T-3598's.
- Layout: `docs/users/flai-reference.md`, generated from the help T-3598 changes.
- Design: `docs/operators/runbooks/change-the-board.md`, which states rules 1 and 2 and limitation 4, and is named by I-0380. `design/system/workflow.md` § The workflow file. `design/system/flai-cli.md` and `docs/users/flai.md`, for the blocker and the check rule.
- Co-change (`flai touches suggest`, 20%), and the `doc.save` route row: `design/system/flaiover-dashboard.md`. `docs/users/flaiover.md` covers the editor's saves.
- Layout: `design/issues/I-0380-…md` and `design/issues/summary.md`, written by `flai issue close` (criterion 2).
- None declared before planning. No folder touches.

The story names no `after`. It reads no output of an open story: S-0755, which made flai read a project's own workflow file, is accepted.

The figures:

- Forecast 1h30m, against flai's 50m. flai's figure is 119 s per unit over 82 done large-band remediation stories, about 55m at the final size of 28. It was raised because there are five tasks across `preview`, `mergequeue`, `docedit`, `check`, and `workitem`, each with its own test. The close-out also runs integration and smoke for a change under `flai/`: S-0828's close-out took 9m43s today. Delivery is flai's 21:47Z plus the 40m added.
- Cost of delay value 37.50 USD a week, as `flai cod` gives it from flai's input, 15m lost per 168h cycle at 150 USD an hour. It stands: one occurrence, and S-0756, the epic's end-to-end story, will meet instances 1 and 4.

### Accepted by the orchestrator

- Verified: 0beb074eac7ab5c6feee3d9ff326afe2aa5432d1
- At: 2026-11-28T00:43:34Z

Verdict: pass at 0beb074eac7ab5c6feee3d9ff326afe2aa5432d1
- 1: flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/mergequeue/eligible.go, flai/cmd/accept.go, flai/internal/docedit/docedit.go, flai/internal/check/check_items.go, flai/internal/check/scope.go, flai/cmd/check.go, flai/internal/workitem/board.go, flai/cmd/board.go, flai/cmd/accept_workflow_file_test.go, flai/internal/preview/accept_test.go, flai/internal/mergequeue/eligible_test.go, flai/internal/docedit/docedit_test.go, flai/internal/hostapi/writes_test.go, flai/internal/check/items_test.go, flai/internal/check/scope_test.go, flai/internal/check/itemwrite_test.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/board_update_test.go, flai/cmd/board_limit_test.go
- 2: design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md, design/issues/summary.md
