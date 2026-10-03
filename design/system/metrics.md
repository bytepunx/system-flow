---
title: Flow metrics
updated: 2026-10-03
status: active
topics: [cli, dashboard, analysis]
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
- **Completion against time and cost**: the items in order of `completed` (then ID), each point carrying `completed`, the ID, and the cumulative count of items, tokens, and cost up to and including it. Per model, the same over the items that model worked on, counting that model's tokens and cost. Reported for scripts; no chart draws it since [ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md).

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
| Minutes per item | Agent seconds over items, in minutes: the mean agent time an item took. Absent also when none of them took agent time ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| Tokens per minute | Tokens over agent minutes |
| Tokens per dollar | Tokens over cost |
| Mean tokens, mean cost | Of a bucket only: the running mean per bucket, the sum from the first bucket of the series to this one over the number of those buckets |

A value is absent when its divisor is zero.

`flai stats --json` carries each item's usage under `items[].usage` and the aggregates under `usage`: `items`, `tokens`, `cost`, `seconds`, `estimated`, `models`, `done`, `by_model`, and since S-0163 `bucket` (`hour`, `day`, or `week`) and `spend`. `spend` has the keys `epic`, `story`, and `task`, each with the values of its items over the window (`items`, `tokens`, `cost`, `seconds`, `estimated`, `tokens_per_item`, `cost_per_item`, `minutes_per_item`, `tokens_per_minute`, `tokens_per_dollar`), the same per model under `models`, and the series under `buckets`: each point has `at`, the same values, `mean_tokens`, `mean_cost`, and its `models`. A rate is `tokens_per_minute`. `tokens_per_hour`, sixty times that, stays beside it on items and models for what was written to flai 1.25.

## Strategic agents (S-0206)

What the planner, the orchestrator, and the analyzer spent, from their activity documents under `wip/agents` ([agent-narrative.md](agent-narrative.md), [ADR-0079](../adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md)). This is apart from item usage: an activity's cost is not on any item.

`flai stats --json` carries them under `strategic`, a list with one entry per document that exists, in the order planner, orchestrator, analyzer, and an empty list when none does. The totals are the document's front matter as written, over all time, not the window's; only `log` is windowed. `flai stats` prints the totals, one line per agent, and nothing when there are no documents. An unreadable document stops `flai stats` with its path and what is wrong, as an unreadable item does.

| Field | Source | Precision |
|-------|--------|-----------|
| `kind` | `kind`: `planner`, `orchestrator`, or `analyzer` | |
| `cost` | `accrued_cost`, in US dollars | Four decimals, as written |
| `seconds` | `accrued_seconds`, wall-clock | Whole seconds |
| `activities` | `tasks_completed`, the number of activities logged | |
| `last_run` | `last_run`, when the newest activity ended; `""` before the first | UTC, to the second |
| `log` | The log entries that ended in the window, from its start to the report's now, both included, oldest first | |
| `log[].at` | When the activity ended, the entry's heading | UTC, to the second |
| `log[].seconds` | The entry's wall-clock seconds | Whole seconds |
| `log[].cost` | The entry's cost, in US dollars | Four decimals, as written |
| `log[].estimated` | `true` when the cost was apportioned or priced rather than reported; absent otherwise | |
| `log[].items` | The IDs of the items the activity touched; `[]` for none | |

Their cost and use per day, beside delivery, are in [Strategic use per day](#strategic-use-per-day).

## Planning, waiting, and claims (S-0205)

What the planner's figures, the operator's cost of delay, waiting, and claims come to, for the charts of E-0016 ([ADR-0081](../adrs/0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md)). Each value is of the report's type, as the other aggregates are: a story's here, an epic's with `--type epic`. Per-item values are under `items[]`, absent when an input is. The aggregates cover the items of the type completed in the window, cancelled ones left out, unless a row says otherwise. A series by the day runs from the day that holds the window's start to the day that holds now, one point per day; one by the week, from the ISO week that holds the window's start to the one that holds now.

### Forecasts and estimates

The planner's `forecast` and the human's `estimate` against what happened ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)). A positive error is later or longer than forecast.

