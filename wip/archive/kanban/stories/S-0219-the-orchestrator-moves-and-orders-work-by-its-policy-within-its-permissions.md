---
id: S-0219
type: story
nature: feature
title: The orchestrator moves and orders work by its policy within its permissions
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-06T03:45:18Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:16Z
    by: alex
  - to: in-progress
    at: 2026-10-06T03:19:18Z
    by: agent-S-0219
  - to: review
    at: 2026-10-06T03:37:22Z
    by: agent-S-0219
  - to: done
    at: 2026-10-06T03:45:18Z
    by: alex
tags: [flai]
topics: [orchestration, planning]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", template, design/system/strategic-agents.md, flai/internal/guard, flai/internal/serve, flai/internal/workitem/board.go, flai/internal/workitem/board_placed_test.go, flai/internal/workitem/policy.go, flai/internal/workitem/policy_test.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/internal/workitem/finalizable.go, flai/internal/workitem/finalizable_test.go, flai/cmd/order.go, flai/cmd/order_test.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go, flai/cmd/promote.go, flai/cmd/promote_drafts_test.go, flai/cmd/move.go, flai/cmd/move_orchestrate_test.go, flai/cmd/edit.go, flai/cmd/edit_test.go, flai/cmd/guard.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/testdata/orchestrate, flai/internal/mcpserver/plan.go, flai/internal/mcpserver/plan_test.go, flai/internal/mcpserver/items_write.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/server.go, design/adrs, design/conventions/strategic-agents.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators, flai/internal/workitem/promotable.go, flai/internal/mcpserver/folder.go, flai/cmd/plan_test.go, flai/cmd/plan_orchestrate_test.go]
after: [S-0218, S-0209]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2409
  models:
    - model: claude-opus-5-5
      input: 918
      output: 341328
      cache_read: 51782921
      cache_write: 1194240
      cost: 23.9634
    - model: claude-sonnet-5-5
      input: 20
      output: 4791
      cache_read: 181306
      cache_write: 60933
      cost: 0.2365
  strategic:
    - kind: planner
      seconds: 73
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 7567
          cache_read: 545410
          cache_write: 68587
          cost: 0.8092
cost_of_delay:
  value: 141.38
  by: planner-S-0219
  at: 2026-10-06T03:15:49Z
forecast:
  duration: 1h40m
  delivery: 2026-10-06T04:37:00Z
  basis: "Kept at 1h40m against flai's 2h3m, which file-level touches inflate; the stream opened 02:57Z and five of eight tasks were done by 03:15Z, leaving the prompt, the fixture test, the docs, and close-out."
  by: planner-S-0219
  at: 2026-10-06T03:15:49Z
---
# S-0219 The orchestrator moves and orders work by its policy within its permissions

## Goal

With its permissions on, the orchestrator keeps the flow: it asks the planner to draft stories for backlog epics, finalizes drafts it judges ready, promotes backlog stories to ready while the limits leave room, and orders the ready column by the policy.

## Acceptance criteria
- [x] With `plan_backlog_epics`, it runs the planner for a backlog epic that has no stories, and for one whose stories are all done when the epic is not
- [x] With `finalize_drafts`, it finalizes a draft whose criteria, touches, forecast, and value are complete and consistent, and leaves one that is not with a thread saying why
- [x] With `promote_to_ready`, it promotes `flai promote --candidates` in order while the ready limit has room, never a draft, never a held story
- [x] With `order_ready`, it applies `flai order --by <policy> --apply` after each change to the ready column, and does not reorder a story the operator ordered by hand in the last day
- [x] Every action is logged with the policy figure that justified it; a refusal from `flai guard` ends the attempt and is logged
- [x] Tests cover each permission on and off against a fixture board

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

Planned by planner-S-0219 on 2026-10-05; revisited on 2026-10-06 while the story's agent worked it, with T-0884, T-0888, T-0890, T-0893, and T-0896 done and T-0899 in progress.

Touches (43, all kept):

- Declared, kept: `flai/internal/harness`, `.claude/agents/orchestrator.md`, `template`, `design/system/strategic-agents.md`, `flai/internal/guard`, `flai/internal/serve`.
- Layout, from the tasks: the files each task changes under `flai/internal/workitem` (the placement record and S-0217's `policy.go`; `plancandidates.go`; `finalizable.go`; `promotable.go`), `flai/cmd` (`order.go`, `plan.go`, `promote.go`, `move.go`, `edit.go`, `guard.go`, and their tests, with the fixture-board test and its `testdata`), and `flai/internal/mcpserver` (`plan.go`, `items_write.go`, `server.go`, `folder.go`). Named file by file rather than by package, so the story holds less of `flai/cmd` and `flai/internal/workitem` against stories running beside it. The story's agent added `promotable.go`, `folder.go`, `plan_test.go`, and `plan_orchestrate_test.go` while working T-0896; every path a task names falls under the story's touches.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/workflow.md`, and `design/adrs` (ADR-0088, from T-0884). On the revisit it lists `flaiover-dashboard.md` (18%), `docs/users/flaiover.md` (14%), `design/issues/summary.md` (11%), `work-management.md`, `work-hierarchy.md`, and `hostapi/writes.go` (7% each), and others below that. None is added: the story adds no dashboard view, no hostapi write, and no item field. The dashboard's drag still records `flaiover` as the placer, which the narrative leaves outside this story.
- Design: `design/conventions/strategic-agents.md` (the convention's "As the orchestrator", landed with the template) and `docs/operators` (the operator's guide S-0218 wrote, with the flag index in `settings.md` that T-0904 brings up to date).

Forecast: 1h40m, kept, delivery 2026-10-06T04:37Z. On the revisit `flai forecast` gave 2h3m, delivery 06:22Z: 151 s a unit of size over 16 done feature stories on this model in the large band, times size 49 (6 criteria, 43 touches). The rise from 1h40m comes from the touches named file by file, which count more than the package-level touches of the stories it learns from, so the figure is not kept. The work so far bears out 1h40m: the stream opened at 02:57Z and five of eight tasks were done by 03:15Z. The rest is a prompt and template release, one test, the documents, and close-out. Delivery is the stream's start plus 1h40m, not flai's pull-order date, since the story is already being worked.

Cost of delay: 141.38 USD a week, as `flai cod` gives it: this story's share of E-0016's 1500 USD a week (the operator's 10h lost per cycle), 1h40m of 17h41m forecast over the epic's 15 open stories without inputs. It replaces the 101.35 of 2026-10-05 only because the epic's other forecasts moved. The story has no inputs of its own, and the epic's are the operator's.
