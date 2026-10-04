---
id: T-0778
type: task
nature: feature
title: flai stats reports forecast, delivery, and estimate error per story and their p50 and p85
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:40:03Z
updated: 2026-10-03T20:47:24Z
transitions:
  - to: ready
    at: 2026-10-03T20:41:51Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T20:41:51Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T20:47:24Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [flai/internal/metrics]
after: [T-0777]
usage:
  source: log
  seconds: 333
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 10664
      cache_read: 1539827
      cache_write: 41684
      cost: 0.7711
---
# T-0778 flai stats reports forecast, delivery, and estimate error per story and their p50 and p85

## Work

In `flai/internal/metrics`, per item `forecast_seconds`, `forecast_error_seconds`, `delivery_error_seconds`, and `estimate_error_seconds`, and the `forecasts` aggregates (p50 and p85 of absolute error, overall, by nature, by agent model), as `metrics.md` § Forecasts and estimates defines them. Waits for the spec task, which defines them; touches `metrics.go`, so the later metric tasks wait for it.

## Done when

- [ ] `Compute` reports each value, absent when its input is
- [ ] Tests pin each value and aggregate on a fixture of their own, to the second

## Notes
