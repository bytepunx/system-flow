---
id: S-0163
type: story
nature: feature
title: Improved Charts
status: ready
parent: E-0011
owner: alex
created: 2026-09-29T19:40:33Z
updated: 2026-09-29T20:14:01Z
transitions:
  - to: ready
    at: 2026-09-29T20:14:01Z
    by: alex
tags: [dashboard, cli]
topics: [client-side, server-side, analytics]
touches: [flaiover/src, flai/cmd]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: high
---
# S-0163 Improved Charts

## Goal

Right now, each story and task tracks how many tokens it has used via front-matter. Unfortunately, the first round of charts designed do a very poor job of representing the efficiency or expenditure.

We need much better and more reasonable charts that help the operator understand the rate of token usage in units of minutes, tasks, stories, and epics. The operator should also be able to understand the cost of completed epics, stories, and tasks (average cost per item type grouped or not grouped by the model used). Finally, the operator should be able to understand the estimated cost accrual rate as a function of time.

These different measures are key because they allow an operator to make important strategic decisions about how they are utilizing their compute, what it's costing, and whether the resulting outcomes are worth while.

The Token Rate chart is the worst of the original charts since agent hours is far too large a unit of time to measure. It should be minutes at the largest. When the original requirements asked for measures "by agent" it meant the specific model rather than the specific agent instance (which is the same as grouping all stories and tasks together).

## Acceptance criteria
- [ ] Show the token expenditure rate as:
      1) average tokens spent per day over time
      2) average per work item type (i.e. per story, per task, per epic) over time
      3) average tokens per USD over time
- [ ] Show the model use as:
      1) average USD expenditure for each model over time
      2) average USD expenditure per work type over time
      3) average token expenditure per work type over time

## Tasks

## Notes
