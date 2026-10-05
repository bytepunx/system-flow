---
id: S-0226
type: story
nature: improvement
title: The orchestrator's cost is recorded on the story or epic each decision concerned
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-05T04:40:46Z
transitions: []
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/mcpserver, design/system/strategic-agents.md, flai/internal/metrics, design/adrs]
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
  delivery: 2026-10-05T17:46:00Z
  basis: "Its own forecast of 45m; 17th in the pull order with an in-progress limit of 3, behind S-0244, S-0258, S-0276, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0223 and S-0224."
  by: flai
  at: 2026-10-05T04:40:46Z
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

## Notes

- Rewritten on TH-0107: the first version restated S-0225's criteria, which S-0225 delivers for every kind.
