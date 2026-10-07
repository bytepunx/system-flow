---
id: T-0956
type: task
nature: feature
title: The forecast-by-model chart maps the p50 absolute forecast error per bucket per model
status: done
parent: S-0212
owner: alex
created: 2026-10-05T05:46:25Z
updated: 2026-10-07T07:09:27Z
transitions:
  - to: ready
    at: 2026-10-07T07:02:57Z
    by: agent-S-0212
  - to: in-progress
    at: 2026-10-07T07:02:57Z
    by: agent-S-0212
  - to: done
    at: 2026-10-07T07:09:27Z
    by: agent-S-0212
stream: S-0212
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
after: [T-0951]
usage:
  source: log
  seconds: 390
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 24202
      cache_read: 2604455
      cache_write: 85043
      cost: 1.5085
---
# T-0956 The forecast-by-model chart maps the p50 absolute forecast error per bucket per model

## Work

Waits for T-0951: it changes the same two files, so the chart builders are written one after another.

- In `flaiover/src/lib/viz/charts.ts`, build `forecast-by-model`. Draw one line per model (`agent.model`, `(none)` without one), and one point per bucket in which that model's stories with a forecast were completed. The point is the p50 of their absolute `forecast_error_seconds`, worked out as flai works out its percentiles.
- Use the bucket control (hour, day, or week) and its UTC buckets. Run the axis from the bucket that holds the window's start to the one that holds now, with half a bucket either side (ADR-0054, as the spend charts do). Give each model its fixed slot and mark (`modelSlot`, `MODEL_SYMBOL` in `palette.ts`).
- Offer the nature filter, but not the model filter, since the chart shows every model.
- Add a function giving the rows of the planning charts' table view: story, completed, nature, model, forecast, actual, forecast error, delivery error, and estimate error. T-0958 draws it.
- Cover the mapping and the rows in `charts.test.ts`: buckets, per-model p50s, an empty bucket, and the window span.

## Done when

- `scripts/flaiover-test.sh` passes, with tests that pin the forecast-by-model mapping and the planning table rows
- Over the whole window and one model, the bucket p50s agree with `forecasts.forecast.by_model` when the window is a single bucket

## Notes
