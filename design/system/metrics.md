---
title: Flow metrics
updated: 2026-09-15
status: active
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

## Data access

The dashboard reads items via its server API, which scans `wip/kanban` and `wip/archive` on request and caches by file modification time. `flai stats` prints the same aggregates as a table or JSON so scripts and CI can use them. Both share the same definitions above; the Go implementation is the reference and `flaiover` has a fixture test that checks its numbers against `flai stats --json` output on the sample repo in the template.
