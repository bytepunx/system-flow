---
id: S-0744
type: story
nature: feature
title: flai move, flai accept, the formula commands, and the acceptance preview act on a transition's effect and who may make it instead of a state's name
status: done
parent: E-0034
owner: alex
created: 2026-11-12T16:00:04Z
updated: 2026-11-15T19:43:52Z
transitions:
  - to: ready
    at: 2026-11-14T17:47:59Z
    by: alex
  - to: in-progress
    at: 2026-11-15T00:56:52Z
    by: agent-S-0744
  - to: review
    at: 2026-11-15T01:47:19Z
    by: agent-S-0744
  - to: in-progress
    at: 2026-11-15T19:29:34Z
    by: flai
  - to: review
    at: 2026-11-15T19:42:20Z
    by: agent-S-0744
  - to: done
    at: 2026-11-15T19:43:52Z
    by: orchestrator
tags: [cli]
topics: [cli]
touches: [flai/cmd/move.go, flai/cmd/move_cancel.go, flai/internal/workitem/transition.go, flai/cmd/accept.go, flai/internal/preview/accept.go, flai/internal/preview/orchestrator.go, flai/internal/preview/cancel.go, flai/cmd/formula_stop.go, flai/cmd/formula.go, flai/cmd/plan.go, flai/cmd/board.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_restore.go, flai/cmd/touches.go, flai/cmd/forge.go, flai/cmd/kluster.go, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/move_workflow_test.go, flai/cmd/project.go, flai/cmd/formula_stop_workflow_test.go, flai/cmd/commands_workflow_test.go, design/system/workflow.md, flai/cmd/accept_workflow_test.go, flai/internal/preview/accept_workflow_test.go, design/issues/I-0377-workitem-s-item-closed-and-planarchive-read-the-baseline-s-done-and-cancelled.md, design/issues/summary.md]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3791
  priced_by: harness
  turns:
    - day: 2026-11-15
      ceremony: 4
      test_runs: 3
      hand_edits: 1
      work: 97
  models:
    - model: claude-opus-5-5
      input: 444
      output: 137554
      cache_read: 42819012
      cache_write: 1257557
      cost: 19.9801
  strategic:
    - kind: orchestrator
      seconds: 195
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 62
          output: 550
          cache_read: 2167850
          cache_write: 300564
          cost: 0.7392
        - model: claude-sonnet-5-5
          input: 16
          output: 957
          cache_read: 169242
          cache_write: 74114
          cost: 0.1919
cost_of_delay:
  value: 106.89
  by: system-flow
  at: 2026-11-14T17:44:06Z
forecast:
  duration: 1h
  delivery: 2026-11-15T01:58:00Z
  basis: "Its own forecast of 1h; 6th in the pull order with an in-progress limit of 3, behind S-0738, S-0741, S-0746, S-0739 and S-0743."
  by: flai
  at: 2026-11-14T23:48:54Z
finalized:
  by: orchestrator
  at: 2026-11-14T17:45:54Z
---
# S-0744 flai move, flai accept, the formula commands, and the acceptance preview act on a transition's effect and who may make it instead of a state's name

## Goal

The command line follows the workflow. `flai move` refuses a move its transition does not let the mover make, `flai accept` is the transition with the acceptance effect from whatever state the workflow says, a formula's stop closes and cancels along the workflow's transitions, and the acceptance preview lists the blockers the workflow's rules give, so the operator's own states move, accept, and stop with no change to the commands.

## Acceptance criteria

- [x] The files of `flai/cmd` that move, accept, stop a formula, plan, show the board, and sync or restore a stream, and `flai/internal/preview`, name no state: each reads a state's kind or a transition's effect and who may make it from the workflow
- [x] `flai move` refuses a transition the mover may not make, naming the transition, who may make it, and the command that does; the move to done stays the acceptance's, as today
- [x] `flai accept` makes the transition with the acceptance effect, from the state it leaves, and refuses a story in any other state with today's words and that state's name
- [x] `flai formula stop` closes the stopped step's task and cancels the steps after it along the workflow's transitions
- [x] The acceptance preview names the blockers the acceptance transition's rules give, word for word as today with the baseline
- [x] Tests run move, accept, formula stop, and the preview against a workflow with a state renamed and one added between `in-progress` and `review`
- [x] With the baseline, every existing `cmd` and `preview` test passes unchanged

## Tasks

- T-3221 flai move refuses a transition the mover may not make, naming it and who may, and cancels along the workflow's transitions
- T-3222 flai formula stop closes the stopped step and cancels the steps after it along the workflow's transitions
- T-3223 The plan, board, stream, touches, role, analyze, curate, forge, kluster, and agent-run commands read a state's kind
- T-3227 flai accept makes the transition with the acceptance effect, and the acceptance preview lists the blockers its rules give

