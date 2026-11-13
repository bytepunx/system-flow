---
id: S-0744
type: story
nature: feature
title: flai move, flai accept, the formula commands, and the acceptance preview act on a transition's effect and who may make it instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:04Z
updated: 2026-11-13T11:42:32Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/cmd/move.go, flai/cmd/move_test.go, flai/cmd/move_cancel.go, flai/cmd/move_cancel_test.go, flai/cmd/move_by_test.go, flai/internal/workitem/transition.go, flai/cmd/accept.go, flai/cmd/accept_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/cancel.go, flai/cmd/formula_stop.go, flai/cmd/formula_stop_test.go, flai/cmd/formula.go, flai/cmd/plan.go, flai/cmd/board.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_restore.go, flai/cmd/touches.go, flai/cmd/role.go, flai/cmd/analyze.go, flai/cmd/curate.go, flai/cmd/forge.go, flai/cmd/kluster.go, flai/cmd/serve_agent_run.go]
after: [S-0737]
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
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 13
          cache_read: 108314
          cache_write: 18019
          cost: 0.0361
draft: true
forecast:
  duration: 1h
  delivery: 2026-11-13T20:11:00Z
  basis: "Its own forecast of 1h; 18th in the pull order with an in-progress limit of 3, behind S-0740, S-0763, S-0766, S-0760, S-0745, S-0759, S-0751, S-0757, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742 and S-0743."
  by: flai
  at: 2026-11-13T11:42:32Z
---
# S-0744 flai move, flai accept, the formula commands, and the acceptance preview act on a transition's effect and who may make it instead of a state's name

## Goal

The command line follows the workflow. `flai move` refuses a move its transition does not let the mover make, `flai accept` is the transition with the acceptance effect from whatever state the workflow says, a formula's stop closes and cancels along the workflow's transitions, and the acceptance preview lists the blockers the workflow's rules give, so the operator's own states move, accept, and stop with no change to the commands.

## Acceptance criteria

- [ ] The files of `flai/cmd` that move, accept, stop a formula, plan, show the board, and sync or restore a stream, and `flai/internal/preview`, name no state: each reads a state's kind or a transition's effect and who may make it from the workflow
- [ ] `flai move` refuses a transition the mover may not make, naming the transition, who may make it, and the command that does; the move to done stays the acceptance's, as today
- [ ] `flai accept` makes the transition with the acceptance effect, from the state it leaves, and refuses a story in any other state with today's words and that state's name
- [ ] `flai formula stop` closes the stopped step's task and cancels the steps after it along the workflow's transitions
- [ ] The acceptance preview names the blockers the acceptance transition's rules give, word for word as today with the baseline
- [ ] Tests run move, accept, formula stop, and the preview against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing `cmd` and `preview` test passes unchanged

## Tasks

- T-3221 flai move refuses a transition the mover may not make, naming it and who may, and cancels along the workflow's transitions
- T-3222 flai formula stop closes the stopped step and cancels the steps after it along the workflow's transitions
- T-3223 The plan, board, stream, touches, role, analyze, curate, forge, kluster, and agent-run commands read a state's kind
- T-3227 flai accept makes the transition with the acceptance effect, and the acceptance preview lists the blockers its rules give

## Notes

- `cmd` holds about 43 references to a state's name in 18 files (`formula_stop.go` 10, `accept.go` 7, `move.go` 4), `preview` 17 (`accept.go` 12). Who may make each transition is a field the workflow gives from S-0737 on; today it is spread over `cmd/move.go`, `item_move`, and the guard.

### Planning

- Touches, all files:
  - Layout: every file of `flai/cmd` with a reference to a state's name, but `stats.go` (S-0752's), and the three files of `flai/internal/preview` that name one, with the tests of the commands whose behaviour changes.
  - `flai/internal/workitem/transition.go`: layout. The helper `flai move` and `flai accept` share, over the workflow package's mover check.
  - Left out from `flai touches suggest`: the user guide and reference (19% to 28%), since no command's flags or help change; `stream_sync_test.go`, `branch.go`, `root.go`, and `storygit/sync.go` (S-0746's), at 8% to 13%.
- `after`: S-0737. `flai move` and `flai accept` read who may make each transition and the acceptance effect from the workflow S-0737 gives, through its mover check, and the tests use its test hook.
- Tasks: T-3221 (move), T-3222 (formula stop), and T-3223 (the other commands) in the first layer; T-3227 (accept and the preview) after T-3221, whose helper in `transition.go` it calls.
- Forecast 1h, as `flai forecast` gives it for size 33; 34 with `move_cancel_test.go` added rounds the same.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
