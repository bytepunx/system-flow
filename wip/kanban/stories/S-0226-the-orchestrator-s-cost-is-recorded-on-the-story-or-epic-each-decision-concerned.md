---
id: S-0226
type: story
nature: improvement
title: The orchestrator's cost is recorded on the story or epic each decision concerned
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-04T04:48:22Z
transitions: []
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/mcpserver, design/system/strategic-agents.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 36.59
  by: planner-E-0016
  at: 2026-10-04T04:48:22Z
forecast:
  duration: 45m
  delivery: 2026-10-05T11:00:00Z
  basis: "S-0225 already makes the strategic usage entry generic per kind and shows it in stats and charts, so what is left is charging each orchestrator activity to the items it names, else the project; small like S-0202 (1197 s)."
  by: planner-E-0016
  at: 2026-10-04T04:45:00Z
---
# S-0226 The orchestrator's cost is recorded on the story or epic each decision concerned

## Goal

ADR-0051 records what story agents spend on the items they work. The orchestrator has no story, so what it spends would be lost to the items' totals and to the charts. Its cost and seconds should land on the story or epic each decision concerned, and the project when none, marked as strategic so the item's own agent cost stays readable.

## Acceptance criteria
- [ ] When a orchestrator activity ends, its usage (tokens, cost, seconds, from the run's log apportioned by activity as ADR-0051 apportions tasks) is added to the item's `usage` under a `strategic` entry per agent kind, `estimated: true`, and summed up the hierarchy as other usage is
- [ ] `flai stats` and the cost charts show strategic cost per item and in totals as its own series, and leave it out of the per-model agent figures
- [ ] An item that gets a orchestrator forecast or estimate before any agent works it shows the planner's expected cost beside it, from the forecast duration and the project's mean cost per agent hour, marked as an estimate
- [ ] `design/system/metrics.md` records the `strategic` usage entry (the ADR of the metrics story covers it); tests pin the apportioning

## Tasks

## Notes
