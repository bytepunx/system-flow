---
id: S-0228
type: story
nature: feature
title: The Workflow menu has Orchestrator and Analyzer pages showing their status, activity log, and runs
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-07T00:51:08Z
transitions:
  - to: ready
    at: 2026-10-06T22:46:46Z
    by: alex
  - to: in-progress
    at: 2026-10-06T23:58:30Z
    by: alex
  - to: review
    at: 2026-10-07T00:49:39Z
    by: agent-S-0228
  - to: done
    at: 2026-10-07T00:51:08Z
    by: alex
tags: [dashboard]
topics: [orchestration, analysis]
touches: [flaiover/src/routes, flaiover/src/lib/components, flaiover/src/routes/api, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/sitemenu.ts, flai/internal/hostapi, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/planner.ts, flaiover/src/lib/planner.test.ts, flaiover/src/lib/strategic.ts, flaiover/src/lib/strategic.test.ts, flaiover/src/lib/activity.ts, flai/internal/serve/stream.go, flai/internal/serve/stream_test.go, flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/agents.go, flai/cmd/serve.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, design/system/flai-cli.md, design/system/dashboard-host-channel.md, docs/users/flai-reference.md, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/agents_test.go, flaiover/src/routes/api/orchestrator, flaiover/src/routes/api/analyzer, flaiover/src/routes/api/agent-stream, flaiover/src/routes/workflow/orchestrator, flaiover/src/routes/workflow/analyzer, flaiover/src/lib/components/AgentStream.svelte, flaiover/src/lib/components/AgentStream.svelte.test.ts, flaiover/src/lib/server/agent.test.ts, design/system/strategic-agents.md, docs/users/flai.md, docs/operators/index.md]
after: [S-0208, S-0218, S-0223, S-0259]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2834
  models:
    - model: claude-opus-5-5
      input: 624
      output: 277460
      cache_read: 42611536
      cache_write: 1031465
      cost: 19.8613
    - model: claude-sonnet-5-5
      input: 14
      output: 4094
      cache_read: 171868
      cache_write: 35915
      cost: 0.1651
  strategic:
    - kind: orchestrator
      seconds: 27
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 256
          cache_read: 532936
          cache_write: 3381
          cost: 0.238
cost_of_delay:
  value: 60.81
  by: planner-S-0228
  at: 2026-10-05T05:44:36Z
forecast:
  duration: 1h
  delivery: 2026-10-07T01:17:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 3, behind S-0273."
  by: flai
  at: 2026-10-06T23:57:14Z
---
# S-0228 The Workflow menu has Orchestrator and Analyzer pages showing their status, activity log, and runs

## Goal

The orchestrator and the analyzer are invisible without a place in the dashboard. Each gets a page under the Workflow menu, laid out like S-0259's Planner page. The page shows whether the agent's host action is enabled and running, its activity document newest first, its runs with cost and outcome, and the actions the operator can take: run the analyzer, or start and stop the orchestrator.

## Acceptance criteria
- [x] `/workflow/orchestrator` and `/workflow/analyzer` are in the two-tier site menu's Workflow tier, beside the Planner page. Each shows the host action's state with the command that enables it, the current run if any (the agent pane that story pages have), the activity log's entries with duration and cost, the front matter totals, and a list of past runs
- [x] The analyzer page has a Run button with the focus. The orchestrator page has a Stop button and the last decisions with their reasons
- [x] Pages update live from `/api/events` as the activity documents change
- [x] `design/system/flaiover-dashboard.md` and the user guide describe them; tests cover each page's data and actions

## Tasks
- T-0917 flai serves the orchestrator's and the analyzer's run streams, and stops and starts the orchestrator on the operator's word
- T-0921 The dashboard's API answers the orchestrator's and the analyzer's state, activity, and runs, and runs the analyzer and stops and starts the orchestrator
- T-0928 The Workflow menu's Orchestrator page shows its state, current run, decisions, and runs, and stops and starts it
- T-0931 The Workflow menu's Analyzer page shows its state, current run, activity, and runs, and runs it with a focus
- T-0938 The design, the user guide, and the flai reference describe the Orchestrator and Analyzer pages and the orchestrator's stop and start

## Notes

- Split on TH-0107: the Planner page is S-0259, which ships first.

### Planning

Planned by planner-S-0228 on 2026-10-05. The plan's thread on S-0228 names the tasks, layers, and assumptions.

Touches, where each came from:

- Declared, kept: `flaiover/src/routes`, `flaiover/src/lib/components`, `flaiover/src/routes/api`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `flaiover/src/lib/sitemenu.ts`, `flai/internal/hostapi`.
- Co-change (`flai touches suggest`, 289 of 866 commits): `design/system/flai-cli.md` (27%), `flaiover/src/lib/server/agent.ts` (11%), `docs/users/flai-reference.md` (9%), `flai/cmd/serve_actions.go` and its test (8%), `design/system/dashboard-host-channel.md` (6%), `flai/internal/serve/agents.go` and `flaiover/src/lib/activity.ts` (4%). These are where the orchestrator's stop and start, and the run streams by role, are served, typed, and described.
- Design: `flai/internal/serve/stream.go` and its test, where S-0259's T-0836 read a planner run's stream and where a run by role is read; `flai/internal/serve/orchestrate.go` and its test, which S-0218's T-0883 creates and where the hold on Stop lives.
- Layout: `flai/cmd/serve.go` for the `flai serve orchestrate stop|start` subcommand; `flaiover/src/lib/planner.ts` and its test, and a new `flaiover/src/lib/strategic.ts` and its test, for what the three pages share; `flaiover/src/lib/sitemenu.test.ts` for the menu entries.
- Left out: `docs/users/flai.md` (25%) and `docs/operators/index.md` (19%) co-change often, but no host action is added, only two writes under `orchestrate`; `design/issues/summary.md` and `design/adrs/README.md` are incidental co-changes.

Figures:

- Forecast 1h, delivery 2026-10-05T20:49Z. `flai forecast` said 14m on the 7 declared touches and 1h4m on the 25 predicted ones (132 s per unit of size over 14 done large feature stories on claude-opus-5-5). Set to 1h: S-0259, one page of this kind with its flai reads, took 29m in progress, and this story is two pages plus the orchestrator's stop and start in flai serve. The delivery is flai's play-out start, 17th in the pull order, plus 1h.
- Cost of delay 60.81 USD a week, from `flai cod`: no inputs of its own, so its share of E-0016's 1500 USD a week, 1h of the epic's open forecast. It was 60.98 from planner-E-0016 and moves only with the forecast.
