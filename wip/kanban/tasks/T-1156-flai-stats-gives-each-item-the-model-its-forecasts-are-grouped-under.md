---
id: T-1156
type: task
nature: feature
title: flai stats gives each item the model its forecasts are grouped under
status: done
parent: S-0212
owner: alex
created: 2026-10-07T06:47:57Z
updated: 2026-10-07T06:48:55Z
transitions:
  - to: ready
    at: 2026-10-07T06:48:06Z
    by: agent-S-0212
  - to: in-progress
    at: 2026-10-07T06:48:06Z
    by: agent-S-0212
  - to: done
    at: 2026-10-07T06:48:55Z
    by: agent-S-0212
stream: S-0212
tags: []
touches: [flai/internal/metrics/metrics.go, flai/internal/metrics/forecast.go, flai/internal/metrics/forecast_test.go, design/system/metrics.md, design/adrs/0111-flai-stats-gives-each-item-the-model-its-forecasts-are-grouped-under.md, design/adrs/README.md]
usage:
  source: log
  seconds: 49
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 19
      output: 10185
      cache_read: 1096030
      cache_write: 35788
      cost: 0.6348
---

# T-1156 flai stats gives each item the model its forecasts are grouped under

## Work

First layer: it waits for nothing. The planning charts filter by the story agent's model and work out per-model percentiles per bucket, and `flai stats --json` did not say which model an item is grouped under ([ADR-0111](../../../design/adrs/0111-flai-stats-gives-each-item-the-model-its-forecasts-are-grouped-under.md)).

- In `flai/internal/metrics`, give `ItemMetrics` a `model` field: the item's `agent.model`, `(none)` without one, and have `forecasts` group by it.
- Add the field to the per-item table in `design/system/metrics.md` § Forecasts and estimates, linking the ADR.

## Done when

- `forecast_test.go` pins `model` on an item with a model, one whose agent names none, and one with no agent
- `flai test` passes on the changed paths

## Notes
