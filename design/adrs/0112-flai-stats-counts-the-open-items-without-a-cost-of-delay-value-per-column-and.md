---
id: ADR-0112
title: "flai stats counts the open items without a cost of delay value per column, and projects the ready column's cost of delay until each story is pulled under the pull order, by cost of delay, and by WSJF"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0081]
topics: [cli, dashboard, planning]
---

# ADR-0112 flai stats counts the open items without a cost of delay value per column, and projects the ready column's cost of delay until each story is pulled under the pull order, by cost of delay, and by WSJF

## Context

ADR-0081 has `flai stats` report cost of delay as sums: the value outstanding per column at the end of each day, and what waiting incurred per day and per week. S-0213 charts them, and adds a chart of what the pull order costs against ordering by cost of delay and by WSJF. Two things are missing:

- An item without a `cost_of_delay.value` adds nothing to a sum, so it is invisible. A column of ten stories, one valued, looks cheap.
- Nothing says what the ready column will cost under its pull order, or what S-0217's `cod` and `wsjf` policy orders would cost instead. `flai order --by` ranks the stories; it does not price the ranking.

`design/system/metrics.md` is the contract between `flai stats` and the dashboard, so both are defined there first.

## Decision

`flai stats` adds two values under `cost_of_delay`, defined in `metrics.md` § Cost of delay and § What the pull order costs.

- `days[].without_value`: per column, the number of the items of the report's type in it at the end of the day with no value, every column present. They are the items `outstanding` sums over. A task carries no value and is never counted.
- `order`: the ready column's cost of delay, projected from now until each story is pulled, under three orders: `current`, the board's pull order; `cod`; and `wsjf`. Each order is a series of points of cumulative cost against the projected pull, with its total. Beside them are the `saving` of the cheaper of `cod` and `wsjf` over `current`, which order that is, and the ready stories `left_out`.
- A story's projected cost runs until it is pulled, not until it is delivered. An item waits only while it is in `backlog` or `ready`, as `cost_of_delay_incurred` counts it, so the projection and the history add up the same way.
- The `cod` and `wsjf` orders are S-0217's, from `workitem.OrderByPolicy`, applied to the placed stories in pull order. There is no second copy of the policies, and no hand placement is kept.
- Each order is played out on lanes as `flai forecast` plays out the pull order: the board's `in-progress` limit, lanes freed by the stories in progress, each story holding its lane for its duration times its busy factor, and `after` respected. With no limit every story is pulled at once, as `flai forecast` pulls them, and every order costs nothing.
- Only the ready stories with both a value and a `forecast.duration` are placed, each with its own duration, so the three orders place the same stories and their totals compare.

## Consequences

- The dashboard draws the cost of delay charts from `flai stats --json` alone, with no arithmetic of its own, and `flai stats` prints the same figures.
- The JSON grows by a map per day and the `order` object. A dashboard reading an older flai gets neither and must show the charts without them.
- `flai stats` now reads the board's order and limit and the planning settings, and replays the forecast's lanes, so `metrics` depends on the `planning` play-out.
- The projection is only as good as each story's own forecast: a stale `forecast.duration` misplaces every story after it.
- Stories without a value or a forecast are left out of every order and named in `left_out`, rather than priced with a guess.
- `cod-order` is the one chart whose time axis is not the window (ADR-0054): it runs from now to the last pull.

## Alternatives considered

- **Project waiting until delivery.** Charges the time in progress and review as waiting, which `cost_of_delay_incurred` does not, so the projection would not match the history.
- **Place stories without a forecast with the forecaster's estimate.** WSJF would then rank on a duration nobody wrote down, and a story without a value would still have to be left out; leaving both out keeps the three totals over one set of stories.
- **One lane, Smith's rule, ignoring the in-progress limit.** Simpler and optimal for WSJF on one machine, but the board pulls several stories at once; the projection would disagree with `flai forecast` and overstate every order's cost.
- **Compute the projection in the dashboard.** The dashboard would need the board, the limit, the forecaster's history, and the policies, and would repeat what flai computes; the designer chose on 2026-10-02 that this arithmetic lives in flai.
- **Count the items without a value only today.** The per-day count costs one pass more and lets a chart show when values were set; the charts show today's.
