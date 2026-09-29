---
id: T-0542
type: task
nature: feature
title: flai stats reports tokens and cost by model
status: done
parent: S-0143
owner: alex
created: 2026-09-29T06:08:16Z
updated: 2026-09-29T06:23:25Z
transitions:
  - to: ready
    at: 2026-09-29T06:20:48Z
    by: agent-S-0143
  - to: in-progress
    at: 2026-09-29T06:20:48Z
    by: agent-S-0143
  - to: done
    at: 2026-09-29T06:23:25Z
    by: agent-S-0143
stream: S-0143
tags: []
touches: [flai/internal/metrics, flai/cmd/stats.go, design/system/metrics.md, design/adrs]
---
# T-0542 flai stats reports tokens and cost by model

## Work

- An ADR: work items record the tokens and cost their agents spent, measured from the agents' logs, summed up the hierarchy at done, and charted by model.
- `design/system/metrics.md` defines per item tokens, cost, and token rate per hour of agent work, and the aggregates: per model, per type; cumulative items done against time and against cost.
- `flai stats --json` carries them per item and per model; the table output adds a usage line.

## Done when

- Metric tests over fixture items with usage pass; `flai stats` on this repository shows usage.
- `make test` and lint pass.

## Notes
