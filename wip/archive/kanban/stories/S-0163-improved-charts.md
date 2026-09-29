---
id: S-0163
type: story
nature: feature
title: Improved Charts
status: done
parent: E-0011
owner: alex
created: 2026-09-29T19:40:33Z
updated: 2026-09-29T21:02:21Z
transitions:
  - to: ready
    at: 2026-09-29T20:14:01Z
    by: alex
  - to: in-progress
    at: 2026-09-29T20:29:22Z
    by: agent-S-0163
  - to: review
    at: 2026-09-29T21:01:26Z
    by: agent-S-0163
  - to: done
    at: 2026-09-29T21:02:21Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side, analytics]
touches: [flaiover/src, flai/cmd, flai/internal/metrics, flai/internal/hostapi, design/system/metrics.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, design/tech/charts.md, docs/users, docs/operators, design/adrs/0053-usage-is-charted-as-spend-over-time-in-buckets-for-every-item-type-at-once-and.md, design/adrs/README.md, design/issues/I-0027-work-items-written-in-the-main-checkout-reach-ci-without-a-markdown-lint.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
usage:
  source: log
  seconds: 1917
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 178
      output: 1236
      cache_read: 19720001
      cache_write: 356115
      cost: 32.726
---
# S-0163 Improved Charts

## Goal

Right now, each story and task tracks how many tokens it has used via front-matter. Unfortunately, the first round of charts designed do a very poor job of representing the efficiency or expenditure.

We need much better and more reasonable charts that help the operator understand the rate of token usage in units of minutes, tasks, stories, and epics. The operator should also be able to understand the cost of completed epics, stories, and tasks (average cost per item type grouped or not grouped by the model used). Finally, the operator should be able to understand the estimated cost accrual rate as a function of time.

These different measures are key because they allow an operator to make important strategic decisions about how they are utilizing their compute, what it's costing, and whether the resulting outcomes are worth while.

The Token Rate chart is the worst of the original charts since agent hours is far too large a unit of time to measure. It should be minutes at the largest. When the original requirements asked for measures "by agent" it meant the specific model rather than the specific agent instance (which is the same as grouping all stories and tasks together).

## Acceptance criteria
- [x] Show the token expenditure rate as:
      1) average tokens spent per day over time
      2) average per work item type (i.e. per story, per task, per epic) over time
      3) average tokens per USD over time
- [x] Show the model use as:
      1) average USD expenditure for each model over time
      2) average USD expenditure per work type over time
      3) average token expenditure per work type over time

## Tasks
- T-0580 The usage metrics define spend per bucket of time, per item, per agent minute, and per dollar, with an ADR
- T-0581 flai stats lays spend out per bucket and per item type, and counts agent time in minutes
- T-0582 The dashboard builds the token and cost charts from the spend series
- T-0583 The charts page offers the new charts with a bucket, a grouping, and a table view each
- T-0584 Every tier passes and the charts are looked at with this repository's numbers

## Notes

Which chart answers each criterion, on the Charts page under usage:

| Criterion | Chart | What it draws |
|-----------|-------|---------------|
| Tokens per day over time | Tokens per day | The tokens of the items done each day, stacked by model, and the mean per day so far |
| Tokens per work item type over time | Tokens per item, comparing item types | The mean tokens per epic, per story, and per task done in each day |
| Tokens per USD over time | Tokens per dollar | Tokens over cost for each day, per model |
| USD for each model over time | Cost per day; Cost per item, comparing models | Each model's dollars per day, and each model's mean dollars per item |
| USD per work type over time | Cost per item, comparing item types | The mean dollars per epic, per story, and per task |
| Tokens per work type over time | Tokens per item | By item type, or by model for one type |

The goal asks for the rate in minutes: Token rate is now tokens per minute of agent work over time, per model. A day can be an hour or a week instead (the `per` control), because every item with usage here was done on one day. An item counts when it entered done, since its usage is one total. The cost of each item and the two completion charts of S-0143 stay; ADR-0053 records all of this.

How the criteria were verified: behaviour tests of the series in `flai/internal/metrics` and of `flai stats --bucket`; unit tests of each chart builder; a test that mounts the charts page and reads its controls, summary, and table; and every new chart drawn by ECharts from this repository's `flai stats --json`, by day and by hour, and from a synthetic month with three models, in light and dark, and looked at. Not done: the page was not opened in a browser against the host's running flai serve.

`flai check --strict` has 0 errors and 4 warnings, none from this story: E-0003, E-0010, and E-0012 are done and not archived, and TH-0032 is answered on an archived story. The markdown lint has one error, in `wip/threads/TH-0035` on main (I-0027). Both stop the smoke tier; every other step of it passes.

Proposed follow-up, not part of this story: the opus and fable colours of S-0143's palette are close for a reader with deuteranopia (the dataviz validator's all-pairs check gives 0.6 in light and 2.2 in dark, where 8 is the target). The new line charts give each model a mark as well as a colour. The stacked bars and the completion charts do not. A story could move fable to a slot that stays apart from opus, sonnet, and haiku.
