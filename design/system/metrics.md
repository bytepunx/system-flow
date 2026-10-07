---
title: Flow metrics
updated: 2026-10-07
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

What strategic agents spent on an item, its `usage.strategic` ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)), is reported beside the agents' figures and apart from them: it is never added to any agent figure, total, per-model value, series, or rate below, and an item that carries only `strategic` has its agents' values at zero and is left out of every agent aggregate, as an empty usage is. [Strategic usage](#strategic-usage) defines its values.

| Value | Definition |
|-------|------------|
| Tokens | Sum over the item's models of `input + output + cache_read + cache_write` |
| Cost | Sum over the item's models of `cost`, in US dollars |
| Agent minutes | `seconds / 60` |
| Token rate | Tokens over agent minutes; absent when `seconds` is 0 |
| Per model | The same for each model the item lists: its tokens, its cost, and its tokens over the item's agent minutes |
| Estimated | The item's `estimated` |
| Empty wakes | The item's `usage.empty_wakes`, 0 without one: on a story, the [empty wakes](#waiting) its agents' logs held when flai serve last measured it; on an epic, its stories' summed; a task carries none ([ADR-0105](../adrs/0105-a-story-s-empty-wakes-the-wait-for-events-calls-of-its-agents-that-timed-out.md), S-0272). Their aggregates are under [Waiting](#waiting), not here |

Aggregates cover the items of the report's type that entered `done` in the window and carry usage; cancelled items are left out:

- **Totals**: items, tokens, cost, agent seconds, and whether any is estimated.
- **Per model**: for each model, the items it worked on, its tokens and cost over them, and its tokens over the agent minutes of those items.
- **Completion against time and cost**: the items in order of `completed` (then ID), each point carrying `completed`, the ID, and the cumulative count of items, tokens, and cost up to and including it. Per model, the same over the items that model worked on, counting that model's tokens and cost. Reported for scripts; no chart draws it since [ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md).

### Strategic usage

What the planner and the orchestrator spent on an item, and the analyzer through an issue (below), is charged to it and to the items above it when the activity is logged: a planner activity to the item it planned ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)), an orchestrator activity split evenly between the work items it named ([ADR-0095](../adrs/0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md)). An entry's tokens and cost are those of its `models`, as an agent model's are.

| Value | Definition |
|-------|------------|
| `items[].usage.strategic` | One entry per kind the item's `usage.strategic` lists, in its order (`planner`, `orchestrator`, `analyzer`): `kind`, `tokens`, `cost`, `seconds`, and `estimated`, `true` since every charge is apportioned. `[]` when the item's usage has no `strategic` |
| `usage.strategic` | Over the items of the report's type that entered `done` in the window and carry `strategic`, cancelled ones left out: `items`, their count; `tokens`, `cost`, and `seconds`, the sums over every kind; `estimated`, whether any entry is; and `kinds`, the same per kind (`kind`, `items` that carry it, `tokens`, `cost`, `seconds`), in the order planner, orchestrator, analyzer, only the kinds present. Always present: zeros and `kinds: []` when no item carries any |

An item with `usage` that carries only `strategic` has `items[].usage` with `tokens`, `cost`, and `seconds` 0 and `models: []`.

Since S-0227 ([ADR-0100](../adrs/0100-an-analyzer-activity-s-usage-is-charged-evenly-to-the-issues-it-names-under.md)) an analyzer activity is split evenly between the issues it named and charged to each, under the issue's own `usage.strategic`, in the shape an item's has. When `flai issue story` makes a story from an issue, the story carries the issue's entries, and so does its epic; the issue keeps them. `flai stats` counts each spend once in its totals:

- An issue's usage counts while no story was made from it.
- Once a story was made from it, the story's entry counts, under the items, and the issue's does not.
- A story was made from an issue when the issue's `## Remediation` holds the line `Story S-nnnn remediates this issue, created from it at <time>.`, which `flai issue story` writes, and that story is among the items the report reads, archived ones included. With more than one such line, the first story the items hold is the one. A story that only links or names the issue carries none of its usage.

