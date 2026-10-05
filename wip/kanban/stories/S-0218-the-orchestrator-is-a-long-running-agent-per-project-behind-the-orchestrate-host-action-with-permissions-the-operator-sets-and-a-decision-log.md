---
id: S-0218
type: story
nature: feature
title: The orchestrator is a long-running agent per project behind the orchestrate host action, with permissions the operator sets and a decision log
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T05:23:10Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:14Z
    by: alex
tags: [flai, dashboard]
topics: [orchestration]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/internal/mcpserver, flai/internal/config, ".claude/agents", template, flai/internal/manifest, flai/internal/guard, flai/cmd, flaiover/src/routes/activity, design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/operators, docs/users/flai.md, ".claude/settings.json", flai/internal/check/check.go, flai/internal/check/orchestration_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flaiover/src/lib/activity.ts, flaiover/src/lib/server/agent.ts, flaiover/src/lib/server/agent.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai-reference.md]
after: [S-0206, S-0207, S-0217]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 102.04
  by: planner-S-0218
  at: 2026-10-05T04:48:03Z
forecast:
  duration: 1h40m
  delivery: 2026-10-05T10:57:00Z
  basis: "Its own forecast of 1h40m; 4th in the pull order with an in-progress limit of 3, behind S-0258, S-0276 and S-0217."
  by: flai
  at: 2026-10-05T05:23:10Z
---
# S-0218 The orchestrator is a long-running agent per project behind the orchestrate host action, with permissions the operator sets and a decision log

## Goal

The orchestrator is an agentic operator: it keeps work moving between backlog, ready, and the agents, answers threads, and decides releases, each only as far as the operator permits. It runs for a project as long as `orchestrate` is enabled, waiting on events between decisions, and every decision is logged with its reason.

## Acceptance criteria
- [ ] An `orchestrate` host action, off by default, starts one orchestrator run per project (`orchestration.agent` gives harness, model, config); `flai serve` restarts it when it ends and the action is still on, and stops it when the action is turned off
- [ ] Permissions live in the manifest under `orchestration.permissions`, each off by default: `plan_backlog_epics` (ask the planner to draft stories for backlog epics), `finalize_drafts`, `promote_to_ready`, `order_ready`, `answer_threads` (`recommend` or `autonomous`), `accept_reviews`, `publish`; and `orchestration.policy` (`throughput` or `cost_of_delay`)
- [ ] `flai guard` enforces the permissions on the orchestrator's calls: a call outside them is refused with the permission that would allow it, and the refusal is logged
- [ ] Its prompt: prime with `--role orchestrate`, read the board and inbox, act within permissions using the deterministic commands, log each decision with `activity_log` (what, why, which policy figure), then `wait_for_events` and repeat
- [ ] Each decision and refusal is in the orchestrator's activity document and visible in the dashboard
- [ ] `design/system/strategic-agents.md`, `flai-cli.md`, and the operator guide describe it; tests cover start and stop with the action, a permitted and a refused call

## Tasks
- T-0881 system-flow.yaml takes orchestration.permissions and orchestration.agent, and flai check reports a bad one
- T-0882 The orchestrator's prompt, its claude-code definition, and settings that run flai guard on its file edits
- T-0883 flai serve runs one orchestrator per project while the orchestrate host action is on, restarts it when it ends, and stops it when the action is turned off
- T-0887 flai guard holds an orchestrator session to its permissions, names the permission a refused call needs, and logs the refusal in orchestrator.md
- T-0889 The dashboard shows the orchestrator's run on the Activity page and the orchestrate action on the Settings page
- T-0892 An ADR, strategic-agents.md, flai-cli.md, the manifest design, and the operator and user guides describe the orchestrator

## Notes

### Planning

Planned by planner-S-0218 on 2026-10-05.

Touches:

- Declared, kept: `flai/internal/serve`, `flai/internal/harness`, `flai/internal/hostapi`, `flai/internal/mcpserver`, `flai/internal/config`, `.claude/agents`, `template`, `flai/internal/manifest`, `flai/internal/guard`, `flai/cmd`, `flaiover/src/routes/activity`, `design/adrs`, `design/system/strategic-agents.md`, `design/system/flai-cli.md`, `design/system/project-manifest.md`, `docs/operators`, `docs/users/flai.md`.
- Design (ADR-0082 § 7, the planner's precedent): `.claude/settings.json`, whose `Edit|Write|NotebookEdit` hook runs the guard only for `FLAI_ROLE=plan` and must cover `orchestrate` too.
- Co-change and layout: `flai/internal/check/check.go` and `flai/internal/check/orchestration_test.go` (6% of the seed commits; `manifest.orchestration` findings, as S-0217's T-0809 has them); `flai/internal/workitem/activity.go` and its test (a refusal's entry in `orchestrator.md`); `flaiover/src/lib/activity.ts` and `flaiover/src/lib/server/agent.ts` with its test (5%; S-0208 changed both for the planner's run); `design/system/flaiover-dashboard.md` (16%), `docs/users/flaiover.md` (13%), and `docs/users/flai-reference.md` (15%).
- Left out: the co-change list's issues, conventions, and workflow documents, which S-0208's history shows but this story's criteria do not reach.

Forecast: 1h40m, against flai's 1h15m (132 s per unit over 14 large feature stories on this model, times size 34). Adjusted up because S-0208, the planner's host action and the closest precedent, took 65m of agent time, and this story adds a permission model the guard reads from the manifest and a run that restarts and stops with its action, about half as much again. Delivery 10:52Z: it starts after S-0217, about 08:20Z, and takes 1h40m times the cycle factor 1.52.

Cost of delay: 102.04 USD a week, flai's figure, kept. The story has no inputs of its own, so it takes its share of E-0016's 1500 USD a week, 1h40m of 24h30m forecast over the epic's 17 open stories without inputs. It replaces planner-E-0016's 121.95, which was worked out over more open stories and a 2h30m duration.

Plan: three layers. T-0881 and T-0882 wait for nothing; T-0883 waits for both and T-0887 for T-0881; T-0889 waits for T-0883 and T-0892 for the four code tasks. The plan's thread on this story gives the assumptions.
