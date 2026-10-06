---
id: S-0226
type: story
nature: improvement
title: The orchestrator's cost is recorded on the story or epic each decision concerned
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-06T11:44:49Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:36Z
    by: alex
tags: [flai]
topics: [orchestration, planning]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/mcpserver, design/system/strategic-agents.md, flai/internal/metrics, design/adrs, flai/cmd/activity.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/work-hierarchy.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 145
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 970
          output: 198
          cache_read: 5667923
          cache_write: 380508
          cost: 1.1704
        - model: claude-opus-5-5
          input: 216
          output: 15545
          cache_read: 7042453
          cache_write: 316523
          cost: 4.085
cost_of_delay:
  value: 33.55
  by: planner-S-0226
  at: 2026-10-05T04:44:17Z
forecast:
  duration: 35m
  delivery: 2026-10-06T13:13:00Z
  basis: "Its own forecast of 35m; 3rd in the pull order with an in-progress limit of 3, behind S-0221, S-0222 and S-0224."
  by: flai
  at: 2026-10-06T11:44:49Z
---
# S-0226 The orchestrator's cost is recorded on the story or epic each decision concerned

## Goal

ADR-0051 records what story agents spend on the items they work. S-0225 charges a planner activity's usage to the item it planned, under a `strategic` entry per kind (ADR-0083). `flai stats` and the cost charts already show that entry apart from the agents' figures.

The orchestrator works for the whole project, not for one item. Each of its activities should be charged to the items that decision concerned, or to a project total when it concerned none, so that none of its cost is lost.

## Acceptance criteria
- [ ] An orchestrator activity's usage is charged when the activity is logged, through `activity_log` or as its run ends. The usage is apportioned to the activity's span as ADR-0083 apportions a planner's. It goes under the `orchestrator` entry of the items the activity named, split evenly between them, and is summed up the hierarchy as ADR-0083 sums it
- [ ] An activity that names no item is charged to a project strategic total. `flai stats` reports that total per kind beside the per-item figures, so that the per-kind totals equal the activity document's totals
- [ ] An ADR extends ADR-0083 to the orchestrator and the project total. `design/system/metrics.md` and `design/system/strategic-agents.md` say how the charge works. Tests pin the split, the roll-up, and the project total

## Tasks
- T-0878 An ADR extends ADR-0083 to the orchestrator's items and a project strategic total
- T-0879 An orchestrator activity's apportioned usage is charged evenly to the items it named and the items above them
- T-0880 flai stats reports the project strategic total per kind beside the per-item figures

## Notes

- Rewritten on TH-0107: the first version restated S-0225's criteria, which S-0225 delivers for every kind.

### Planning

- Touches declared, kept: `flai/internal/usage`, `flai/internal/serve`, `flai/internal/workitem`, `flai/internal/mcpserver`, `flai/internal/metrics`, `design/system/metrics.md`, `design/system/strategic-agents.md`, `design/adrs`.
- Touches added from the code layout, where S-0225 made the planner's charge: `flai/cmd/activity.go` (prints what an activity charged), `flai/cmd/stats.go` and `flai/cmd/check_stats_test.go` (the stats report and its tests).
- Touches added from the design: `design/system/work-hierarchy.md`, which ADR-0083 names as defining the `strategic` entry.
- Touches added from co-change (`flai touches suggest`) and S-0225's own touches: `design/system/flai-cli.md` (34%), `docs/users/flai.md` (33%), `docs/users/flai-reference.md` (16%, the generated `flai stats` help). The dashboard is left out: the criteria ask only for `flai stats`, and the item page already shows `strategic` per kind.
- Forecast 35m, up from flai's 16m (86 s per unit of size, size 11): S-0225, the planner's version of this story, took 52m in progress, a third of it dashboard work this story lacks, and this one adds an ADR. Delivery is flai's 13:49Z moved by the 19m added.
- Cost of delay 33.55 USD a week, as `flai cod` gives it after the forecast change: the story's share of E-0016's 1500 USD a week, 35m of 26h5m over the epic's 17 open stories without inputs. The story has no inputs of its own; it replaces planner-E-0016's 36.59, worked out from the old 45m forecast.
- Topics `orchestration` and `planning` added: the tasks reach the strategic agents' design and ADR-0083.
