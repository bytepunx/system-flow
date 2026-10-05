---
id: S-0219
type: story
nature: feature
title: The orchestrator moves and orders work by its policy within its permissions
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T05:47:33Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:16Z
    by: alex
tags: [flai]
topics: [orchestration, planning]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", template, design/system/strategic-agents.md, flai/internal/guard, flai/internal/serve, flai/internal/workitem/board.go, flai/internal/workitem/board_placed_test.go, flai/internal/workitem/policy.go, flai/internal/workitem/policy_test.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/internal/workitem/finalizable.go, flai/internal/workitem/finalizable_test.go, flai/cmd/order.go, flai/cmd/order_test.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go, flai/cmd/promote.go, flai/cmd/promote_drafts_test.go, flai/cmd/move.go, flai/cmd/move_orchestrate_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/guard.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/testdata/orchestrate, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/server.go, design/adrs, design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators]
after: [S-0218, S-0209]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 101.35
  by: planner-S-0219
  at: 2026-10-05T04:48:56Z
forecast:
  duration: 1h40m
  delivery: 2026-10-05T13:53:00Z
  basis: "Its own forecast of 1h40m; 3rd in the pull order with an in-progress limit of 3, behind S-0276, S-0217 and S-0218."
  by: flai
  at: 2026-10-05T05:47:33Z
---
# S-0219 The orchestrator moves and orders work by its policy within its permissions

## Goal

With its permissions on, the orchestrator keeps the flow: it asks the planner to draft stories for backlog epics, finalizes drafts it judges ready, promotes backlog stories to ready while the limits leave room, and orders the ready column by the policy.

## Acceptance criteria
- [ ] With `plan_backlog_epics`, it runs the planner for a backlog epic that has no stories, and for one whose stories are all done when the epic is not
- [ ] With `finalize_drafts`, it finalizes a draft whose criteria, touches, forecast, and value are complete and consistent, and leaves one that is not with a thread saying why
- [ ] With `promote_to_ready`, it promotes `flai promote --candidates` in order while the ready limit has room, never a draft, never a held story
- [ ] With `order_ready`, it applies `flai order --by <policy> --apply` after each change to the ready column, and does not reorder a story the operator ordered by hand in the last day
- [ ] Every action is logged with the policy figure that justified it; a refusal from `flai guard` ends the attempt and is logged
- [ ] Tests cover each permission on and off against a fixture board

## Tasks
- T-0884 flai order records who placed a story by hand and when, and flai order --by --apply keeps a hand placement of the last day
- T-0888 flai plan --candidates lists the epics the planner should plan: backlog epics with no stories, and open epics whose stories are all done
- T-0890 flai promote --drafts lists each draft story with what it lacks to be finalized: criteria, touches, forecast, and value
- T-0893 flai guard passes the orchestrator's plan, finalize, promote, and order calls only while their permission is on, and its new reads always
- T-0896 The orchestrator can start the planner for an epic, finalize a complete draft, and promote only a candidate story to ready while the ready limit has room
- T-0899 The orchestrator's prompt and definition say how it plans, finalizes, promotes, and orders by its permissions, logging each action with its policy figure
- T-0901 A fixture-board test runs the orchestrator's calls for each of its four permissions, on and off, through the guard and flai
- T-0904 strategic-agents.md, flai-cli.md, workflow.md, and the user and operator guides describe how the orchestrator moves and orders work by its permissions

## Notes

### Planning

Planned by planner-S-0219 on 2026-10-05.

Touches:

- Declared, kept: `flai/internal/harness`, `.claude/agents/orchestrator.md`, `template`, `design/system/strategic-agents.md`, `flai/internal/guard`, `flai/internal/serve`.
- Layout, from the tasks: the files each task changes under `flai/internal/workitem` (the placement record and S-0217's `policy.go`; `plancandidates.go`; `finalizable.go`), `flai/cmd` (`order.go`, `plan.go`, `promote.go`, `move.go`, `edit.go`, `guard.go`, and their tests, with the fixture-board test and its `testdata`), and `flai/internal/mcpserver` (`plan.go`, `items_write.go`, `server.go`). Named file by file rather than by package, so the story holds less of `flai/cmd` and `flai/internal/workitem` against stories running beside it.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md` (30%), `docs/users/flai.md` (26%), `docs/users/flai-reference.md` (11%), `design/system/workflow.md` (10%), `design/adrs` (its README 12%). The rest it listed, such as `work-management.md`, `flaiover-dashboard.md`, `CLAUDE.md`, and `hostapi/writes.go`, nothing in the plan changes: the story adds no dashboard view and no hostapi write.
- Design: `design/conventions/strategic-agents.md` (the convention's "As the orchestrator", landed with the template) and `docs/operators` (the operator's guide S-0218 writes); `design/adrs` for where the hand-placement record lives (T-0884).

Forecast: 1h40m, delivery 2026-10-05T13:24Z, as `flai forecast` gave it once the touches were predicted, over 14 done feature stories on this model in the large band. Kept: it replaces the 1h30m written before the tasks existed, and sits between S-0210's 54m and S-0209's 98m of agent time, the closest stories in kind, for eight tasks in four layers. The file-level touches raise its size above what package-level ones would; the eight tasks bear the figure out. The delivery assumes S-0217 and S-0218, which it waits for, keep their forecasts.

Cost of delay: 101.35 USD a week, as `flai cod` gave it with the 1h40m written: this story's share of E-0016's 1500 USD a week (the operator's 10h lost per cycle), 1h40m of 24h40m over the epic's 17 open stories without inputs. Kept; the story has no inputs of its own, and the epic's are the operator's.
