---
id: S-0217
type: story
nature: feature
title: "flai exposes the orchestrator's deterministic operations as commands: ordering by policy, promotion candidates, and release evaluation"
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T07:09:05Z
transitions:
  - to: ready
    at: 2026-10-03T20:34:00Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:09Z
    by: alex
  - to: ready
    at: 2026-10-05T04:41:11Z
    by: alex
  - to: in-progress
    at: 2026-10-05T06:07:02Z
    by: agent-S-0217
  - to: review
    at: 2026-10-05T07:07:35Z
    by: agent-S-0217
  - to: done
    at: 2026-10-05T07:09:05Z
    by: alex
tags: [flai]
touches: [flai/cmd, flai/internal/workitem, flai/internal/release, flai/internal/hostapi, flai/internal/mcpserver, flai/internal/manifest, flai/internal/guard, flai/internal/check/check.go, flai/internal/check/orchestration_test.go, flaiover/src/lib/server/agent.ts, design/system/flai-cli.md, design/system/project-manifest.md, design/system/server-performance.md, design/issues, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [S-0199, S-0210]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4073
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 554
      output: 8472
      cache_read: 22537056
      cache_write: 792999
      cost: 20.1038
    - model: claude-sonnet-5
      input: 62
      output: 15349
      cache_read: 1875334
      cache_write: 155702
      cost: 0.8499
cost_of_delay:
  value: 97.56
  by: planner-E-0016
  at: 2026-10-04T04:48:17Z
forecast:
  duration: 2h
  delivery: 2026-10-05T09:02:00Z
  basis: "Its own forecast of 2h; 1st in the pull order with an in-progress limit of 3, with nothing ahead of it."
  by: flai
  at: 2026-10-05T05:59:10Z
---
# S-0217 flai exposes the orchestrator's deterministic operations as commands: ordering by policy, promotion candidates, and release evaluation

## Goal

The designer chose (2026-10-02) that the orchestrator's arithmetic lives in flai, testable and the same in the dashboard, and that the long-running agent makes only judgement calls. These commands are that arithmetic, usable by the operator without any agent.

## Acceptance criteria
- [x] `flai order --by cod|wsjf|throughput|fifo [--apply]` computes the ready column's order by the policy (cost of delay value; value over forecast duration; shortest forecast first; creation order) and prints it with the figures; `--apply` writes `board.md`'s order
- [x] `flai promote --candidates [--limit]` lists backlog stories that could move to ready: finalized, not held, with forecast and value, ordered by the project's policy, and says why each other story is not a candidate
- [x] `flai release --evaluate` says whether the release policy is met: `threshold` (unreleased cost of delay value or count of accepted stories over the manifest's figure), `theme` (every story of a named epic or tag accepted), or `judgement` (never met by itself), with the figures
- [x] Each has a hostapi read and an MCP tool; the policy and thresholds are manifest settings (`orchestration.policy`, `orchestration.release`)
- [x] `design/system/flai-cli.md` and the user guide describe them; tests pin each ordering and evaluation on fixtures

## Tasks
- T-0809 system-flow.yaml takes orchestration.policy and orchestration.release, and flai check reports a bad one
- T-0810 flai order --by computes the ready column's order by cod, wsjf, throughput, or fifo with its figures, and --apply writes it
- T-0811 flai promote --candidates lists the backlog stories that could go to ready, by the project's policy, and why each other one cannot
- T-0812 flai release --evaluate says whether the release policy is met, with its figures
- T-0813 The policy order, the promotion candidates, and the release evaluation each have a hostapi read and an MCP tool that the guard passes as reads
- T-0814 flai-cli.md, the user guide, and the reference describe flai order --by, flai promote --candidates, and flai release --evaluate

## Notes