## Notes

- `cmd` holds about 43 references to a state's name in 18 files (`formula_stop.go` 10, `accept.go` 7, `move.go` 4), `preview` 17 (`accept.go` 12). Who may make each transition is a field the workflow gives from S-0737 on; today it is spread over `cmd/move.go`, `item_move`, and the guard.
- Criterion 3 is verified as written: on a board whose accepted state is `shipped`, `flai accept` walks the story to `shipped` from the state the acceptance leaves and refuses a story in `testing` with today's words. The acceptance's next step, the archive, still stops there, since `workitem`'s `Item.Closed` and `PlanArchive` read the names `done` and `cancelled`, outside this story's files. I-0377 records it. On a board that keeps the name `done`, the whole acceptance runs.
- Prose in the commands' help and in a few messages still names the baseline's states, such as "stays in review" and "done and archived". Criterion 1 is about the code's reads, and every comparison with a state's name is gone from the story's files.

### Planning

- Touches, all files:
  - Layout: every file of `flai/cmd` with a reference to a state's name, but `stats.go` (S-0752's), and the three files of `flai/internal/preview` that name one, with the tests of the commands whose behaviour changes.
  - `flai/internal/workitem/transition.go`: layout. The helper `flai move` and `flai accept` share, over the workflow package's mover check.
  - Left out from `flai touches suggest`: the user guide and reference (19% to 28%), since no command's flags or help change; `stream_sync_test.go`, `branch.go`, `root.go`, and `storygit/sync.go` (S-0746's), at 8% to 13%.
- `after`: S-0737. `flai move` and `flai accept` read who may make each transition and the acceptance effect from the workflow S-0737 gives, through its mover check, and the tests use its test hook.
- Tasks: T-3221 (move), T-3222 (formula stop), and T-3223 (the other commands) in the first layer; T-3227 (accept and the preview) after T-3221, whose helper in `transition.go` it calls.
- Forecast 1h, as `flai forecast` gives it for size 33; 34 with `move_cancel_test.go` added rounds the same.
- Cost of delay: 106.89 USD a week, as `flai cod` gives it: this story's share of E-0034's 1500.00 USD a week, the operator's `time_lost_per_cycle: 10h` (TH-1115) at 150 USD an hour, by its forecast of 1h over the 14h2m of the epic's 19 open stories. Kept: the epic's outcome needs every story, so a share by size is the split that fits.
- 2026-11-15T19:29:34Z: moved to in-progress: alex asked its agent to rebase story/S-0744 onto main and resolve what holds its acceptance, 1 blocker: S-0744's verify does not cover the main branch it would land on: its last flai verify, at 18bbb55, ran over main at 449ee65, and 44 of the 67 paths main brought...; its agent syncs it with flai stream sync S-0744, takes each blocker's recommended action, closes out again, and moves it to review

### Accepted by the orchestrator

- Verified: b21e28ecde7feaa061a327317a10443237347442
- At: 2026-11-15T19:43:52Z

Verdict: meets all 7 criteria at b21e28ecde7feaa061a327317a10443237347442; close-out passed every step over main e20b6d42f.
- 1: flai/cmd/move.go, flai/cmd/move_cancel.go, flai/cmd/accept.go, flai/cmd/formula_stop.go, flai/cmd/formula.go, flai/cmd/plan.go, flai/cmd/board.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_restore.go, flai/cmd/touches.go, flai/cmd/forge.go, flai/cmd/kluster.go, flai/internal/preview/accept.go, flai/internal/preview/cancel.go, flai/internal/preview/orchestrator.go, flai/internal/workitem/transition.go
- 2: flai/cmd/move.go, flai/internal/workitem/transition.go, flai/cmd/move_workflow_test.go
- 3: flai/cmd/accept.go, flai/internal/workitem/transition.go, flai/cmd/accept_workflow_test.go
- 4: flai/cmd/formula_stop.go, flai/cmd/formula.go, flai/cmd/formula_stop_workflow_test.go
- 5: flai/internal/preview/accept.go, flai/internal/preview/accept_workflow_test.go
- 6: flai/cmd/move_workflow_test.go, flai/cmd/accept_workflow_test.go, flai/cmd/formula_stop_workflow_test.go, flai/internal/preview/accept_workflow_test.go, flai/cmd/commands_workflow_test.go
- 7: flai/cmd/move_workflow_test.go, flai/cmd/accept_workflow_test.go, flai/cmd/formula_stop_workflow_test.go, flai/cmd/commands_workflow_test.go, flai/internal/preview/accept_workflow_test.go
