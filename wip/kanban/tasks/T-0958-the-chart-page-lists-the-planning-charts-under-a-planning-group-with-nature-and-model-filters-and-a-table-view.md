---
id: T-0958
type: task
nature: feature
title: The chart page lists the planning charts under a Planning group, with nature and model filters and a table view
status: backlog
parent: S-0212
owner: alex
created: 2026-10-05T05:46:37Z
updated: 2026-10-05T05:46:37Z
transitions: []
stream: S-0212
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/components/ForecastTable.svelte, flaiover/src/lib/components/ForecastTable.svelte.test.ts]
after: [T-0956]
---
# T-0958 The chart page lists the planning charts under a Planning group, with nature and model filters and a table view

## Work

Waits for T-0956: the page draws the three chart kinds and the table rows the chart tasks add, and needs all of them in place.

- In `flaiover/src/routes/charts/[kind]/+page.svelte`, list `PLANNING_KINDS` as a third group, Planning, after flow and usage in the page header. Give each chart a title.
- Show the nature and model selects when `controls()` asks for them, as the type, epic, and bucket selects are shown. Keep the window in `chartWindow`, and ask `/api/stats` for the window and bucket as the other charts do.
- Add `flaiover/src/lib/components/ForecastTable.svelte`, the table view under the planning charts, after `SpendTable.svelte`, with the rows T-0956 gives. When no story in the window has a forecast, say so, and say how one is set (the planner, or `flai edit --forecast-duration`).
- Add a summary strip for the planning charts: stories with a forecast, and p50 and p85 of absolute forecast and delivery error, from `forecasts`.
- Cover the page in `charts.svelte.test.ts` and the table in `ForecastTable.svelte.test.ts`: the Planning group, the filters, the empty state, and the rows.

## Done when

- `scripts/flaiover-test.sh` passes, with lint, type check, and the new page and table tests
- The Planning group's three charts open from the chart page's header and draw against a fixture report

## Notes
