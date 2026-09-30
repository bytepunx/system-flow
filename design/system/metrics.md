---
title: Flow metrics
updated: 2026-09-30
status: active
topics: [cli, dashboard]
---

# Flow metrics

Everything the dashboard charts is derived from work item front matter. This document defines each derived value so `flai stats` and `flaiover` compute the same numbers.

## Derived timestamps per item

| Value | Definition |
|-------|------------|
| `created` | Front matter `created` |
| `committed` | `at` of the first transition to `ready` |
| `started` | `at` of the first transition to `in-progress` |
| `completed` | `at` of the last transition, while the item is `done` or `cancelled`; none while it is open, so an item moved back out of `cancelled` is not completed until it closes again ([ADR-0055](../adrs/0055-a-story-moves-back-one-column-from-ready-in-progress-review-or-cancelled-and.md)) |
| `state_intervals` | For each transition, the interval from its `at` to the next transition's `at` (or now), labelled with the state entered |
| `blocked_total` | Sum of `blocked` intervals, open intervals end at now |

## Per item durations

| Metric | Formula | Answers |
|--------|---------|---------|
| Lead time | `completed - created` | How long from idea to delivered |
| Cycle time | `completed - started` | How long once someone started |
| Queue time | `started - committed` | How long ready work waits |
| Time in state | Sum of `state_intervals` per state | Where the process spends time |
| Blocked time | `blocked_total` | How much of cycle time was waiting |
| Estimate error | `(cycle time - estimate) / estimate` | How good estimates are, when present |

Cancelled items are excluded from lead and cycle time aggregates but included in the cancellation rate.

## Aggregates

Computed over a window (default 30 days, by `completed`) and groupable by `type`, `nature`, and `parent`.

- **Throughput**: items completed per week, one bucket per ISO week of the window, a week with nothing done among them with zero (S-0166).
- **Cycle time distribution**: median, 85th percentile, and max. The 85th percentile is the number used for forecasting.
- **WIP**: count of items in `in-progress` plus `review` at each point in time, reconstructed from transitions.
- **Flow efficiency**: `(cycle time - blocked time) / cycle time`, averaged.
- **Cancellation rate**: cancelled over cancelled plus done.
- **Time-in-state share**: total time per state as a share of total lead time, the chart that shows which part of the process to optimise.

## Usage (S-0143, S-0163)

What agents spent on items, from each item's `usage` front matter ([work-hierarchy.md](work-hierarchy.md), [ADR-0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)), and how it is laid out over time ([ADR-0053](../adrs/0053-usage-is-charted-as-spend-over-time-in-buckets-for-every-item-type-at-once-and.md)). An item without `usage`, or with an empty one, has no usage values and is left out of every usage aggregate.

| Value | Definition |
|-------|------------|
| Tokens | Sum over the item's models of `input + output + cache_read + cache_write` |
| Cost | Sum over the item's models of `cost`, in US dollars |
| Agent minutes | `seconds / 60` |
| Token rate | Tokens over agent minutes; absent when `seconds` is 0 |
| Per model | The same for each model the item lists: its tokens, its cost, and its tokens over the item's agent minutes |
| Estimated | The item's `estimated` |

Aggregates cover the items of the report's type that entered `done` in the window and carry usage; cancelled items are left out:

- **Totals**: items, tokens, cost, agent seconds, and whether any is estimated.
- **Per model**: for each model, the items it worked on, its tokens and cost over them, and its tokens over the agent minutes of those items.
- **Completion against time and cost**: the items in order of `completed` (then ID), each point carrying `completed`, the ID, and the cumulative count of items, tokens, and cost up to and including it. Per model, the same over the items that model worked on, counting that model's tokens and cost.

### Spend over time (S-0163)

Spend is laid out in buckets, for epics, for stories, and for tasks, whatever the report's type. Each type is summed over its own items: a story's usage holds its tasks' and is never added to them.

| Term | Definition |
|------|------------|
| Bucket | An hour, a day, or an ISO week starting Monday, in UTC, named by its start. A day unless asked otherwise. An hour needs a window of 31 days or less |
| An item's bucket | The one that holds the moment it entered `done`. Its whole usage counts there |
| Series | For a type, one point per bucket from the bucket of the first item of that type done in the window with usage, to the bucket that holds now. A bucket with no such item is a point of zeros. A type with no such item has no points |

A set of items, whether those of a bucket, of a type over the window, or of either that one model worked on, has these values:

| Value | Definition |
|-------|------------|
| Items, tokens, cost, seconds | Their count and the sums of their tokens, cost, and agent seconds. For a model: the items it worked on, its own tokens and cost, and those items' agent seconds |
| Estimated | Whether any of them is |
| Tokens per item | Tokens over items: the mean an item took |
| Cost per item | Cost over items |
| Tokens per minute | Tokens over agent minutes |
| Tokens per dollar | Tokens over cost |
| Mean tokens, mean cost | Of a bucket only: the running mean per bucket, the sum from the first bucket of the series to this one over the number of those buckets |

A value is absent when its divisor is zero.

