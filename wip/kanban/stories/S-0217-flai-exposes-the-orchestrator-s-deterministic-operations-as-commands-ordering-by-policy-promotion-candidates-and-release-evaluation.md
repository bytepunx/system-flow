---
id: S-0217
type: story
nature: feature
title: "flai exposes the orchestrator's deterministic operations as commands: ordering by policy, promotion candidates, and release evaluation"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-04T04:48:17Z
transitions:
  - to: ready
    at: 2026-10-03T20:34:00Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:09Z
    by: alex
tags: [flai]
touches: [flai/cmd, flai/internal/workitem, flai/internal/release, flai/internal/hostapi, flai/internal/mcpserver, flai/internal/manifest, flai/internal/guard, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [S-0199, S-0210]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 97.56
  by: planner-E-0016
  at: 2026-10-04T04:48:17Z
forecast:
  duration: 2h
  delivery: 2026-10-04T12:00:00Z
  basis: "Three commands, each with a hostapi read and an MCP tool and manifest settings, extending flai order and flai release; sized above S-0205 (4208 s); serial after S-0210 on flai/cmd."
  by: planner-E-0016
  at: 2026-10-04T04:44:32Z
---
# S-0217 flai exposes the orchestrator's deterministic operations as commands: ordering by policy, promotion candidates, and release evaluation

## Goal

The designer chose (2026-10-02) that the orchestrator's arithmetic lives in flai, testable and the same in the dashboard, and that the long-running agent makes only judgement calls. These commands are that arithmetic, usable by the operator without any agent.

## Acceptance criteria
- [ ] `flai order --by cod|wsjf|throughput|fifo [--apply]` computes the ready column's order by the policy (cost of delay value; value over forecast duration; shortest forecast first; creation order) and prints it with the figures; `--apply` writes `board.md`'s order
- [ ] `flai promote --candidates [--limit]` lists backlog stories that could move to ready: finalized, not held, with forecast and value, ordered by the project's policy, and says why each other story is not a candidate
- [ ] `flai release --evaluate` says whether the release policy is met: `threshold` (unreleased cost of delay value or count of accepted stories over the manifest's figure), `theme` (every story of a named epic or tag accepted), or `judgement` (never met by itself), with the figures
- [ ] Each has a hostapi read and an MCP tool; the policy and thresholds are manifest settings (`orchestration.policy`, `orchestration.release`)
- [ ] `design/system/flai-cli.md` and the user guide describe them; tests pin each ordering and evaluation on fixtures

## Tasks

## Notes
