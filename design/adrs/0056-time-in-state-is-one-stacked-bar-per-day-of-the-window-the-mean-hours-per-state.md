---
id: ADR-0056
title: "Time in state is one stacked bar per day of the window, the mean hours per state of the items completed that day"
status: accepted
date: 2026-09-30
supersedes: []
superseded_by: []
refines: [ADR-0054]
topics: [dashboard]
---

# ADR-0056 Time in state is one stacked bar per day of the window, the mean hours per state of the items completed that day

## Context

ADR-0054 has every chart span the window chosen: a time axis runs from the window's start to the report's now. Time in state had no time axis. It drew one stacked bar per item completed in the window, on a category axis of item IDs. Changing the window changed which items it held, but the axis showed no time range (S-0168).

A bar per item placed at its completion time would give it a time axis. On this repository items are often completed minutes apart, and on 2026-09-29, 30 stories were completed in one day. Their bars would overlap or shrink to hairlines on any window longer than a day. The designer chose a bar per day in TH-0039.

## Decision

Time in state is one stacked bar per UTC day of the window in which items were completed.

1. **A day's bar is a mean.** Each state's part of the bar is the mean hours that the items completed that day spent in the state. A day with nothing completed has no bar.
2. **The axis is the window, in days.** It runs from the day that holds the window's start to the day that holds the report's now, with half a day either side, as the series in buckets of ADR-0053 and ADR-0054 do.
3. **The items stay readable.** A bar's tooltip names the items it averages. The table under the chart keeps one row per item completed in the window. The share of lead time per state under the chart is unchanged.

`flai stats --json` is unchanged: the dashboard groups the items it already sends.

## Consequences

- Changing the window moves the axis and changes the days drawn.
- A bar no longer shows one item's time in state. The table does, and the tooltip says which items a bar averages.
- A day with one slow item and many fast ones shows the mean, which the slow item pulls up. Cycle time shows each item on its own.

## Alternatives considered

- One stacked bar per item at its completion time: closest to the chart before, but bars completed close together overlap or become too thin to read.
- Bars by the week: they would read on a year, but hide the day-to-day change on the windows of a day to a month, which are used most.
- A median per day: steadier against one slow item, but the parts of a stacked bar of medians do not add up to a meaningful whole.
