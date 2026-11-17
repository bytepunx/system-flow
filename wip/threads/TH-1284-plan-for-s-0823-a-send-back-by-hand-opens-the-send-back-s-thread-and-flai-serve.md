---
id: TH-1284
title: "Plan for S-0823: a send back by hand opens the send back's thread, and flai serve starts the agent on it"
anchor:
  path: wip/kanban/stories/S-0823-a-story-the-operator-moves-from-review-to-in-progress-by-hand-gets-no-agent-and.md
  item: S-0823
status: open
participants: [planner-S-0823]
created: 2026-11-17T08:31:22Z
updated: 2026-11-17T08:31:22Z
---

# TH-1284 Plan for S-0823: a send back by hand opens the send back's thread, and flai serve starts the agent on it

On wip/kanban/stories/S-0823-a-story-the-operator-moves-from-review-to-in-progress-by-hand-gets-no-agent-and.md.

## Entries

### 2026-11-17T08:31:22Z planner-S-0823
S-0823 is planned as three tasks in three layers. Its touches, forecast, and cost of delay are written, with the reasons under `### Planning` in its Notes.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3541 flai move along the send-back transition by hand opens the send back's thread with the move's reason | `flai/internal/sendback/sendback.go` and test, `flai/cmd/move.go`, new `flai/cmd/move_sendback_test.go`, `docs/users/flai-reference.md` |
| 2 | T-3543 flai serve starts the agent of a story sent back by hand and tells it who sent it back and why (after T-3541) | `flai/internal/harness/harness.go` and test, `flai/internal/serve/agents_restart.go`, `agents_start.go`, `agents_test.go` |
| 3 | T-3545 The workflow design and the user guides name the send back by hand, and I-0397 is closed (after T-3541, T-3543) | `design/system/workflow.md`, `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flaiover.md`, the I-0397 file, `design/issues/summary.md` |

Figures: forecast 36m, delivery 2026-11-17T15:01:00Z, from `flai forecast` and kept. Cost of delay 150 USD a week from `flai cod` on the 1h input from I-0397, kept. No `after` on another story.

Assumptions:

- Of the issue's two ways, the story takes the thread. The baseline's send-back transition already lists `reason-given`, and the dashboard's move dialog already asks for the reason, so every send back by hand has one and no board marker is needed.
- A send back by hand is a move along the `send-back` transition through `flai move`, which the dashboard's `item.move` runs, by anyone but `flai` and the story's own agent. The MCP `item_move` is left alone, since the operator does not move through it.
- The thread keeps the `Acceptance sent back:` prefix and `flai` as opener, so the existing `startAcceptSentBack` and `resolveAcceptSentBack` start the agent and resolve the thread. T-3543 only changes the words the agent and the journal get.
- No ADR is planned. The change extends ADR-0248's send-back thread to a move by hand, and `design/system/workflow.md` records it. The story's agent may add an ADR if it judges the decision new.
- The orchestrator's request named a thread "asked" that does not exist, and S-0823 had no earlier plan or thread. So this is a first plan, not a revision.
