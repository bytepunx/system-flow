---
id: T-1183
type: task
nature: improvement
title: Charts and the chart tables show their times and dates in local time, over flai's UTC buckets
status: done
parent: S-0329
owner: alex
created: 2026-10-07T19:49:43Z
updated: 2026-10-07T20:37:55Z
transitions:
  - to: ready
    at: 2026-10-07T20:26:27Z
    by: agent-S-0329
  - to: in-progress
    at: 2026-10-07T20:26:27Z
    by: agent-S-0329
  - to: done
    at: 2026-10-07T20:37:55Z
    by: agent-S-0329
stream: S-0329
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", flaiover/src/lib/components/SpendTable.svelte, flaiover/src/lib/components/SpendTable.svelte.test.ts, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/WaitTable.svelte.test.ts, flaiover/src/lib/components/ForecastTable.svelte, flaiover/src/lib/components/ForecastTable.svelte.test.ts]
after: [T-1180]
usage:
  source: log
  seconds: 688
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 165
      output: 60144
      cache_read: 10235408
      cache_write: 263175
      cost: 4.8636
---
# T-1183 Charts and the chart tables show their times and dates in local time, over flai's UTC buckets

## Work

The charts draw their time axes in UTC and label hours `... UTC`; the tables beside them cut dates with `.slice(0, 10)`.

- `flaiover/src/lib/viz/charts.ts`:
  - Drop `useUTC: true` from the five options that set it (Time in State, Cost over Time, Hold Time, Spend by Model, Strategic Use), so axis ticks and tooltips read local time; check the other time axes (Cycle Time, Burn-up) do the same.
  - `bucketLabel` labels an hour in local time through `localTime`, and a day or a week by its date, without `UTC`.
  - The day and week buckets stay flai's: `flai stats` counts them per UTC day and per UTC week from Monday (`design/system/metrics.md`, ADR-0056), and this story does not change that contract. Draw a day's or a week's bar at the local midnight of its date, not at UTC midnight, so the bar sits on its day's tick in any zone.
- `flaiover/src/routes/charts/[kind]/+page.svelte`: the tables' `completed.slice(0, 10)` go through `localDate`; the hourly rows use `bucketLabel`; the notice `holds the items that entered done in it, in UTC` keeps saying the buckets are UTC days, since they are.
- `flaiover/src/lib/components/SpendTable.svelte`, `WaitTable.svelte`, and `ForecastTable.svelte`: their times and dates through `localTime` and `localDate`. `WaitTable.svelte`'s comment says it shows UTC as flai records it; it no longer does.
- Update `charts.test.ts` (`useUTC` expectations, `bucketLabel`, positions of day bars) and the table and page tests to the tests' pinned zone.

It waits for T-1180 because it calls its formatter. It shares no file with the other two view tasks, so the three run together.

## Done when

- No chart axis, tooltip, or table on the charts pages shows a time in UTC, and an hour bucket's label is local.
- Day and week bars sit on their own date's tick in the tests' pinned zone.
- `flai test` on the files above passes.

## Notes

- Rebucketing by the viewer's local day would change `flai stats` and the metrics contract, which needs an ADR; it is out of this story unless the operator asks for it on the plan thread.
