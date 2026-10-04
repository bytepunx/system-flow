---
id: S-0225
type: story
nature: improvement
title: The planner's cost is recorded on the epic or story it planned
status: in-progress
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-04T04:20:34Z
transitions:
  - to: ready
    at: 2026-10-04T00:41:02Z
    by: alex
  - to: in-progress
    at: 2026-10-04T03:59:44Z
    by: agent-S-0225
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, flai/internal/metrics, flai/cmd/stats.go, flai/cmd/show.go, design/system/metrics.md, design/system/work-hierarchy.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md, design/adrs/README.md, docs/users/flai.md, docs/users/flaiover.md, flaiover/src/routes/items, flaiover/src/lib/usage.ts, flaiover/src/lib/viz, flai/cmd/activity.go, flai/cmd/check_stats_test.go]
after: [S-0208]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1312
  estimated: true
  models:
    - model: claude-haiku-4-5-20251001
      input: 364
      output: 59
      cache_read: 2510303
      cache_write: 173513
      cost: 0.7371
    - model: claude-opus-5-5
      input: 346
      output: 2494
      cache_read: 20443237
      cache_write: 535759
      cost: 8.5571
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
- T-0803 Usage carries a strategic entry per agent kind, apart from the agents' figures, charged to an item and the items above it
- T-0804 Usage carries a strategic entry per agent kind, apart from the agents' figures, charged to an item and the items above it
- T-0805 A planner activity's apportioned usage is charged to the item it planned and the items above it
- T-0806 flai stats reports strategic usage as its own figures, the mean cost per agent hour, and each item's expected cost
- T-0807 flai show and the item page show the expected cost and the strategic usage
- T-0808 The cost charts draw strategic cost as its own series

## Notes
