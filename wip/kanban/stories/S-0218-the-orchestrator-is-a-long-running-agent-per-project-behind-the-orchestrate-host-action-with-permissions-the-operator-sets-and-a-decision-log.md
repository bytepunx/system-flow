---
id: S-0218
type: story
nature: feature
title: The orchestrator is a long-running agent per project behind the orchestrate host action, with permissions the operator sets and a decision log
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T00:33:39Z
transitions: []
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/harness, flai/internal/hostapi, flai/internal/mcpserver, flai/internal/config, ".claude/agents", template, flai/internal/manifest, flai/internal/guard, flai/cmd, flaiover/src/routes/activity, design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/operators, docs/users/flai.md]
after: [S-0206, S-0207, S-0217]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 121.95
  by: planner-E-0016
  at: 2026-10-04T04:48:17Z
forecast:
  duration: 2h30m
  delivery: 2026-10-05T10:11:00Z
  basis: "Its own forecast of 2h30m; 16th in the pull order with an in-progress limit of 3, behind S-0249, S-0266, S-0268, S-0253, S-0257, S-0244, S-0262, S-0258, S-0260, S-0212, S-0213, S-0214, S-0215, S-0216 and S-0217."
  by: flai
  at: 2026-10-05T00:33:39Z
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

## Notes