`flai stats --json` carries each item's usage under `items[].usage` and the aggregates under `usage`: `items`, `tokens`, `cost`, `seconds`, `estimated`, `models`, `done`, `by_model`, and since S-0163 `bucket` (`hour`, `day`, or `week`) and `spend`. `spend` has the keys `epic`, `story`, and `task`, each with the values of its items over the window (`items`, `tokens`, `cost`, `seconds`, `estimated`, `tokens_per_item`, `cost_per_item`, `tokens_per_minute`, `tokens_per_dollar`), the same per model under `models`, and the series under `buckets`: each point has `at`, the same values, `mean_tokens`, `mean_cost`, and its `models`. A rate is `tokens_per_minute`. `tokens_per_hour`, sixty times that, stays beside it on items and models for what was written to flai 1.25.

## Charts

Every chart spans the window chosen ([ADR-0054](../adrs/0054-every-chart-spans-the-window-chosen-its-time-axis-runs-from-the-window-s-start.md), S-0166). A time axis runs from the window's start to the report's now, whatever the data: a series by the day from the day that holds the start, a series in buckets from the bucket that holds the start to the one that holds now, with half a bucket either side. A chart per item plots only the items completed in the window, and time in state groups them by the day they were completed; `items` in `flai stats --json` holds every item of the type, and the dashboard picks them.

| Chart | Data | Notes |
|-------|------|-------|
| Kanban board | Current status of active items grouped by column, with age in column and blocked flag | Age in column is now minus last transition |
| Cycle time scatter | One point per story completed in the window, x completed date, y cycle time, with 50th and 85th percentile lines | Filter by nature |
| Burn-up | Per epic or whole repo, cumulative stories created versus done, per day of the window | Scope line and done line, forecast line from throughput |
| Cumulative flow diagram | Stacked count of items per state per day of the window | Widening bands show where work piles up |
| Time in state | Stacked bar per day of the window with stories completed: the mean hours per state of those stories ([ADR-0056](../adrs/0056-time-in-state-is-one-stacked-bar-per-day-of-the-window-the-mean-hours-per-state.md)), and the aggregate share | The process optimisation chart. The table lists each story |
| Throughput | Bar per week of the window | With nature breakdown |
| Aging WIP | Active items by age since started, against the 85th percentile | Flags items likely to be late |
| Estimate vs actual | Scatter, only items completed in the window with `estimate` | |
| Token rate | Per model, one point per bucket with agent time: the model's tokens per agent minute over the items done in the bucket | Of the report's type. Replaces the point per item in tokens per agent hour (S-0163) |
| Tokens per bucket | Bar per bucket, the tokens of the items done in it, stacked by model, with the running mean per bucket as a line | Of the report's type. Titled by the bucket: tokens per day |
| Tokens per item | One point per bucket with items: tokens per item, one series per type | Or one series per model, for the report's type |
| Tokens per dollar | Per model, one point per bucket with cost: tokens over cost | Of the report's type |
| Cost per bucket | Bar per bucket, the cost of the items done in it, stacked by model, with the running mean per bucket as a line | Of the report's type. Estimated costs marked |
| Cost per item | One point per bucket with items: cost per item, one series per type | Or one series per model, for the report's type |
| Cost by item | Bar per item completed in the window with usage, stacked by model | Estimated costs marked |
| Completion against time | Per model, cumulative items done over time | The per model `done` series |
| Completion against cost | Per model, cumulative items done over cumulative cost | The same series, x cost |

## Precision rules

So that `flai stats` and the dashboard agree to the second:

- Percentiles use nearest rank on the sorted values: the p-th percentile is the value at position `ceil(p/100 × n)`.
- The backlog interval starts at `created`. An open item's current interval ends at now. Closed items accrue nothing after `completed`.
- Queue time exists only when both `committed` and `started` exist and `started` is not before `committed`.
- Cycle time exists only when the item was started. An item cancelled from backlog has a lead time but no cycle time and is excluded from cycle aggregates.
- Burn-up counts an item in scope from its `created` day and excludes cancelled items. Cumulative flow uses the item's state at the end of each day (23:59:59 UTC). Both run from the day that holds the window's start, or the first item's `created` day if that is later, to today (S-0166).
- Throughput per week is completed items in the window divided by window days over seven. Weekly buckets are ISO weeks starting Monday, from the one that holds the window's start to the one that holds now.
- Flow efficiency averages `(cycle - blocked) / cycle` over completed items with a positive cycle time.
- Time-in-state share divides total seconds per state by total lead time, over completed items in the window.
- A bucket holds the moments from its start up to, not including, the next one's. A week's bucket starts on the Monday of the ISO week, at 00:00:00 UTC. The running mean divides by the number of buckets from the first of the series, empty ones counted.

## Data access

The dashboard reads items via its server API, which scans `wip/kanban` and `wip/archive` on request and caches by file modification time. `flai stats` prints the same aggregates as a table or JSON so scripts and CI can use them. Both share the same definitions above; the Go implementation is the reference and `flaiover` has a fixture test that checks its numbers against `flai stats --json` output on the sample repo in the template.
