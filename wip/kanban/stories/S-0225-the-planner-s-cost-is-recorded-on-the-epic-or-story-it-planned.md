---
id: S-0225
type: story
nature: improvement
title: The planner's cost is recorded on the epic or story it planned
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-02T11:54:44Z
transitions: []
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md]
after: [S-0208]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0225 The planner's cost is recorded on the epic or story it planned

## Goal

ADR-0051 records what story agents spend on the items they work. The planner has no story, so what it spends would be lost to the items' totals and to the charts. Its cost and seconds should land on the epic or story it planned, marked as strategic so the item's own agent cost stays readable.

## Acceptance criteria
- [ ] When a planner activity ends, its usage (tokens, cost, seconds, from the run's log apportioned by activity as ADR-0051 apportions tasks) is added to the item's `usage` under a `strategic` entry per agent kind, `estimated: true`, and summed up the hierarchy as other usage is
- [ ] `flai stats` and the cost charts show strategic cost per item and in totals as its own series, and leave it out of the per-model agent figures
- [ ] An item that gets a planner forecast or estimate before any agent works it shows the planner's expected cost beside it, from the forecast duration and the project's mean cost per agent hour, marked as an estimate
- [ ] `design/system/metrics.md` records the `strategic` usage entry (the ADR of the metrics story covers it); tests pin the apportioning

## Tasks

## Notes
