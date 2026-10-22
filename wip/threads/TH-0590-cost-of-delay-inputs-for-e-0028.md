---
id: TH-0590
title: Cost of delay inputs for E-0028
anchor:
  path: wip/kanban/epics/E-0028-charts-improvements.md
  item: E-0028
status: resolved
participants: [planner-E-0028, orchestrator]
created: 2026-10-21T20:53:46Z
updated: 2026-10-22T21:15:43Z
---

# TH-0590 Cost of delay inputs for E-0028

On wip/kanban/epics/E-0028-charts-improvements.md.

## Entries

### 2026-10-21T20:53:46Z planner-E-0028
E-0028 has no cost of delay inputs, and `flai cod` refuses a story whose epic has none, so its stories get no value until you set them. The inputs are yours; I will not write them.

Recommended: `time_lost_per_cycle: 1h` on E-0028. Each week the charts stay as they are costs about an hour of your reading time, at `planning.hour_rate` 150 over the default 168h cycle that is 150 USD a week, which the two stories then share by their forecast durations.

Alternatives: `penalty_per_week: 500`, as E-0023 carries, if you weigh the charts higher; or `penalty_per_week: 50` if this is a convenience.

To set it: `flai edit E-0028 --time-lost-per-cycle 1h` (or `--penalty-per-week <n>`), or the cost of delay panel on the epic's page. I plan the stories and tasks meanwhile and write each story's value once the input is there.

### 2026-10-21T20:53:56Z orchestrator (recommendation)
Recommended for the operator: take the planner's figure. Run `flai edit E-0028 --time-lost-per-cycle 1h`, which works out to 150 USD a week at `planning.hour_rate` 150 over the 168h cycle.

- **Why 1h fits:** E-0028 reworks how the charts are reached and explained, as E-0027 reworked the main layout, and you gave E-0027 2h. This is a narrower change, to one page, so half that is in proportion.
- **Choose 500 USD a week** if you read the charts daily to steer the work, so that finding a chart costs you more than an hour a week.
- **Choose 50 USD a week** if this is polish.

The input is yours to set: flai guard refuses me an epic's input (I-0187). Once it is set, `flai cod` gives each E-0028 story its share by forecast duration, and the planner writes those values.

Source: design/system/strategic-agents.md § Cost of delay

### 2026-10-22T20:18:07Z planner-E-0028
Both stories were finalized and S-0484 started on 2026-10-22T20:17Z with E-0028 still carrying no cost of delay inputs, so S-0484 and S-0485 have no value. When you set an input on E-0028 (`flai edit E-0028 --time-lost-per-cycle 1h`, or the panel on its page), `flai cod S-0485` gives the value from the forecasts, and a replan of S-0485 writes it; S-0484 is in progress and takes its value the same way if you ask for one. This run ends here.

### 2026-10-22T21:15:43Z orchestrator
Resolved: S-0485 was accepted
