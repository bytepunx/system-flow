---
id: TH-0148
title: "S-0213 plan: seven tasks in five layers, metrics contract first"
anchor:
  path: wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md
  item: S-0213
status: resolved
participants: [planner-S-0213, alex]
created: 2026-10-05T05:48:02Z
updated: 2026-10-05T05:50:56Z
---

# TH-0148 S-0213 plan: seven tasks in five layers, metrics contract first

On wip/kanban/stories/S-0213-charts-show-cost-of-delay-outstanding-incurred-and-what-the-pull-order-costs.md.

## Entries

### 2026-10-05T05:48:02Z planner-S-0213
Plan for S-0213 (planner-S-0213). S-0205 already gives `flai stats --json` the cost of delay outstanding per column per day and incurred per day and week. Two things are missing: a count of items without a value, and any projection of the pull order's cost. Both are additions to `metrics.md`, so they need an ADR first.

Tasks, by layer:

- Layer 1: T-0953, which defines `cost_of_delay.days[].without_value` and `cost_of_delay.order` in `metrics.md`, with an ADR.
- Layer 2: T-0955, in which `flai stats` counts the items without a value per column per day.
- Layer 3: T-0959, in which `flai stats` projects the ready column's cost under the pull order, by cost of delay, and by WSJF. It waits for T-0955 because both change the types in `costofdelay.go`.
- Layer 4: T-0965 and T-0967, which run together.
  - T-0965: `flai stats` prints the saving and the count without a value, and `docs/users/flai.md` describes them.
  - T-0967: `charts.ts` gets the types, the three kinds, the builders, and the mapping tests.
- Layer 5: T-0969 and T-0970, which run together.
  - T-0969: the charts page gets a planning row and states the saving and the counts.
  - T-0970: `flaiover-dashboard.md` and `docs/users/flaiover.md` describe the charts.

Assumptions:

- The projection counts cost until each ready story is pulled, not until it is delivered. This matches `cost_of_delay_incurred`, which counts only time in backlog and ready. The lanes are the in-progress limit, laid out as `flai forecast` lays them (`planning.playOut`).
- The cod and wsjf orders reuse S-0217's policy function (T-0810, `workitem/policy.go`). That is why the story waits for S-0217.
- `cod-order` draws three lines, current, cod, and wsjf. The criterion says two lines and puts WSJF in parentheses. I read that as WSJF being a third order on the same chart, and the saving is stated against the cheaper of the two. If you meant two lines, with WSJF only stated, T-0967 and T-0969 drop one series.
- `cod-order`'s axis runs from now to the projected horizon, not over the window. It is a projection, and ADR-0054's window covers the history charts.
- "The Charts menu lists them under Planning" means a planning row on the charts page, next to flow and usage. If S-0212 lands first, they join its row. `sitemenu.ts` stays as declared but is not expected to change.
- Ready stories with no value or no forecast are left out of the projection and listed by ID. Tasks are not counted in `without_value`.

Figures:

- Forecast: kept at 1h30m rather than `flai forecast`'s 31m. This is S-0205's shape of work.
- Cost of delay: 91.22 USD a week from `flai cod`, the story's share of E-0016 by forecast. It replaces 73.17.
- Touches: added `flai/internal/statsread`, `flai/internal/planning/forecast.go`, `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`, and `docs/users/flai.md`. Every declared touch is kept.

Nothing to split, merge, or drop.

### 2026-10-05T05:50:56Z alex
Resolved.
