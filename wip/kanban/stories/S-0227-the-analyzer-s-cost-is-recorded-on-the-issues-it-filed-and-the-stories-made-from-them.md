---
id: S-0227
type: story
nature: improvement
title: The analyzer's cost is recorded on the issues it filed and the stories made from them
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-04T04:48:23Z
transitions: []
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/issues, flai/internal/metrics, design/system/continuous-improvement.md, design/adrs]
after: [S-0223]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.98
  by: planner-E-0016
  at: 2026-10-04T04:48:23Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T13:00:00Z
  basis: "Issues carry no usage today, so this adds a usage field to issues with an ADR and carries it to the story made from an issue, beyond S-0225's generic strategic entry; like S-0203 (1555 s) doubled."
  by: planner-E-0016
  at: 2026-10-04T04:45:02Z
---
# S-0227 The analyzer's cost is recorded on the issues it filed and the stories made from them

## Goal

ADR-0051 records what story agents spend on the items they work. The analyzer has no story, so what it spends would be lost to the items' totals and to the charts. Its cost and seconds should land on the issues it filed and the stories made from them, marked as strategic so the item's own agent cost stays readable.

## Acceptance criteria
- [ ] When a analyzer activity ends, its usage (tokens, cost, seconds, from the run's log apportioned by activity as ADR-0051 apportions tasks) is added to the item's `usage` under a `strategic` entry per agent kind, `estimated: true`, and summed up the hierarchy as other usage is
- [ ] `flai stats` and the cost charts show strategic cost per item and in totals as its own series, and leave it out of the per-model agent figures
- [ ] An item that gets a analyzer forecast or estimate before any agent works it shows the planner's expected cost beside it, from the forecast duration and the project's mean cost per agent hour, marked as an estimate
- [ ] `design/system/metrics.md` records the `strategic` usage entry (the ADR of the metrics story covers it); tests pin the apportioning

## Tasks

## Notes