| Value | Definition |
|-------|------------|
| `strategic_issues` | One entry per issue that carries `usage.strategic`, whatever its status, in order of ID: `id`, `title`, `status`; `story`, the story made from it, `""` for none; `counted`, `true` while no story was made from it; and `strategic`, its entries as `items[].usage.strategic` gives an item's. All time, not windowed. `[]` when no issue carries any, or the project has no issues folder |
| `strategic[].issues` | What the issues carry of the kind that no story does: the sum, over the issues with `counted` true, of the cost of their entry of the kind and of its `seconds`. `strategic[].project` leaves it out ([Strategic agents](#strategic-agents-s-0206)) |

An issue's usage is never in the agents' figures, the per-model figures, `usage.strategic`, `usage.spend`, or `strategic_days`: an issue is not an item and is never done in a window. Its usage reaches those only through the story made from it.

### Cost per agent hour and expected cost

| Value | Definition |
|-------|------------|
| `usage.cost_per_agent_hour` | The project's mean cost of an hour of agent work: the sum of the agents' cost over the sum of their `seconds` in hours, over every story, whatever its status, archived ones included, whose `usage` has `source: log` and `seconds` above 0. Not windowed, and of stories whatever the report's type. What strategic agents spent is left out. Absent when no story qualifies |
| `items[].expected_cost` | What the item is expected to cost: `cost`, its `forecast.duration`, or else its `estimate`, in hours, times `usage.cost_per_agent_hour`; `from`, `forecast` or `estimate`, the one used; and `estimated: true`. Present on any item, open or closed. Absent when the item has neither duration, or there is no cost per agent hour |

### Spend over time (S-0163)

Spend is laid out in buckets, for epics, for stories, and for tasks, whatever the report's type. Each type is summed over its own items: a story's usage holds its tasks' and is never added to them.

| Term | Definition |
|------|------------|
| Bucket | An hour, a day, or an ISO week starting Monday, in UTC, named by its start. A day unless asked otherwise. An hour needs a window of 31 days or less |
| An item's bucket | The one that holds the moment it entered `done`. Its whole usage counts there, its agents' figures and its strategic ones apart |
| Series | For a type, one point per bucket from the bucket of the first item of that type done in the window with usage, agents' or strategic, to the bucket that holds now. An item that carries only strategic usage starts the series as one with agents' usage does, so its bucket is a point whose agents' values are zero. A bucket with no such item is a point of zeros. A type with no such item has no points |

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

Each type over the window, and each point of its series, also has `strategic`: `items`, the items among them that carry strategic usage, and `tokens`, `cost`, and `seconds`, what strategic agents spent on those, summed over every kind ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)). It is always present, zeros when nothing was spent, and no other value of the type or the point counts it: the values above, the models, and the running means are the agents' alone.

`flai stats --json` carries each item's usage under `items[].usage`, with its `strategic` list and, since S-0272, `empty_wakes`, and its expected cost under `items[].expected_cost`; and the aggregates under `usage`: `items`, `tokens`, `cost`, `seconds`, `estimated`, `models`, `done`, `by_model`, since S-0163 `bucket` (`hour`, `day`, or `week`) and `spend`, and since S-0225 `strategic` and `cost_per_agent_hour`. `spend` has the keys `epic`, `story`, and `task`, each with the values of its items over the window (`items`, `tokens`, `cost`, `seconds`, `estimated`, `tokens_per_item`, `cost_per_item`, `minutes_per_item`, `tokens_per_minute`, `tokens_per_dollar`), the same per model under `models`, its `strategic`, and the series under `buckets`: each point has `at`, the same values, `mean_tokens`, `mean_cost`, its `models`, and its `strategic`. A rate is `tokens_per_minute`. `tokens_per_hour`, sixty times that, stays beside it on items and models for what was written to flai 1.25.

## Strategic agents (S-0206)

