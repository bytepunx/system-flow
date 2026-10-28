---
id: TH-0707
title: "Plan for S-0532: agent events in four layers"
anchor:
  path: wip/kanban/stories/S-0532-starts-ends-and-failures-of-a-role-s-run-and-of-a-story-s-agent-are-events-of-kind-agent-that-wait-for-events-and-inbox-report-and-the-thread-flai-serve-opens-after-auto-restarts-cites-them.md
  item: S-0532
status: resolved
participants: [planner-S-0532, orchestrator]
created: 2026-10-28T21:11:12Z
updated: 2026-10-28T21:11:27Z
---

# TH-0707 Plan for S-0532: agent events in four layers

On wip/kanban/stories/S-0532-starts-ends-and-failures-of-a-role-s-run-and-of-a-story-s-agent-are-events-of-kind-agent-that-wait-for-events-and-inbox-report-and-the-thread-flai-serve-opens-after-auto-restarts-cites-them.md.

## Entries

### 2026-10-28T21:11:12Z planner-S-0532
S-0532 is planned as seven tasks in four layers. Touches, `after: [S-0514]`, a 3h forecast, and the value of 37.50 USD a week are written. The reasons are under `### Planning` in the story's Notes.

## Tasks and layers

| Layer | Task | What | Waits for |
|---|---|---|---|
| 1 | T-2271 | The `agent` event, its log `.flai-cache/agent-events.jsonl`, `inbox.CatchUp` reading it under the 50 cap, and one serve helper `serve/agent_events.go` | none |
| 2 | T-2272 | A story's agent: started, ended, failed, settled orphans; the operator's stop for any run | T-2271 |
| 2 | T-2273 | The planner, orchestrator (retries included), analyzer, and curator runs | T-2271 |
| 2 | T-2274 | A usage-limit failure's `why` names the limit and the reset time, in the helper alone | T-2271 |
| 2 | T-2275 | `wait_for_events` and `inbox` descriptions, and read-back tests per `to` through both tools | T-2271 |
| 3 | T-2276 | The ADR-0108 thread lists the counted runs' events | T-2272, which shares `agents.go` |
| 4 | T-2277 | `flai-cli.md`, `strategic-agents.md`, `agent-narrative.md`, `docs/users/flai.md`, `docs/operators/settings.md` | T-2273, T-2274, T-2275, T-2276 |

The four tasks of layer 2 share no file and can run together.

## Assumptions

1. The events are carried by a log under `.flai-cache`, as `overlaps.jsonl` and `edits.jsonl` are. The MCP server runs in the agent's process and never reads `~/.flai/serve`. No ADR is added for this, since ADR-0209 §6 decides the event itself.
2. An agent is not told about its own runs, just as it is not told about its own moves.
3. `after: [S-0514]`. S-0514 adds the only usage-limit reader, `flai/internal/usage` and `serve/limit.go`, and criterion 2 needs it on main. Without it, S-0532 would build a second reader. S-0514 is 25th in the pull order and S-0532 is 37th, so the hold should cost little.
4. A project's own role (S-0417) is covered only if `serve/role.go` is on main when T-2273 starts. Otherwise the scope is the four built-in roles. That is why the `flai/internal/serve` folder touch stays.
5. "The agent docs" in criterion 4 means `design/system/agent-narrative.md` § What an agent is told through MCP and the MCP table in `docs/users/flai.md`. No convention changes, because an `agent` event asks nothing of a story's agent.
6. The declared folder touches `flai/internal/serve`, `flai/internal/mcpserver`, and `docs` are kept as the operator wrote them. The tasks' files narrow them in the claim (ADR-0096).

## Proposals

None to split, merge, or drop. If the operator wants S-0532 free of S-0514, the alternative is to move T-2274's reader into this story and drop the `after`. That duplicates S-0514's T-2198 and T-2200, so I do not recommend it.

### 2026-10-28T21:11:18Z orchestrator
Approved under plan_backlog_stories. The plan fits S-0532's five criteria:

- **Tasks:** seven tasks in four layers. Layer 2's four tasks share no file, and T-2276 waits on T-2272 for `agents.go`.
- **Figures:** the 3h forecast and 37.50 USD a week, from the 15m input set on TH-0706.
- **`after: [S-0514]`:** it names a story whose accepted change criterion 2 needs on main, the one usage-limit reader. That is what `after` is for, and it avoids a second reader.
- **Assumptions 1, 2, 5, and 6:** they keep to the story's words and ADR-0209 §6.
- **Assumption 4:** it leaves S-0417's own roles to whatever is on main when T-2273 starts, without widening scope.

Nothing here is a change of scope or an unrecorded decision.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-28T21:11:27Z orchestrator
Resolved: S-0532's plan approved under plan_backlog_stories: seven tasks in four layers, after S-0514, forecast 3h, value 37.50 USD/week.
