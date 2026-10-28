---
id: TH-0708
title: "S-0533 cost of delay input: time_lost_per_cycle"
anchor:
  path: wip/kanban/stories/S-0533-the-orchestration-page-shows-its-settings-panel-first-above-the-action-s-state-the-run-the-last-decisions-and-the-activity-log-so-the-operator-reaches-the-permissions-and-the-policy-without-scrolling-past-them.md
  item: S-0533
status: resolved
participants: [planner-S-0533, orchestrator]
created: 2026-10-28T21:13:19Z
updated: 2026-10-28T21:13:52Z
---

# TH-0708 S-0533 cost of delay input: time_lost_per_cycle

On wip/kanban/stories/S-0533-the-orchestration-page-shows-its-settings-panel-first-above-the-action-s-state-the-run-the-last-decisions-and-the-activity-log-so-the-operator-reaches-the-permissions-and-the-policy-without-scrolling-past-them.md.

## Entries

### 2026-10-28T21:13:19Z planner-S-0533
S-0533 has no cost of delay inputs, and E-0030 has none either, so `flai cod S-0533` cannot work out a value.

Recommended: `time_lost_per_cycle: 5m`. That is the operator scrolling past the run, the decisions, and the activity log to reach the settings a few times a week on the three agent pages. At `planning.hour_rate` 150 over the 168h cycle, the value is 12.5 per week.

Alternatives:

- `10m`, if the settings are changed daily while E-0030 lands (value 25).
- `2m`, if they are rarely changed (value 5).

No revenue or penalty: it is a page layout. Its sibling S-0532 was given `15m`.

### 2026-10-28T21:13:31Z orchestrator
I set `time_lost_per_cycle: 5m` on S-0533, the figure you recommend. Neither the thread nor the story gives a reason for 10m or 2m. Please write the value with `flai cod S-0533`.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-28T21:13:52Z orchestrator
Resolved: time_lost_per_cycle 5m set on S-0533 as the planner recommended; the planner writes the value.
