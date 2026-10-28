---
id: S-0532
type: story
nature: improvement
title: Starts, ends, and failures of a role's run and of a story's agent are events of kind agent that wait_for_events and inbox report, and the thread flai serve opens after auto_restarts cites them
status: backlog
parent: E-0030
owner: alex
created: 2026-10-28T21:00:58Z
updated: 2026-10-28T21:33:08Z
transitions: []
tags: [flai, serve, events]
topics: [orchestration, cli]
touches: [flai/internal/serve, flai/internal/mcpserver, design/system/flai-cli.md, design/system/strategic-agents.md, docs, flai/internal/inbox/agentevents.go, flai/internal/inbox/agentevents_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/serve/agent_events.go, flai/internal/serve/agent_events_test.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/stop.go, flai/internal/serve/stop_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/limit.go, flai/internal/serve/limit_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, design/system/agent-narrative.md, docs/users/flai.md, docs/operators/settings.md]
after: [S-0514]
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
      seconds: 422
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 55
          output: 835
          cache_read: 4547428
          cache_write: 20603
          cost: 1.0052
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-10-28T21:03:12Z
  value: 37.5
  by: planner-S-0532
  at: 2026-10-28T21:10:53Z
forecast:
  duration: 3h
  delivery: 2026-10-29T15:13:00Z
  basis: "Its own forecast of 3h; 37th in the pull order with an in-progress limit of 5, behind S-0416, S-0419, S-0423, S-0418, S-0421, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517, S-0518, S-0519, S-0520, S-0523, S-0524, S-0525, S-0527 and S-0529."
  by: flai
  at: 2026-10-28T21:33:08Z
---
# S-0532 Starts, ends, and failures of a role's run and of a story's agent are events of kind agent that wait_for_events and inbox report, and the thread flai serve opens after auto_restarts cites them

## Goal

ADR-0209 §6. Today an agent that restarts three times is visible only in the serve journal and in the thread flai serve opens after `agent.auto_restarts` (ADR-0108). The orchestrator, the analyzer, and the operator reading the inbox see nothing of it, and nothing tells a running agent why a run around it ended. This story makes each start, end, and failure of a role's run (planner, orchestrator, analyzer, curator, a project's own role) and of a story's agent an event of kind `agent`, with `to` one of `started`, `ended`, `failed`, `stopped`, the item or role it concerns, and `why` as the journal words it (the exit, the usage-limit refusal, the operator's stop, the auto-restart count), reported by `wait_for_events` and `inbox` as the other events are, under the 50-event cap, and cited by the ADR-0108 thread.

## Acceptance criteria
- [ ] A role's run and a story's agent starting, ending, failing, or being stopped produces one event of kind `agent` with its item or role, `to`, and `why`, read back through `wait_for_events` and `inbox`, with a test for each kind
- [ ] A usage-limit refusal is a `failed` event whose `why` names the limit and the reset time Claude Code reports
- [ ] The thread ADR-0108 opens after `auto_restarts` lists the `agent` events of those runs
- [ ] The MCP tool descriptions, the design (`flai-cli.md`, `strategic-agents.md`), and the agent docs name the event kind and its fields
- [ ] `flai check --strict` is clean

## Tasks
- T-2271 An agent event of kind agent is appended to a log under .flai-cache and read back by inbox.CatchUp, once, under the 50-event cap
- T-2272 flai serve appends a started, ended, failed, or stopped agent event for each story's agent run and for any run the operator stops
- T-2273 The planner's, the orchestrator's, the analyzer's, and the curator's runs append a started, ended, or failed agent event naming their role and item
- T-2274 A run that ended on the harness's usage limit is a failed agent event whose why names the limit and the reset time Claude Code reported
- T-2275 The wait_for_events and inbox tool descriptions name the agent kind and its fields, and a test per to reads an agent event back through each tool
- T-2276 The thread flai serve opens after agent.auto_restarts lists the agent events of the runs it counted
- T-2277 The design, the user docs, and the agent docs name the agent event kind, its fields, its log, and what the auto_restarts thread lists

