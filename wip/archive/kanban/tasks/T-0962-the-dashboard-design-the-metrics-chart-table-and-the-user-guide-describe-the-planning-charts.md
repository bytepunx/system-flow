---
id: T-0962
type: task
nature: feature
title: The dashboard design, the metrics chart table, and the user guide describe the planning charts
status: done
parent: S-0212
owner: alex
created: 2026-10-05T05:46:50Z
updated: 2026-10-07T07:17:35Z
transitions:
  - to: ready
    at: 2026-10-07T07:09:28Z
    by: agent-S-0212
  - to: in-progress
    at: 2026-10-07T07:09:29Z
    by: agent-S-0212
  - to: done
    at: 2026-10-07T07:17:18Z
    by: agent-S-0212
stream: S-0212
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, design/system/metrics.md, docs/users/flaiover.md]
after: [T-0956]
usage:
  source: log
  seconds: 469
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 26
      output: 14190
      cache_read: 1526989
      cache_write: 49860
      cost: 0.8844
---
# T-0962 The dashboard design, the metrics chart table, and the user guide describe the planning charts

## Work

Waits for T-0956: the three charts' mapping is settled by then, so the docs describe what was built. It touches only docs, so it runs together with T-0958.

- In `design/system/flaiover-dashboard.md`, add rows for `/charts/forecast-accuracy`, `/charts/delivery-accuracy`, and `/charts/forecast-by-model` to the views table. Say what each plots and from which `flai stats --json` fields (`items[].forecast_error_seconds`, `delivery_error_seconds`, `estimate_error_seconds`, and `forecasts`). Change the chart groups paragraph from two groups to three: flow, usage, and planning. Name the nature and model filters, and `ForecastTable.svelte`.
- In `design/system/metrics.md`, add the three charts to the Charts table, each with its data and notes. This maps metrics S-0205 already defines, so it needs no ADR.
- In `docs/users/flaiover.md` § Charts, add the Planning group and its three charts to the table, and say what a positive error means and how to read the on-time share. Say that a story needs a forecast, from the planner or `flai edit`, to appear.

## Done when

- `scripts/lint-md.sh` passes on the three documents
- Each document names the three charts as the dashboard draws them, and `flai check --strict` reports nothing new

## Notes
