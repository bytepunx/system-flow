---
id: ADR-0081
title: "flai stats reports forecast error, cost of delay, waiting, holds, touches drift, and strategic use per day from the files and git"
status: accepted
date: 2026-10-03
supersedes: []
superseded_by: []
refines: [ADR-0074]
---

# ADR-0081 flai stats reports forecast error, cost of delay, waiting, holds, touches drift, and strategic use per day from the files and git

## Context

E-0016 adds a planner, an orchestrator, and an analyzer, and charts that judge them: forecast against actual, cost of delay outstanding and incurred, time agents spend waiting, holds and touches drift, and what the strategic agents cost against delivery (S-0212 to S-0216). `design/system/metrics.md` is the contract between `flai stats` and the dashboard, so the numbers are defined there first (S-0205). Some inputs have no history: an item keeps only its current `cost_of_delay.value`, `touches`, and `after`, and the launcher records no hold.

## Decision

`flai stats` reports forecast error, cost of delay, waiting, holds, touches drift, and strategic use per day, derived from the work item files, the threads, the activity documents, the board, and git, as `metrics.md` § Planning, waiting, and claims defines them.

- The errors are in seconds, keyed with the unit as every duration is: `forecast_error_seconds`, `delivery_error_seconds`, and `estimate_error_seconds`. The relative `estimate_error` stays as it was, beside them. Aggregates are p50 and p85 of the absolute error, by nature and by `agent.model`.
- Cost of delay is incurred while an item is in `backlog` or `ready`, at its current value per week, prorated by the second. Amounts are rounded to two decimals once summed.
- Waiting on a thread runs from its first entry to the first later entry by another author, and counts while the story is in progress, overlapping threads once. Waiting on review is the time in `review`.
- Held time is replayed: at every transition of any item, the hold rules of ADR-0046 and S-0130 are applied to the states of the time and today's touches and `after`.
- Touches drift compares the files of the commits whose subject names the story, merges and the `wip` folder left out, with its touches and its tasks', read as a claim reads them.
- Strategic use per day sums the activity log entries by the day they ended, beside the usage cost and cycle time per item completed that day.

## Consequences

- The dashboard charts of E-0016 read `flai stats --json` and need no figures of their own.
- `flai stats` reads threads, the board, and git, as well as items; a project that is not a git repository has no drift, and stats says so.
- A value, touches, or `after` changed late rewrites the history these metrics derive from. Recording holds or older values would make them exact; nothing needs that yet.
- Held time costs a pass over the items per transition while a story was ready.

## Alternatives considered

- **Record holds from the launcher and replay nothing.** Exact from now on, but no history before it, and it needs flai serve on the host to have run throughout.
- **`forecast_error` and friends without the unit.** The criteria's words, but every other duration carries `_seconds`, and `estimate_error` already names a ratio.
- **Waiting summed per thread.** Two questions asked at once would count the same wait twice.
- **Cost of delay per day of the window only.** The weekly total would then miss the first week's days before the window, which the throughput week counts.
