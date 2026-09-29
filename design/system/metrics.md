---
title: Flow metrics
updated: 2026-09-29
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
| `completed` | `at` of the transition to `done` or `cancelled` |
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

- **Throughput**: items completed per week.
- **Cycle time distribution**: median, 85th percentile, and max. The 85th percentile is the number used for forecasting.
- **WIP**: count of items in `in-progress` plus `review` at each point in time, reconstructed from transitions.
- **Flow efficiency**: `(cycle time - blocked time) / cycle time`, averaged.
- **Cancellation rate**: cancelled over cancelled plus done.
- **Time-in-state share**: total time per state as a share of total lead time, the chart that shows which part of the process to optimise.

## Usage (S-0143)

What agents spent on items, from each item's `usage` front matter ([work-hierarchy.md](work-hierarchy.md), [ADR-0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)). An item without `usage`, or with an empty one, has no usage values and is left out of every usage aggregate.

| Value | Definition |
|-------|------------|
| Tokens | Sum over the item's models of `input + output + cache_read + cache_write` |
| Cost | Sum over the item's models of `cost`, in US dollars |
| Agent hours | `seconds / 3600` |
| Token rate | Tokens over agent hours; absent when `seconds` is 0 |
| Per model | The same for each model the item lists: its tokens, its cost, and its tokens over the item's agent hours |
| Estimated | The item's `estimated` |

Aggregates cover the items of the report's type that entered `done` in the window and carry usage; cancelled items are left out:

- **Totals**: items, tokens, cost, agent seconds, and whether any is estimated.
- **Per model**: for each model, the items it worked on, its tokens and cost over them, and its tokens over the agent hours of those items.
- **Completion against time and cost**: the items in order of `completed` (then ID), each point carrying `completed`, the ID, and the cumulative count of items, tokens, and cost up to and including it. Per model, the same over the items that model worked on, counting that model's tokens and cost.

`flai stats --json` carries each item's usage under `items[].usage` and the aggregates under `usage` (`items`, `tokens`, `cost`, `seconds`, `estimated`, `models`, `done`, `by_model`).

## Charts

| Chart | Data | Notes |
|-------|------|-------|
| Kanban board | Current status of active items grouped by column, with age in column and blocked flag | Age in column is now minus last transition |
| Cycle time scatter | One point per completed story, x completed date, y cycle time, with 50th and 85th percentile lines | Filter by nature |
| Burn-up | Per epic or whole repo, cumulative stories created versus done over time | Scope line and done line, forecast line from throughput |
| Cumulative flow diagram | Stacked count of items per state per day | Widening bands show where work piles up |
| Time in state | Stacked bar per completed story, or aggregate share | The process optimisation chart |
| Throughput | Bar per week | With nature breakdown |
| Aging WIP | Active items by age since started, against the 85th percentile | Flags items likely to be late |
| Estimate vs actual | Scatter, only items with `estimate` | |
| Token rate | One point per item with usage and agent time, per model: x completed date (started for an open item), y the model's tokens per agent hour | By type: epic, story, task |
| Cost | Bar per item with usage, stacked by model | Estimated costs marked |
| Completion against time | Per model, cumulative items done over time | The per model `done` series |
| Completion against cost | Per model, cumulative items done over cumulative cost | The same series, x cost |

## Precision rules

So that `flai stats` and the dashboard agree to the second:

- Percentiles use nearest rank on the sorted values: the p-th percentile is the value at position `ceil(p/100 × n)`.
- The backlog interval starts at `created`. An open item's current interval ends at now. Closed items accrue nothing after `completed`.
- Queue time exists only when both `committed` and `started` exist and `started` is not before `committed`.
- Cycle time exists only when the item was started. An item cancelled from backlog has a lead time but no cycle time and is excluded from cycle aggregates.
- Burn-up counts an item in scope from its `created` day and excludes cancelled items. Cumulative flow uses the item's state at the end of each day (23:59:59 UTC).
- Throughput per week is completed items in the window divided by window days over seven. Weekly buckets are ISO weeks starting Monday.
- Flow efficiency averages `(cycle - blocked) / cycle` over completed items with a positive cycle time.
- Time-in-state share divides total seconds per state by total lead time, over completed items in the window.

## Data access

The dashboard reads items via its server API, which scans `wip/kanban` and `wip/archive` on request and caches by file modification time. `flai stats` prints the same aggregates as a table or JSON so scripts and CI can use them. Both share the same definitions above; the Go implementation is the reference and `flaiover` has a fixture test that checks its numbers against `flai stats --json` output on the sample repo in the template.