## Notes

Evidence: ADR-0209's context. On 2026-10-26 and 27, five projects' orchestrators were started and refused 8,986 times, and only the serve journal recorded it.

### Planning

Planned by planner-S-0532 on 2026-10-28. The tasks, their layers, and the assumptions are on the plan's thread.

The design: `wait_for_events` and `inbox` both read `inbox.CatchUp`, which merges work item changes with two logs under `.flai-cache` (`edits.jsonl`, `overlaps.jsonl`). flai serve's run state lives in `~/.flai/serve/agents.json`, which the MCP server does not read. So flai serve appends each run's event to a third log, `.flai-cache/agent-events.jsonl`, through one helper, and `CatchUp` reads it. No usage-limit detection exists in flai today; S-0514 adds the reader (`flai/internal/usage`, `flai/internal/serve/limit.go`), so this story names S-0514 in `after` to use it rather than build a second one.

Touches, and where each came from:

- `flai/internal/serve`, `flai/internal/mcpserver`, `docs`: declared, and kept. The tasks name the files below each folder, so the story's claim narrows each folder to those files (ADR-0096). `flai/internal/serve` also stays because T-2273 adds `flai/internal/serve/role.go` if S-0417 has put it on main by then, a file no task can name yet. `flai/internal/mcpserver` and `docs` stay because the operator declared them; every file below them the plan expects is named.
- `design/system/flai-cli.md`, `design/system/strategic-agents.md`: declared; criterion 4.
- `flai/internal/inbox/inbox.go`, `inbox_test.go`: layout; `CatchUp`, the kinds, `Describe`, and the 50 cap.
- `flai/internal/inbox/agentevents.go`, `agentevents_test.go`, `flai/internal/serve/agent_events.go`, `agent_events_test.go`: layout, new files for the log and the one helper every run path calls.
- `flai/internal/serve/agents.go` (co-change), `agents_test.go`, `stop.go`, `stop_test.go`: layout; a story's agent's start, end, orphan settling, the operator's stop, and the ADR-0108 thread in `stopRestarting`.
- `flai/internal/serve/plan.go`, `orchestrate.go`, `analyze.go`, `curate.go` and their tests: layout; each role's start and end.
- `flai/internal/serve/limit.go`, `limit_test.go`: design; S-0514's reader of a run's usage limit, which T-2274 uses.
- `flai/internal/mcpserver/server.go`, `server_test.go`, `folder.go`, `folder_test.go`: layout; the `wait_for_events` and `inbox` descriptions and their read-back tests.
- `design/system/agent-narrative.md`: design, § What an agent is told through MCP, the agent docs criterion 4 names.
- `docs/users/flai.md`, `docs/operators/settings.md`: design; the MCP tool table and the `agent.auto_restarts` entries.
- Left out of the suggestion: `flaiover-dashboard.md`, `design/adrs/README.md`, `workflow.md`, `template/CHANGELOG.md`, `design/issues/summary.md`, the host API, `serve_actions.go`, `root.go`, `check.go`, the harness and the manifest: co-change of past serve stories, but nothing in them changes here. No ADR is added: ADR-0209 §6 decides the event, and the log file is an implementation detail recorded in `flai-cli.md`.

Forecast: `flai forecast` gave 15m (median 87 s per unit of size over 20 done improvement stories on claude-opus-5-5, size 10: 5 criteria, 5 touches). Raised to 3h: seven tasks across the inbox, five serve run paths, the MCP tools, and five documents, beside S-0514's 2h30m for a similar spread. Delivery is flai's 2026-10-29T11:19Z moved by the added 2h45m, 37th in the pull order with an in-progress limit of 5.

Cost of delay: `flai cod` gave 37.50 USD a week from the input the orchestrator set on TH-0706 (15m of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week). The value stands as computed; the input is not changed.