What the planner, the orchestrator, and the analyzer spent, from their activity documents under `wip/agents` ([agent-narrative.md](agent-narrative.md), [ADR-0079](../adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md)). A planner activity's cost is also charged to the item it planned and the items above it, and an orchestrator activity's split evenly between the work items it named and charged to each and the items above it, under their `usage.strategic` ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md), [ADR-0095](../adrs/0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md), [Strategic usage](#strategic-usage)): the activity documents and the items are two views of one spend, never added together. Since S-0227 an analyzer activity's cost is charged to the issues it named, and carried to the story made from an issue ([Strategic usage](#strategic-usage)). What of a document's totals neither the items nor the issues carry is its kind's project strategic total: the activities that named no work item or issue, the planner's with no planned item, those logged before their kind was charged, and any charge that failed. Each entry splits its totals between `items`, `issues`, and `project`, so that the three add up to the document's, to a hundredth of a cent; the project total is worked out here, not written anywhere.

`flai stats --json` carries them under `strategic`, a list with one entry per document that exists, in the order planner, orchestrator, analyzer, and an empty list when none does. The totals are the document's front matter as written, over all time, not the window's; only `log` is windowed. `flai stats` prints the totals, one line per agent with what the items carry, what the issues carry, and the project total beside them, and nothing when there are no documents; then one line per issue that carries strategic usage, with its cost and seconds per kind and whether it counts on the issue or on its story. An unreadable document stops `flai stats` with its path and what is wrong, as an unreadable item does.

| Field | Source | Precision |
|-------|--------|-----------|
| `kind` | `kind`: `planner`, `orchestrator`, or `analyzer` | |
| `cost` | `accrued_cost`, in US dollars | Four decimals, as written |
| `seconds` | `accrued_seconds`, wall-clock | Whole seconds |
| `activities` | `tasks_completed`, the number of activities logged | |
| `last_run` | `last_run`, when the newest activity ended; `""` before the first | UTC, to the second |
| `items.cost`, `items.seconds` | What the items carry of the kind: the sum, over the items at the top of the hierarchy, those with no parent or a parent that does not exist, archived ones included, of the cost of their `usage.strategic` entry of the kind, its models' costs summed, and of its `seconds`. Each charge is on its item and every item above it, so the items at the top carry every charge once. Zeros when no item carries the kind | Four decimals; whole seconds |
| `issues.cost`, `issues.seconds` | What the issues carry of the kind that no story does: the sum, over the issues from which no story among the items was made, of the cost of their `usage.strategic` entry of the kind, its models' costs summed, and of its `seconds` (S-0227). Zeros when no such issue carries the kind | Four decimals; whole seconds |
| `project.cost`, `project.seconds` | The kind's project strategic total: `cost` less `items.cost` and `issues.cost`, and `seconds` less `items.seconds` and `issues.seconds`, never below zero. Zero below zero, when the items and the issues carry more than the document holds, as when a document is restored from an older commit | Four decimals; whole seconds |
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
| `model` | `agent.model`, `(none)` without an agent or a model, on every item ([ADR-0111](../adrs/0111-flai-stats-gives-each-item-the-model-its-forecasts-are-grouped-under.md)) |

`forecasts` has `forecast`, `delivery`, and `estimate`, each over the items with that error: `count`, `p50_seconds`, and `p85_seconds` of the absolute error, and the same per nature under `by_nature` and per story agent model, the item's `model`, under `by_model`. A set with no items has `count` 0 and no percentiles.

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

The time a story's agent waited for someone else: on its threads while it was in progress, and in review. The orchestrator's replies (S-0220) are told apart from the rest ([ADR-0091](../adrs/0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md)): a recommendation ends no wait, and the waits that the orchestrator's answers and the operator's confirmations of a recommendation ended are counted apart. Beside them, the empty wakes count the turns a story's agent spent waking from `wait_for_events` to nothing (S-0272), the turns that S-0272's ending on an open question saves.

| Value | Definition |
|-------|------------|
| A thread's wait | From its first entry to the first later entry by another author that is not a recommendation ([ADR-0090](../adrs/0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md)). A thread nobody else answered waits until `updated` once resolved, and until now while open; one whose only replies are recommendations waits on in the same way, until the operator confirms one or someone answers |
| Who ended a wait | The orchestrator, when the entry that ended it is by `orchestrator`, the author name of the orchestrator's run (ADR-0087); confirmed, when that entry confirms a recommendation, a line of it starting with the words `Confirmed the recommendation of` and a space; anyone else otherwise |
| `items[].wait_threads_seconds` | The seconds of the union of the waits of the threads anchored to the item or one of its tasks that fall in its `in-progress` intervals, whoever ended them. Absent with no such thread |
| `items[].wait_threads_orchestrator_seconds` | The same over only the waits the orchestrator ended: the seconds of their union that fall in the item's `in-progress` intervals. Time in which such a wait overlaps one someone else ended counts here, so it is never more than `wait_threads_seconds`. Present, 0 or more, whenever `wait_threads_seconds` is |
| `items[].wait_review_seconds` | The seconds it spent in `review`, the open interval up to now. Absent if it was never in review |
| An empty wake | A call by a story's agent itself to the flai MCP tool `wait_for_events` whose result reports `timed_out: true`, no events, and no changed paths ([ADR-0105](../adrs/0105-a-story-s-empty-wakes-the-wait-for-events-calls-of-its-agents-that-timed-out.md), S-0272). In Claude Code's stream-json log, a `tool_use` named `mcp__flai__wait_for_events` in an event without `parent_tool_use_id`, whose `tool_result` (by `tool_use_id`) is not `is_error` and whose JSON has `timed_out` true, `events` and `changed` empty, and `end` not true. A sub-agent's call, a call refused or failed, a call that answered `end: true`, and a call with no result in the log are not. flai serve counts them over every log of the agents it started for a story each time it measures the story's usage, and writes the count as `usage.empty_wakes`; `flai stats` reads that, never the logs |
| `items[].usage.empty_wakes` | The item's `usage.empty_wakes`, 0 when its usage has none; with `items[].usage`, absent when the item carries no usage ([Usage](#usage-s-0143-s-0163)) |
| `waiting.empty_wakes` | `count`, the sum of `items[].usage.empty_wakes` over the items of the report's type completed in the window, cancelled ones left out, the items `waiting.weeks[]` counts; and `mean`, `count` over the number of those items that carry agents' usage, those `usage.items` counts, absent when none does. Always present, `count` 0 when none carries any. Stories measured before S-0272's release carry none and count 0 |

`waiting.weeks[]` has `week`, `start`, `items` (those completed in the week), `empty_wakes`, the sum of `items[].usage.empty_wakes` over those items, 0 when none, and `threads` and `review`, each with `total_seconds`, their sum over those items, and `mean_seconds`, the sum over `items`. `threads` also has `orchestrator` and `confirmed`, present in every week, each with `count`, the number of the waits of those items that the orchestrator, or a confirmation, ended and that have some part in the item's `in-progress` intervals, and `total_seconds`, the sum over those items of the seconds of the union of those waits that fall in the item's `in-progress` intervals; for `orchestrator` that is the sum of `wait_threads_orchestrator_seconds`. Each part is a union of its own waits, so with overlapping waits the parts can add up to more than `threads.total_seconds`; each is at most it.

`flai stats` prints the waiting section only when an item of the type was completed in the window. After its lines for threads and review it prints one line for the empty wakes, indented and labelled as they are, `empty wakes  total 12 · mean 1.5 per item with usage`: `waiting.empty_wakes.count`, then its `mean` to one decimal, the part from the `·` on left out when `mean` is absent.

### Claims and touches

What claims cost in holds and parallelism, and how far a story's declared touches were from what it changed ([ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)). Since S-0214, the stories held per day, the time held per week by reason, and the weekly share of stories whose touches were exact, for the Parallelism, Hold Time, and Touches Drift charts ([ADR-0113](../adrs/0113-flai-stats-reports-held-stories-per-day-held-time-by-reason-per-week-and-the.md)).

| Value | Definition | Precision |
|-------|------------|-----------|
| `items[].held_seconds` | The seconds a story in `ready` was held by the hold rules (`overlap`, `no-touches`, `after`), replayed at every creation and transition of a story or task, up to now, from the states of the time and today's touches and `after`: flai records neither a hold nor older touches. Absent for an item that is not a story or was never in ready | Whole seconds |
| `claims.limit` | The board's `in-progress` limit today; absent when it has none | |
| `claims.days[]` | `date`; `in_progress`, the items in `in-progress` at the end of the day (23:59:59 UTC); and `held`, the stories in `ready` that the hold rules hold then, as the replay of `items[].held_seconds` finds them: under the states of the last creation or transition at or before that moment, today's up to now. Only stories are held, so `held` is 0 in a report of another type | `YYYY-MM-DD`; whole counts |
| `claims.weeks[]` | One point per ISO week: `week` and `start`, as `cost_of_delay.weeks[]` has them, then `held_seconds`, `stories`, `exact`, and `exact_share` | `start` is `YYYY-MM-DD`, a Monday |
| `claims.weeks[].held_seconds` | Per reason, `overlap`, `after`, and `no-touches`, the seconds of the week that stories in `ready` spent held under it, summed over the stories of the report's type, as the replay of `items[].held_seconds` finds them; every reason present, 0 when none. A hold counts under the one reason it is named by, `held (after)` or the like: `after` when a story it names in `after` is not done, whatever else holds it; otherwise the reason of the first story in progress, by ID, whose claim holds it, `no-touches` when that story's claim or its own is empty and `overlap` when the two overlap. So the reasons sum to the time held in the week. The window's first week is whole, and this week ends at now | Whole seconds |
| `claims.weeks[].stories` | The stories of `claims.drift[]` completed in the week, from the window's start to now, cancelled ones left out | Whole count |
| `claims.weeks[].exact` | Those of them whose touches were exact: `outside_count` and `unchanged_count` both 0 | Whole count |
| `claims.weeks[].exact_share` | `exact` over `stories`. Absent when `stories` is 0: a week without such a story has no share | From 0 to 1, not rounded |
| `claims.drift[]` | One entry per story of the report that a commit names, by ID: `id`, `committed` (the files its commits changed), `outside` (those under none of its touches), and `unchanged` (its touches no committed file is under), each a sorted list of paths, with `outside_count` and `unchanged_count` | Sorted lists of paths; whole counts |

A story's commits are those on the main branch and the `story/` branches whose subject names it in brackets, `[S-nnnn]`, merges left out; their files under the `wip` folder are left out, since flai writes them. Its touches are its own and those of its tasks not cancelled, a project's name or tag read as its path, as a claim reads them. A file is under a touch that is the file or a folder that holds it. When git cannot be read, `claims.drift` is absent and `flai stats` says so on stderr, and so are `stories`, `exact`, and `exact_share` in every week; `held_seconds` is still there.

### Strategic use per day

`strategic_days[]` lays the strategic agents' activity beside delivery, for the Strategic Cost and Strategic Use charts. One point per day, each with:

| Field | Definition | Precision |
|-------|------------|-----------|
| `date` | The day | `YYYY-MM-DD` |
| `agents` | Per kind with an entry that ended that day: `cost`, `seconds`, and `estimated` when any entry was | Four decimals; whole seconds |
| `cost`, `seconds` | The sums over the kinds | Four decimals; whole seconds |
| `completed` | The items of the report's type completed that day | |
| `cost_per_item` | The agents' usage cost of those on which agents spent, over their number; strategic usage is left out | Four decimals; absent when agents spent on none |
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
| $ / bucket | Bar per bucket, the cost of the items done in it, stacked by model, with the running mean per bucket as a line. What strategic agents spent on those items, `buckets[].strategic.cost`, is stacked on top as a series of its own, `strategic`, when any bucket has it ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)) | Of the report's type. Titled by the bucket: $ / Day. Estimated costs marked; strategic cost is estimated. The running mean is the agents' |
| $ / Work Type | One point per bucket with items: cost per item, one series per type | Every type at once ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| $ / Item | Bar per item completed in the window with usage, stacked by model. What strategic agents spent on it, its `usage.strategic` costs summed, is stacked on top as a series of its own, `strategic`, when any item has it; an item with only strategic usage has that alone ([ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md)) | Estimated costs marked; strategic cost is estimated. An item's total is the agents' |
| Avg. Time / Model | Per model, one point per bucket in which it worked on items that took agent time: their minutes per item | Of the report's type. With two models or more, a dashed line for all of them ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| Avg. Cost / Model | Per model, one point per bucket in which it worked on items: its cost per item | Of the report's type. With two models or more, a dashed line for all of them ([ADR-0057](../adrs/0057-the-dashboard-charts-agent-time-and-cost-per-item-per-model-over-time-instead.md)) |
| Forecast Accuracy | One point per story done in the window with `forecast_error_seconds`, x completed, y that error; `estimate_error_seconds` as a second series. Lines at plus and minus the p50 and p85 of `forecasts.forecast` | Filter by nature and model. Under one filter the lines are its `by_nature` or `by_model` entry; under both, the nearest-rank percentiles of the points shown ([Forecasts and estimates](#forecasts-and-estimates)) |
| Delivery Accuracy | One point per story done in the window with `delivery_error_seconds`, in days, x completed. Lines at plus and minus the p50 and p85 of `forecasts.delivery`. On a second axis, the share of those stories with an error of 0 or less, per ISO week of the window | Filter by nature and model. A week with no such story has no share, a gap, not 0 |
| Forecast Error / Model | Per model, the item's `model`, one point per bucket: the p50 of the absolute `forecast_error_seconds` of its stories done in the bucket | Filter by nature. A bucket with no such story is a gap ([ADR-0111](../adrs/0111-flai-stats-gives-each-item-the-model-its-forecasts-are-grouped-under.md)) |
| Parallelism | Per day of the window, `claims.days[].in_progress`, with `held` as a second series, and a line at `claims.limit` | No limit line when `claims.limit` is absent ([Claims and touches](#claims-and-touches)) |
| Hold Time | Stacked bar per week of the window, `claims.weeks[].held_seconds` in hours, one series per reason: `overlap`, `after`, and `no-touches`, the empty claim | A hold with more than one reason counts under the one it is named by, so a bar's height is the time held ([ADR-0113](../adrs/0113-flai-stats-reports-held-stories-per-day-held-time-by-reason-per-week-and-the.md)) |
| Touches Drift | Stacked bar per story of `claims.drift[]` completed in the window, x its `completed` in `items[]`: `outside_count` and `unchanged_count`. On a second axis, `claims.weeks[].exact_share` per week | Cancelled stories left out. A week with no such story has no share, a gap, not 0. Without `claims.drift`, the chart says git could not be read |

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
- The planning, waiting, and claims values (S-0205) are seconds between timestamps, whole since timestamps are, and their means are not rounded. Empty wakes, held stories, and stories with exact touches are whole counts, and the mean of the empty wakes and the share of exact touches are not rounded either. A day is a UTC day, from 00:00:00 up to, not including, the next; a day's or a week's share of an interval is the part of it inside the day or week, and today's ends at now. Cost of delay amounts are rounded to two decimals once summed, and strategic costs to four.
- Strategic usage costs, per item, in totals, per kind, and in spend, are rounded to four decimals once summed, as activity costs are written; their tokens and seconds are whole. The cost per agent hour is rounded to four decimals, and an expected cost is the duration in hours times that rounded rate, rounded to four decimals, so that it can be recomputed from the JSON. Durations are Go durations, as `forecast_seconds` and `estimate_seconds` read them.
- A bucket holds the moments from its start up to, not including, the next one's. A week's bucket starts on the Monday of the ISO week, at 00:00:00 UTC. The running mean divides by the number of buckets from the first of the series, empty ones counted.

## Data access

The dashboard reads items via its server API, which scans `wip/kanban` and `wip/archive` on request and caches by file modification time. `flai stats` prints the same aggregates as a table or JSON so scripts and CI can use them. Both share the same definitions above; the Go implementation is the reference and `flaiover` has a fixture test that checks its numbers against `flai stats --json` output on the sample repo in the template.
