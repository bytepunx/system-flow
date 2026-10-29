---
id: S-0534
type: story
nature: improvement
title: A usage-limit refusal sets a budget state on the host that holds every start on that provider until the reset or a backoff, probes once, and reaches every inbox, the journal, the dashboard, and the ADR-0205 channel
status: in-progress
parent: E-0030
owner: alex
created: 2026-10-28T21:02:16Z
updated: 2026-10-29T05:47:10Z
transitions:
  - to: ready
    at: 2026-10-28T22:21:09Z
    by: alex
  - to: in-progress
    at: 2026-10-29T05:47:10Z
    by: agent-S-0534
tags: [flai, serve, budget, rate-limit]
topics: [orchestration, cli, dashboard]
touches: [flai/internal/serve, flai/internal/mcpserver, flaiover/src, design/system/flai-cli.md, design/system/project-manifest.md, design/system/flaiover-dashboard.md, docs, flai/internal/serve/budget.go, flai/internal/serve/budget_test.go, flai/internal/serve/serve.go, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/curate.go, flai/internal/serve/curate_test.go, flai/internal/serve/testdata/usage-limit-result.jsonl, flai/internal/serve/start.go, flai/internal/serve/start_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/notify.go, flai/internal/serve/notify_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/inbox_test.go, flai/cmd/mcp.go, flai/cmd/mcp_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/contract_test.go, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/BudgetBanner.svelte, flaiover/src/lib/components/BudgetBanner.svelte.test.ts, flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, design/system/strategic-agents.md, design/system/agent-narrative.md, docs/operators/index.md, docs/operators/runbooks/backup.md, docs/users/flai.md, docs/users/flaiover.md]
after: [S-0532]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 553
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-29
      hand_edits: 1
      work: 21
  models:
    - model: claude-opus-5-5
      input: 88
      output: 503
      cache_read: 4434790
      cache_write: 308986
      cost: 1.622
  strategic:
    - kind: planner
      seconds: 724
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 148
          output: 530
          cache_read: 9418517
          cache_write: 225102
          cost: 0.6314
        - model: claude-opus-5-5
          input: 46
          output: 245
          cache_read: 3735627
          cache_write: 221315
          cost: 0.9426
    - kind: orchestrator
      seconds: 101
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 283
          cache_read: 2909882
          cache_write: 18986
          cost: 0.6446
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h
    by: orchestrator
    at: 2026-10-28T21:17:38Z
  value: 150
  by: planner-S-0534
  at: 2026-10-28T21:28:34Z
forecast:
  duration: 4h30m
  delivery: 2026-10-29T11:20:00Z
  basis: "Its own forecast of 4h30m; 2nd in the pull order with an in-progress limit of 5, behind S-0417 and S-0537."
  by: flai
  at: 2026-10-29T05:16:39Z
---
# S-0534 A usage-limit refusal sets a budget state on the host that holds every start on that provider until the reset or a backoff, probes once, and reaches every inbox, the journal, the dashboard, and the ADR-0205 channel

## Goal

