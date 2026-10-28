---
id: T-2256
type: task
nature: remediation
title: "metrics.md § Charts says what each chart draws: Cycle Time per epic coloured by nature, Burn-up scope and done, the dashed all-models line on every chart per model, and rows for Agent Waiting, Strategic Cost, and Strategic Use"
status: backlog
parent: S-0527
owner: alex
created: 2026-10-28T15:00:51Z
updated: 2026-10-28T15:01:36Z
transitions: []
stream: S-0527
tags: [design, metrics]
touches: [design/system/metrics.md]
after: [T-2255]
---
# T-2256 metrics.md § Charts says what each chart draws: Cycle Time per epic coloured by nature, Burn-up scope and done, the dashed all-models line on every chart per model, and rows for Agent Waiting, Strategic Cost, and Strategic Use

## Work

Rewrite the rows of the § Charts table in `design/system/metrics.md` to what the dashboard draws, citing the ADR T-2255 writes, and bump `updated`. It waits for T-2255 because it cites that ADR's number, and metrics.md changes only with an ADR (ADR-0054).

- Cycle Time: Data, one point per item completed in the window, x completed date, y cycle time, with 50th and 85th percentile lines, coloured by nature; Notes, chosen per epic or whole, no nature filter.
- Burn-up: Notes, the scope line and the done line; drop the forecast line from throughput.
- Tokens / Min and Tokens / $: Notes, with two models or more a dashed line for all of them, citing ADR-0057 and the new ADR, as the Avg. Time / Model and Avg. Cost / Model rows do.
- Add a row for Agent Waiting: a stacked bar per ISO week of the window of `waiting.weeks[]`, threads and review in hours, review split into checks and reviewers under pull-request acceptance, with the mean wait per item as a dashed line, and `waiting.longest[]` as the table under it (ADR-0114, ADR-0168), under the explainers' and `docs/users/flaiover.md`'s words; a week with no item has no point on the line.
- Add a row for Strategic Cost: a bar per day of `strategic_days[]`, `agents[].cost` stacked, the line `cost_per_item`, and the dashed line of the sum of `cost` over the sum of `completed` with its share of `usage.cost` over `usage.items` (ADR-0083).
- Add a row for Strategic Use: a bar per day of `strategic_days[].agents[].seconds` in hours, the lines mean cycle time, `cycle_time_seconds`, and mean waiting, which the dashboard works out from `items[]` as the mean over the items of the report's type completed that day of `wait_threads_seconds` plus `wait_review_seconds`, a wait an item lacks counted as zero (TH-0140). Say the same under § Strategic use per day, after its table, so that a reader of `strategic_days[]` knows the waiting is not a field of it.

Keep every other row as it is. Match each row's words to the chart's explainer in `flaiover/src/lib/viz/explain.ts` and to `design/system/flaiover-dashboard.md`, so the three agree; where they do not, the chart as `flaiover/src/lib/viz/charts.ts` draws it decides, and the ADR records it. Use the titles `flaiover/src/lib/viz/catalog.ts` gives the kinds, since T-2257's test matches rows by title.

## Done when

- § Charts has a row, under its catalog title, for every one of the dashboard's chart kinds, and the Cycle Time, Burn-up, Tokens / Min, and Tokens / $ rows say what `charts.ts` draws, with no forecast line and no nature filter.
- § Strategic use per day says the mean waiting Strategic Use draws is the dashboard's sum, not a field of `strategic_days[]`.
- The markdown lint and `flai check --strict` are clean on `design/system/metrics.md`, and `flaiover/src/lib/viz/explain.test.ts` still finds every explainer's anchor.

## Notes

Drafted by planner-S-0527. S-0502 (in progress on 2026-10-28) touches `design/system/metrics.md` too, in § Waiting; sync before editing and keep its change.
