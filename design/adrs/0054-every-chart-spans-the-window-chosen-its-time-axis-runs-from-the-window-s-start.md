---
id: ADR-0054
title: "Every chart spans the window chosen: its time axis runs from the window's start to now, and it plots only what the window holds"
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
refines: [ADR-0053]
topics: [cli, dashboard]
---

# ADR-0054 Every chart spans the window chosen: its time axis runs from the window's start to now, and it plots only what the window holds

## Context

The charts page of the dashboard offers a window, from a day to a year, and sends it to `flai stats` (S-0166). Changing it changed nothing visible:

- `items` holds every item of the report's type, and the charts per item (cycle time, time in state, estimate against actual, cost by item) and their tables plotted every completed one.
- Burn-up and cumulative flow ran from the day the first item was created to today, and throughput held only the weeks in which something was done.
- Every time axis was fitted to the points it drew, so a series that began inside the window, such as the spend over time of ADR-0053, drew the same axis for any window that held it.

On this repository on 2026-09-29, `flai stats --since 1d` and `--since 365d` gave the same 166 items, the same 15 days of cumulative flow, and the same burn-up. `design/system/metrics.md` defines each chart's data and changes only with an ADR.

## Decision

Every chart spans the window chosen.

1. **The series over time cover the window.** Burn-up and cumulative flow have one point per day from the day that holds the window's start, or the day the first item was created if that is later, to today. Throughput has one bucket per ISO week from the week that holds the window's start to the week that holds now, a week in which nothing was done among them with zero.
2. **The charts per item plot the items completed in the window.** `items` stays every item of the type, since open items and their ages need it. The dashboard plots, and lists in the tables, only those completed between the window's start and the report's now.
3. **Each time axis is the window.** It runs from the window's start to the report's now, whatever the data. A series by the day starts on the day that holds the window's start. A series in buckets, as in ADR-0053, runs from the bucket that holds the window's start to the one that holds now, with half a bucket either side so that the first and last bars are whole.
4. **The spend over time of ADR-0053 is unchanged.** Its series still start at the first bucket with spend, so the running mean still counts from there. Only its axis widens to the window.

## Consequences

- A narrower window shows fewer points on an axis that ends where the window does. A wider one shows the empty time before the first item or the first spend, rather than hiding it.
- Throughput over a year is 53 or 54 bars, most of them zero before a project's first week.
- `flai stats --json` grows by at most one throughput bucket per week of the window, and shrinks when the window is shorter than the project's history.
- The page draws only the answer to the window it asked for last, so a slow answer to an earlier choice cannot replace a later one.
- A dashboard on a flai older than this draws the older flai's burn-up, cumulative flow, and throughput on an axis that spans the window, so their points before the window fall off it.

## Alternatives considered

- Filtering `items` by the window in flai: aging and the open items need every item, and a second list would duplicate the first.
- Starting burn-up and cumulative flow at the window's start even before the first item existed: the days would carry only zeros, and the axis already shows the window.
- Starting the spend series at the window's start: the running mean per bucket would then count the empty buckets before the first spend and change its meaning.
- Trimming every series in the dashboard alone: `flai stats` is the reference implementation, and scripts would still get the whole history for a window of a day.
