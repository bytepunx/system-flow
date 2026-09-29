---
id: S-0143
type: story
nature: feature
title: Stories and Tasks track tokens and cost
status: done
parent: E-0011
owner: alex
created: 2026-09-29T00:24:52Z
updated: 2026-09-29T06:43:49Z
transitions:
  - to: ready
    at: 2026-09-29T03:18:02Z
    by: alex
  - to: in-progress
    at: 2026-09-29T06:03:20Z
    by: agent-S-0143
  - to: review
    at: 2026-09-29T06:32:15Z
    by: agent-S-0143
  - to: done
    at: 2026-09-29T06:43:49Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, flai/internal/usage, flai/internal/workitem, flai/internal/serve, flai/internal/metrics, design/system, docs, design/adrs, design/tech/charts.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0143 Stories and Tasks track tokens and cost

## Goal

Epics, Stories, and Tasks gain frontmatter to include tokens used and cost or estimated cost of the expenditure. Charts use this to show users the following (all should group by agent model doing the work):

- token rates (epic, task, and story) per hour
- cost or estimated cost (epic, task, and story)
- epic, story, and task completion as a function of time and cost

## Acceptance criteria
- [x] epics, stories, and tasks track tokens spent
- [x] epics, stories, and tasks track cost
- [x] when any item (epic, story, task) moves into done, it triggers a cascade up the hierarchy to update its parent's aggregate front-matter tracking

## Tasks
- T-0539 flai measures an agent's tokens and cost from its log
- T-0540 Work items carry usage, and a move into done sums it up the hierarchy
- T-0541 flai serve records usage from its agents' logs, and flai serve agent usage records it by hand
- T-0542 flai stats reports tokens and cost by model
- T-0543 The dashboard charts token rate, cost, and completion against time and cost by model

## Notes

- Decided in ADR-0051 (written as accepted; accepting the story accepts it). Items carry `usage`: `source` (`log` or `sum`), `seconds` of agent work, `estimated`, and per model `input`, `output`, `cache_read`, `cache_write` tokens and `cost` in US dollars.
- Where the numbers come from: the Claude Code stream-json logs `flai serve` keeps for the agents it starts. Each session's newest `result` is its reported total; a task gets its story's session totals in the share of the input and cache tokens of the calls made while it was in progress, marked estimated; a run with no result is priced at the rate the logs report, marked estimated. No price table.
- Who writes it: `flai serve`, when an agent it started ends and when a task of a running agent's story enters done; `flai serve agent usage --write` by hand (`--all` for stories worked before). Every move into done (`flai move`, MCP `item_move`, `flai accept`) sums each unmeasured ancestor from its children, up to the epic.
- Verified: criteria 1 and 2 by `flai serve agent usage --all --write` on a scratch clone of this repository against the host's real logs (44 items written, `flai check --strict` clean, `flai stats` totals by model); criterion 3 by `internal/workitem/usage_test.go` and `internal/serve/usage_test.go`.
- Not measured: sessions run by hand, other harnesses. A flai older than this refuses items with `usage`: upgrade the host's flai before it measures anything, and before running the backfill in this repository.