| Per item | Definition |
|----------|------------|
| `forecast_seconds` | `forecast.duration` |
| `forecast_error_seconds` | Cycle time minus `forecast.duration` |
| `delivery_error_seconds` | `completed` minus `forecast.delivery` |
| `estimate_error_seconds` | Cycle time minus `estimate`. `estimate_error` stays beside it, the same over the estimate |

`forecasts` has `forecast`, `delivery`, and `estimate`, each over the items with that error: `count`, `p50_seconds`, and `p85_seconds` of the absolute error, and the same per nature under `by_nature` and per story agent model (`agent.model`, `(none)` without one) under `by_model`. A set with no items has `count` 0 and no percentiles.

### Cost of delay

The `cost_of_delay.value` of an item is what a week of waiting for it costs, in the project's currency ([ADR-0080](../adrs/0080-a-cost-of-delay-stamps-its-inputs-and-its-value-apart.md)). An item waits while it is in `backlog` or `ready`. Its value today is used for all of its history: flai keeps no older values.

| Value | Definition |
|-------|------------|
| `items[].cost_of_delay` | Its `value`; absent without one |
| `items[].cost_of_delay_incurred` | Its value times the seconds it spent in `backlog` or `ready`, up to now, over 604800 (a week) |
| `cost_of_delay.days[].outstanding` | Per column, `backlog`, `ready`, `in-progress`, and `review`, the sum of the values of the items in it at the end of the day (23:59:59 UTC), every column present |
| `cost_of_delay.days[].incurred` | Over every item, its value times the seconds of the day it spent in `backlog` or `ready`, over 604800. Today ends at now |
| `cost_of_delay.weeks[]` | `week`, `start`, and `incurred`, the same over the week's seconds. The window's first week is whole |

Each day point has `date`; each amount is rounded to two decimals once summed. Items without a value add nothing.

### Waiting

The time a story's agent waited for someone else: on its threads while it was in progress, and in review.

| Value | Definition |
|-------|------------|
| A thread's wait | From its first entry to the first later entry by another author. A thread nobody else answered waits until `updated` once resolved, and until now while open |
| `items[].wait_threads_seconds` | The seconds of the union of the waits of the threads anchored to the item or one of its tasks that fall in its `in-progress` intervals. Absent with no such thread |
| `items[].wait_review_seconds` | The seconds it spent in `review`, the open interval up to now. Absent if it was never in review |

`waiting.weeks[]` has `week`, `start`, `items` (those completed in the week), and `threads` and `review`, each with `total_seconds`, their sum over those items, and `mean_seconds`, the sum over `items`.

### Claims and touches

What claims cost in holds and parallelism, and how far a story's declared touches were from what it changed ([ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)).

| Value | Definition |
|-------|------------|
| `items[].held_seconds` | The seconds a story in `ready` was held by the hold rules (`overlap`, `no-touches`, `after`), replayed at every creation and transition of a story or task, up to now, from the states of the time and today's touches and `after`: flai records neither a hold nor older touches. Absent for an item that is not a story or was never in ready |
| `claims.limit` | The board's `in-progress` limit today; absent when it has none |
| `claims.days[]` | `date`, and `in_progress`, the items in `in-progress` at the end of the day |
| `claims.drift[]` | One entry per story of the report that a commit names, by ID: `id`, `committed` (the files its commits changed), `outside` (those under none of its touches), and `unchanged` (its touches no committed file is under), each a sorted list of paths, with `outside_count` and `unchanged_count` |

A story's commits are those on the main branch and the `story/` branches whose subject names it in brackets, `[S-nnnn]`, merges left out; their files under the `wip` folder are left out, since flai writes them. Its touches are its own and those of its tasks not cancelled, a project's name or tag read as its path, as a claim reads them. A file is under a touch that is the file or a folder that holds it. When git cannot be read, `claims.drift` is absent and `flai stats` says so on stderr.

### Strategic use per day

`strategic_days[]` lays the strategic agents' activity beside delivery, for the Strategic Cost and Strategic Use charts. One point per day, each with:

