---
id: T-0950
type: task
nature: feature
title: The chart data types carry the forecast errors, and the forecast-accuracy chart maps them with nature and model filters
status: done
parent: S-0212
owner: alex
created: 2026-10-05T05:46:04Z
updated: 2026-10-07T06:56:28Z
transitions:
  - to: ready
    at: 2026-10-07T06:49:20Z
    by: agent-S-0212
  - to: in-progress
    at: 2026-10-07T06:49:20Z
    by: agent-S-0212
  - to: done
    at: 2026-10-07T06:56:28Z
    by: agent-S-0212
stream: S-0212
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
usage:
  source: log
  seconds: 428
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 27794
      cache_read: 2991012
      cache_write: 97665
      cost: 1.7324
---
# T-0950 The chart data types carry the forecast errors, and the forecast-accuracy chart maps them with nature and model filters

## Work

First layer: it waits for nothing, and the other chart tasks build on the types and controls it adds.

- In `flaiover/src/lib/viz/charts.ts`, add to `ItemMetrics` the per-item fields `flai stats --json` already sends since S-0205: `forecast_seconds`, `forecast_error_seconds`, `delivery_error_seconds`, and `estimate_error_seconds`. Add `forecasts` to `Report`: `forecast`, `delivery`, and `estimate`, each `count`, `p50_seconds`, `p85_seconds`, `by_nature`, and `by_model`, as `design/system/metrics.md` § Forecasts and estimates defines. Let `normalise()` accept a report from an older flai without them.
- Add a `PLANNING_KINDS` list with `forecast-accuracy`, `delivery-accuracy`, and `forecast-by-model`, and add it to `KINDS`.
- Give `controls()` a nature and a model filter for the planning charts. The model is the story agent's `agent.model`, `(none)` without one, as flai groups it.
- Build `forecast-accuracy`: one point per story completed in the window (`completedIn`), x completed, y `forecast_error_seconds`. Add a second series of `estimate_error_seconds` for the stories that have one. Draw p50 and p85 lines of absolute error. Take them from `forecasts.forecast`, or from its `by_nature` or `by_model` entry under one filter, so they match `flai stats --json` to the second. With both filters set, work them out from the points shown, the same way flai does.
- Cover the mapping in `charts.test.ts` on a fixture report: points, both series, the lines under each filter, and the window span of ADR-0054.

## Done when

- `scripts/flaiover-test.sh` passes, with tests that pin the forecast-accuracy mapping and its filters to the second
- A report without `forecasts` or the error fields maps to an empty chart, not an error

## Notes
