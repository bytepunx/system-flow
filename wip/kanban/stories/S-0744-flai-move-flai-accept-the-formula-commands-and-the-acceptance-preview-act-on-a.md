---
id: S-0744
type: story
nature: feature
title: flai move, flai accept, the formula commands, and the acceptance preview act on a transition's effect and who may make it instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:04Z
updated: 2026-11-12T16:00:04Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/cmd/move.go, flai/cmd/move_test.go, flai/cmd/move_cancel.go, flai/cmd/move_by_test.go, flai/cmd/accept.go, flai/cmd/accept_test.go, flai/cmd/formula_stop.go, flai/cmd/formula_stop_test.go, flai/cmd/formula.go, flai/cmd/plan.go, flai/cmd/board.go, flai/cmd/stream.go, flai/cmd/stream_sync.go, flai/cmd/stream_restore.go, flai/cmd/touches.go, flai/cmd/role.go, flai/cmd/analyze.go, flai/cmd/curate.go, flai/cmd/forge.go, flai/cmd/kluster.go, flai/cmd/serve_agent_run.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/cancel.go, flai/internal/workitem/transition.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

## Notes

- `cmd` holds about 43 references to a state's name in 18 files (`formula_stop.go` 10, `accept.go` 7, `move.go` 4), `preview` 17 (`accept.go` 12). Who may make each transition is a field the workflow gives from S-0737 on; today it is spread over `cmd/move.go`, `item_move`, and the guard.
