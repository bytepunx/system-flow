---
id: T-0951
type: task
nature: feature
title: The delivery-accuracy chart maps each story's delivery error in days and the weekly share delivered on time
status: backlog
parent: S-0212
owner: alex
created: 2026-10-05T05:46:13Z
updated: 2026-10-05T05:46:13Z
transitions: []
stream: S-0212
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
after: [T-0950]
---
# T-0951 The delivery-accuracy chart maps each story's delivery error in days and the weekly share delivered on time

## Work

Waits for T-0950: it builds on the types, `PLANNING_KINDS`, and the filters it adds, and it changes the same two files.

- In `flaiover/src/lib/viz/charts.ts`, build `delivery-accuracy`. Plot one point per story completed in the window, x completed, y `delivery_error_seconds` in days. Positive is later than the forecast date, as `design/system/metrics.md` defines.
- Add the share of those stories delivered on or before the forecast date (`delivery_error_seconds` at most 0) per ISO week, on a second axis. Run the weeks from the week that holds the window's start to the week that holds now, as ADR-0054 requires. A week with no forecast stories has no share, not 0.
- Apply the nature and model filters of T-0950. Show the p50 and p85 of absolute delivery error from `forecasts.delivery` the same way, so the figures match `flai stats --json`.
- Cover the mapping in `charts.test.ts`: the days, the weekly share and its empty weeks, and the window span.

## Done when

- `scripts/flaiover-test.sh` passes, with tests that pin the delivery-accuracy mapping, its weekly share, and its window span
- A story without `forecast.delivery` is left out of the points and the share

## Notes
