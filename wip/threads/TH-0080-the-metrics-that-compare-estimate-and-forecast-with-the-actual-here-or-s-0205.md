---
id: TH-0080
title: "The metrics that compare estimate and forecast with the actual: here or S-0205?"
anchor:
  path: wip/kanban/stories/S-0199-work-items-carry-planning-data-draft-cost-of-delay-inputs-and-value-and-a-forecast-with-who-set-it.md
  item: S-0199
status: resolved
participants: [agent-S-0199, alex]
created: 2026-10-03T06:13:37Z
updated: 2026-10-03T06:48:57Z
---

# TH-0080 The metrics that compare estimate and forecast with the actual: here or S-0205?

On wip/kanban/stories/S-0199-work-items-carry-planning-data-draft-cost-of-delay-inputs-and-value-and-a-forecast-with-who-set-it.md.

## Entries

### 2026-10-03T06:13:37Z agent-S-0199
S-0199's third criterion ends "and the metrics compare both to the actual". S-0205's first criterion defines those metrics: `forecast_error`, `delivery_error`, and `estimate_error`, with p50 and p85, its own ADR, and its own metrics.md change. S-0205 already waits for S-0199 (`after`).

Recommended: leave the comparison to S-0205. S-0199 records in ADR-0074 that `estimate` stays the human's figure and `forecast` the planner's, and that the metrics compare each with the actual (S-0205). I would check the criterion with a note naming S-0205. Building the metrics here as well would change metrics.md twice and take S-0205's first criterion.

The alternative is to add `forecast_error` and `estimate_error` per story to `flai stats` here, and leave the aggregates and `delivery_error` to S-0205.

Which do you want? I'm going on with the docs and the dashboard task meanwhile.

### 2026-10-03T06:20:12Z alex
use recommendation

### 2026-10-03T06:48:57Z agent-S-0199
Resolved: The designer chose the recommendation: S-0205 builds the metrics that compare estimate and forecast with the actual. S-0199's third criterion is checked with a note naming S-0205, and ADR-0074 says so.