ADR-0209 §5. When the API refuses a run with the account's usage limit, the run ends with exit 1 and flai serve starts it again a minute later (`orchestrateRetry`, ADR-0087 §1), or up to `agent.auto_restarts` times for a story's agent (ADR-0108), with no backoff and no knowledge of why. Each of the five projects the host serves learns of the limit by failing in turn, and nothing tells the planner, the curator, the analyzer, or the operator. This story gives flai serve one budget state per provider on the host, kept in the serve folder beside `upgrade-asks.json`: set on a usage-limit refusal (Claude Code's `api_error: usage_limit_reached`, with the reset time it reports, or a backoff doubling from `5m` to `1h` when it reports none), it holds every start of any role and any story's agent on that provider across every project the host serves, counts none of them as a failed run, starts one probe at the reset time or after the backoff, clears on the probe's success and extends on its refusal, and is a journal entry, a `budget` field in every agent's `inbox` with the reset time, a banner on the dashboard, and a reminder through the ADR-0205 channel when one is set. Reasons queued for the orchestrator (the following story) and stories' agents wait where they are.

## Acceptance criteria
- [ ] A run ending with the usage-limit refusal sets the provider's budget state with the reset time the API reported, or the backoff, and `flai serve journal` records it; a test reproduces the refusal from a recorded Claude Code result line
- [ ] While the state is set, no role's run and no story's agent is started on that provider on any served project, none counts against `auto_restarts` or `orchestrateRetry`, and a start the operator asks for by hand says why it is held
- [ ] One probe run starts at the reset time or after the backoff; success clears the state and starts what waited; refusal doubles the backoff up to `1h`
- [ ] `inbox` carries `budget` with the provider and the reset time while the state is set, and the MCP instructions tell an agent to hold rather than retry
- [ ] The dashboard shows a banner on every page naming the provider and the reset time, and the ADR-0205 channel gets one reminder per state change
- [ ] The design (`flai-cli.md`, `project-manifest.md`, `flaiover-dashboard.md`), the serve docs, and the operator docs describe the state, its file, and the probe

## Tasks
- T-2282 flai serve keeps a budget state per provider in budget.json in the serve folder: set from a usage limit, held, probed once, cleared, or extended, each change journaled
- T-2283 A story's agent or a role's run that ends on the usage limit sets its provider's budget state and is not counted a failed run against auto_restarts or orchestrateRetry
- T-2284 While its provider's budget state is set, no story's agent is started or restarted on any served project, and a start or restart asked for by hand says it is held and until when
- T-2285 While its provider's budget state is set, no planner, orchestrator, analyzer, or curator run starts, on its trigger, schedule, or retry, and a run asked for by hand says it is held and until when
- T-2286 inbox carries budget with the provider, the limit, and the reset time while a provider's budget state is set, and the MCP instructions tell an agent to hold rather than retry
- T-2287 The host API's agent.status and flai serve status carry each provider's budget state, with its limit, until, and probe
- T-2288 The ADR-0205 channel gets one notice per budget state change, set, extended, or cleared, with the provider, the limit, and until, and none when no channel is set
- T-2289 The dashboard shows a banner on every page while a provider's budget state is set, naming the provider, the limit, and the reset time in the viewer's time
- T-2290 The design, the serve docs, and the operator docs describe the budget state, budget.json, the hold, the probe, the inbox field, the banner, and the notice

## Notes

Evidence: ADR-0209's context. The account's weekly limit was reached on Oct 14, 20, and 25, the five-hour session limit on Oct 28 after three busy hours. A provider is a name set with `flai serve agent provider` (ADR-0129); the state is per provider because the limit is the account's, not the project's.

### Planning

Planned by planner-S-0534 on 2026-10-28. The tasks, their layers, and the assumptions are on the plan's thread.

The design: no usage-limit reading exists in flai today. S-0514 adds it (`flai/internal/usage`, `flai/internal/serve/limit.go`), and S-0532, which this story names in `after`, waits for S-0514, so both are on main before this story starts. The state is a new file in the serve folder, `budget.json`, written as `upgrade-asks.json` and `reminders.json` are (`serve.go`, `remind.go`). Every launcher reads it, so one refusal holds every served project. The MCP server reads it beside the host configuration, as `flai serve` finds it (`serve.DirFor`). The dashboard gets it on `agent.status`.

Touches, and where each came from:

- `flai/internal/serve`, `flai/internal/mcpserver`, `flaiover/src`, `docs`: declared folder touches, kept. The tasks name every file the plan expects below them, so the story's claim narrows each folder to those files (ADR-0096). `flai/internal/serve` also stays because T-2285 holds the start path S-0417 adds for a project's own role, a file no task can name until S-0417 is on main.
- `design/system/flai-cli.md`, `project-manifest.md`, `flaiover-dashboard.md`: declared; criterion 6.
- `flai/internal/serve/budget.go`, `budget_test.go`, `testdata/usage-limit-result.jsonl`, `flaiover/src/lib/components/BudgetBanner.svelte`, `BudgetBanner.svelte.test.ts`: layout; new files.
- `flai/internal/serve/serve.go`: layout; the serve folder's paths.
- `flai/internal/serve/agents.go`, `start.go`, `restart.go` and their tests: layout; a story's agent's start paths, `judgeRun`, `restartEnded`, and the hand starts.
- `flai/internal/serve/plan.go`, `replan.go`, `orchestrate.go`, `orchestrate_hold.go`, `analyze.go`, `analysis_schedule.go`, `curate.go`, `curate_schedule.go` and their tests: layout; each role's start, end, retry, schedule, and hand start.
- `flai/internal/serve/notify.go`, `notify_test.go`: design (ADR-0205 §3); the channel's sender.
- `flai/internal/inbox/inbox.go`, `inbox_test.go`, `flai/internal/mcpserver/server.go`, `folder.go` and their tests, `mcpserver/inbox_test.go`, `flai/cmd/mcp.go`, `mcp_test.go`: layout; the `Inbox` struct, the MCP instructions, and the options the server is given.
- `flai/cmd/serve_actions.go` (co-change, 5%), `serve_actions_test.go` (co-change, 4%), `flai/internal/hostapi/hostapi.go`, `hostapi_test.go`, `contract_test.go`: layout; `agent.status` and `flai serve status`.
- `flaiover/src/lib/server/agent.ts`, `agent.test.ts`, `routes/api/agent/+server.ts`, `agent.test.ts`, `lib/hostflai.svelte.ts`, `routes/+layout.svelte`, `layout.svelte.test.ts`: layout; the status the dashboard polls and where the other banners sit.
- `design/system/strategic-agents.md` (co-change, 3%), `design/system/agent-narrative.md` (co-change, 3%): the roles' run outcomes and what an agent is told through MCP.
- `docs/operators/index.md`, `docs/operators/runbooks/backup.md`, `docs/users/flai.md`, `docs/users/flaiover.md`: design; the serve files, the host actions, the reminders, and the dashboard's banners.
- Left out of the suggestion: `design/adrs/README.md`, `design/system/workflow.md`, `design/issues/summary.md`, `template/CHANGELOG.md`, the work-management convention and its template copy, `flai/internal/hostapi/writes.go` and its test, `flai/cmd/root.go`, `flai/cmd/serve.go`, `flai/internal/check/check.go`, `work-hierarchy.md`, the harness, the manifest, and `dashboard-host-channel.md`. They are co-change of past serve stories, but nothing in them changes here. No ADR is added, since ADR-0209 §5 decides this. No manifest key, check, template file, or command is added. The host API's writes only build argv; the refusals live in the serve functions.

Forecast: `flai forecast` gave 20m (median 89 s per unit of size over 100 done improvement stories on claude-opus-5-5, size 13: 6 criteria, 7 touches). Raised to 4h30m: nine tasks across a new host state, five run-end paths, every start path, the MCP inbox, the host API, a dashboard banner, the notify channel, and nine documents. S-0532 has 3h for seven tasks of a similar spread, and S-0514 has 2h30m. Delivery is flai's 2026-10-29T15:49Z moved by the added 4h10m; the story is 38th in the pull order with an in-progress limit of 5.

Cost of delay: `flai cod` gave 150.00 USD a week from the input the orchestrator set on TH-0710: 1h of time lost per 168h cycle at 150 USD an hour, 1.00 cycles a week. The value stands as computed, and the input is not changed.
