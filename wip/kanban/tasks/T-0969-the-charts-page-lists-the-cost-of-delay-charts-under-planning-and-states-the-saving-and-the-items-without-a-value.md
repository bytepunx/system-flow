---
id: T-0969
type: task
nature: feature
title: The charts page lists the cost of delay charts under Planning and states the saving and the items without a value
status: done
parent: S-0213
owner: alex
created: 2026-10-05T05:47:24Z
updated: 2026-10-07T08:04:27Z
transitions:
  - to: ready
    at: 2026-10-07T08:00:07Z
    by: agent-S-0213
  - to: in-progress
    at: 2026-10-07T08:00:07Z
    by: agent-S-0213
  - to: done
    at: 2026-10-07T08:04:27Z
    by: agent-S-0213
stream: S-0213
tags: [dashboard]
touches: ["flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-0967]
usage:
  source: log
  seconds: 260
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 48
      output: 287
      cache_read: 2389086
      cache_write: 86488
      cost: 1.1106
---
# T-0969 The charts page lists the cost of delay charts under Planning and states the saving and the items without a value

## Work

In `flaiover/src/routes/charts/[kind]/+page.svelte`, the Charts menu's groups are the rows `flow` and `usage`. Add a `planning` row with the three cost of delay kinds, or add them to the row S-0212 made if it has landed. On these charts:

- state the count of items without a value under the chart, per column, so their absence is visible, and say none when there are none
- on `cod-order`, state the saving: what the current pull order costs against the cheaper of cost of delay and WSJF, in the project's currency, and which order that is; name the ready stories left out
- when the host's flai sends no `cost_of_delay`, say it is older than the dashboard, as the spend charts do
- give the table view under each chart its rows: per day, per week, or per projected pull

The site menu in `sitemenu.ts` keeps Charts as one page under Status. No change there unless S-0212 made one.

It waits for T-0967, whose kinds and builders it draws.

## Done when

- `charts.svelte.test.ts` shows the planning row with the three kinds, the stated saving, the counts without a value, and the older-flai message
- `npm run test`, `npm run check`, and `npm run lint` in `flaiover` pass

## Notes
