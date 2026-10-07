---
id: ADR-0114
title: "flai stats lists the longest waits of the window, each with its story, its thread, and who was awaited"
status: accepted
date: 2026-10-07
supersedes: []
superseded_by: []
refines: [ADR-0081]
---

# ADR-0114 flai stats lists the longest waits of the window, each with its story, its thread, and who was awaited

## Context

S-0215 charts how long agents wait for the operator, on threads and in review, and lists the longest waits under the chart with the story, the thread, and who was awaited. `flai stats --json` reports waiting only as totals: per item (`wait_threads_seconds`, `wait_review_seconds`) and per ISO week of the items completed in it (`waiting.weeks[]`), as [ADR-0081](0081-flai-stats-reports-forecast-error-cost-of-delay-waiting-holds-touches-drift-and.md) and [ADR-0091](0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md) set them. Neither names a thread or who ended a wait. `design/system/metrics.md` defines each value `flai stats` reports and changes only with an ADR.

## Decision

`flai stats --json` carries `waiting.longest[]`: the ten longest waits that have a part in the window, longest first, each with its item, its kind (`thread` or `review`), its thread, when it started and ended, the seconds of it inside the window (and, for a thread, inside the item's `in-progress` intervals), and who was awaited: the author of the entry that ended a thread's wait, or the `by` of the move out of `review`. It covers every item of the report's type that is not cancelled, whatever its state, so a story still waiting is listed.

## Consequences

- The dashboard lists the longest waits from `/api/stats` alone, matching `flai stats --json` to the second.
- A wait is listed on its own, not merged with the item's other waits, so two overlapping waits of one story both appear, and their seconds can add up to more than the story's `wait_threads_seconds`.
- The list's waits are of items in any state, while `waiting.weeks[]` counts only items completed in the window: the table can show a wait that no week's bar holds yet.
- The report grows by at most ten short entries.

## Alternatives considered

- Listing the waits of the completed items only, as the weeks do: a story waiting now, the operator's bottleneck, would be missing until it is accepted.
- Merging each item's waits into one row per item: the thread and who was awaited would be lost.
- Reading the threads in the dashboard: the charts read `/api/stats` alone, and a second source could disagree with the report's window.
