---
id: ADR-0053
title: "Usage is charted as spend over time in buckets, for every item type at once, and agent time is counted in minutes"
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0051]
topics: [cli, dashboard]
---

# ADR-0053 Usage is charted as spend over time in buckets, for every item type at once, and agent time is counted in minutes

## Context

ADR-0051 has work items record the tokens and cost their agents spent, and charts them by model: a token rate in tokens per hour of agent work, one point per item; the cost of each item; and items done against time and against cost. The designer found that these say little about how fast tokens and money are spent, or what an item costs on average (S-0163, E-0011):

- An hour of agent work is too large a unit. The 13 stories done with usage in this repository on 2026-09-29 took 4 to 20 minutes of agent work each, so every rate was an extrapolation to a period no item lasted. A minute is the largest unit that fits.
- Nothing shows what is spent per day, what an item of each type costs on average, or how many tokens a dollar buys, and how each changes over time.
- "By agent" in the first request meant the model that did the work, not the agent's session.

An item's `usage` is a total with its seconds of agent work. Nothing in it says when, within the item, the tokens were spent. The only moment an item's spend can be placed at is one of its transitions.

`design/system/metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR.

## Decision

Usage is charted as spend over time, in buckets, for each item type, and agent time is counted in minutes.

1. **Placed at completion, in buckets.** An item's whole usage is placed at the moment it entered done. A bucket is an hour, a day, or an ISO week starting Monday, in UTC, a day unless asked otherwise; an hour is offered only over a window of 31 days or less. The series of a type runs from the bucket of the first item of that type done in the window that carries usage to the bucket that holds now, and a bucket in which nothing was done is in it with zeros.
2. **Every type at once.** The report carries the series of epics, of stories, and of tasks, whatever type the rest of it is about, so that an average per item can be compared across types. Each type is summed over its own items only: a story's usage holds its tasks', and is never added to them.
3. **What a bucket says.** The items done in it, their tokens, cost, and seconds of agent work, and whether any of it is estimated; from these the tokens per item and cost per item (their means over the items), the tokens per minute of agent work, and the tokens per dollar, each absent when its divisor is zero; the same for each model over the items it worked on; and the running mean of tokens and of cost per bucket, from the first bucket of the series to this one.
4. **Minutes.** The token rate is tokens over minutes of agent work. `tokens_per_minute` is reported wherever `tokens_per_hour` is; `tokens_per_hour` stays in `flai stats --json`, for the dashboards and scripts written to flai 1.25, and nothing charts it.
5. **The charts.** The token rate becomes tokens per agent minute over time, per model. Five charts are added: tokens per bucket and cost per bucket, stacked by model, each with its running mean; tokens per item and cost per item, one series per type, or per model for one type; and tokens per dollar, per model. The cost of each item and the two completion charts of ADR-0051 stay as they are.

## Consequences

- The operator reads a rate of spend (per hour, day, or week, and per agent minute), an average per story, task, and epic, and what a dollar buys, each as it changes over time and each by model.
- A story accepted a day after it was worked counts on the day it was accepted. The buckets say when work was delivered, not when the tokens were burned.
- Items that carry usage and are not done, or were cancelled, are in none of these series, as with the aggregates of ADR-0051: what an abandoned item cost is not shown.
- The report grows by one point per bucket and type. Over 365 days by day that is at most 1095 points.
- A dashboard newer than the host's flai gets no series and says which flai it needs. A dashboard older than the host's flai draws its charts as before, since nothing was taken out of the report.
- `flai stats` gains `--bucket` and `stats.get` gains `bucket`.

## Alternatives considered

- Spreading an item's usage over the intervals it was in progress: more faithful to when tokens were spent, but the usage does not say how they fell within those intervals, a story's seconds of agent work are a small part of its time in progress, and a point would then no longer be a set of items with a mean.
- Reading the agents' logs again for a series by the minute: only `flai serve` has the logs, they are not kept for ever, and work items are the one source the metrics are derived from.
- A trailing mean over seven buckets: it needs the buckets before the window and a length per bucket size. The running mean needs neither and ends at the mean of the window.
- Buckets by day only: every item with usage in this repository was done on one day, so each chart would be one point.
- Removing `tokens_per_hour`: it breaks the dashboard released with flai 1.25 for no gain.
- Removing the charts of ADR-0051: the designer named the token rate as wrong, not the others. They can go when asked.
- The dashboard asking for three reports, one per type: three requests and three computations for one chart, and `flai stats --json` would not carry what the dashboard draws.