| Field | Definition | Precision |
|-------|------------|-----------|
| `date` | The day | `YYYY-MM-DD` |
| `agents` | Per kind with an entry that ended that day: `cost`, `seconds`, and `estimated` when any entry was | Four decimals; whole seconds |
| `cost`, `seconds` | The sums over the kinds | Four decimals; whole seconds |
| `completed` | The items of the report's type completed that day | |
| `cost_per_item` | The usage cost of those carrying usage, over their number | Four decimals; absent when none carries usage |
| `cycle_time_seconds` | The mean cycle time of those with one | Absent when none has one |

## Charts

Every chart spans the window chosen ([ADR-0054](../adrs/0054-every-chart-spans-the-window-chosen-its-time-axis-runs-from-the-window-s-start.md), S-0166). A time axis runs from the window's start to the report's now, whatever the data: a series by the day from the day that holds the start, a series in buckets from the bucket that holds the start to the one that holds now, with half a bucket either side. A chart per item plots only the items completed in the window, and time in state groups them by the day they were completed; `items` in `flai stats --json` holds every item of the type, and the dashboard picks them.

| Chart | Data | Notes |
|-------|------|-------|
| Kanban board | Current status of active items grouped by column, with age in column and blocked flag | Age in column is now minus last transition |
| Cycle Time | One point per story completed in the window, x completed date, y cycle time, with 50th and 85th percentile lines | Filter by nature |
| Burn-up | Per epic or whole repo, cumulative stories created versus done, per day of the window | Scope line and done line, forecast line from throughput |
| Cumulative Flow | Stacked count of items per state per day of the window | Widening bands show where work piles up |
| Time in State | Stacked bar per day of the window with stories completed: the mean hours per state of those stories ([ADR-0056](../adrs/0056-time-in-state-is-one-stacked-bar-per-day-of-the-window-the-mean-hours-per-state.md)), and the aggregate share | The process optimisation chart. The table lists each story |
| Throughput | Bar per week of the window | With nature breakdown |
| Tokens / Min | Per model, one point per bucket with agent time: the model's tokens per agent minute over the items done in the bucket | Of the report's type. Replaces the point per item in tokens per agent hour (S-0163) |
| Tokens / bucket | Bar per bucket, the tokens of the items done in it, stacked by model, with the running mean per bucket as a line | Of the report's type. Titled by the bucket: Tokens / Day |
| Tokens per item | One point per bucket with items: tokens per item, one series per type | Every type at once ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| Tokens / $ | Per model, one point per bucket with cost: tokens over cost | Of the report's type |
| $ / bucket | Bar per bucket, the cost of the items done in it, stacked by model, with the running mean per bucket as a line | Of the report's type. Titled by the bucket: $ / Day. Estimated costs marked |
| $ / Work Type | One point per bucket with items: cost per item, one series per type | Every type at once ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| $ / Item | Bar per item completed in the window with usage, stacked by model | Estimated costs marked |
| Avg. Time / Model | Per model, one point per bucket in which it worked on items that took agent time: their minutes per item | Of the report's type. With two models or more, a dashed line for all of them ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| Avg. Cost / Model | Per model, one point per bucket in which it worked on items: its cost per item | Of the report's type. With two models or more, a dashed line for all of them ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |

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
- The planning, waiting, and claims values (S-0205) are seconds between timestamps, whole since timestamps are, and their means are not rounded. A day is a UTC day, from 00:00:00 up to, not including, the next; a day's or a week's share of an interval is the part of it inside the day or week, and today's ends at now. Cost of delay amounts are rounded to two decimals once summed, and strategic costs to four.
- A bucket holds the moments from its start up to, not including, the next one's. A week's bucket starts on the Monday of the ISO week, at 00:00:00 UTC. The running mean divides by the number of buckets from the first of the series, empty ones counted.

## Data access

The dashboard reads items via its server API, which scans `wip/kanban` and `wip/archive` on request and caches by file modification time. `flai stats` prints the same aggregates as a table or JSON so scripts and CI can use them. Both share the same definitions above; the Go implementation is the reference and `flaiover` has a fixture test that checks its numbers against `flai stats --json` output on the sample repo in the template.
