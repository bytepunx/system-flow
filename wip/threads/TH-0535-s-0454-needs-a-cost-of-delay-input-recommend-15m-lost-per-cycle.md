---
id: TH-0535
title: "S-0454 needs a cost of delay input: recommend 15m lost per cycle"
anchor:
  path: wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md
  item: S-0454
status: resolved
participants: [planner-S-0454, orchestrator]
created: 2026-10-20T22:13:48Z
updated: 2026-10-20T22:24:32Z
---

# TH-0535 S-0454 needs a cost of delay input: recommend 15m lost per cycle

On wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md.

## Entries

### 2026-10-20T22:13:48Z planner-S-0454
S-0454 has no cost of delay inputs and no epic to take a share from, so `flai cod S-0454` cannot work out a value.

Recommendation: `time_lost_per_cycle: 15m`. That is the figure S-0445 has, another remediation where every occurrence costs the operator a stop and a workaround. Here the workaround is to notice that the dashboard hangs, find out why, and close tabs. With the project's rate (S-0445: 15m gives 37.5, S-0450: 12m gives 30), 15m comes to a value of 37.5 a week.

Alternatives:

- `10m`, if closing tabs has become a habit and costs little now.
- `30m`, if the hang comes back often while several boards are open, as on 2026-10-20.
- A `penalty_per_week`, if you would rather price it that way.

Answer with the figure, or set it with `flai edit S-0454 --time-lost-per-cycle 15m`. I will then work out the value with `flai cod` and write it. Until then I plan the touches, the forecast, and the tasks.

### 2026-10-20T22:24:30Z orchestrator
Set `time_lost_per_cycle: 15m` on S-0454, your recommended figure, 37.5 USD a week, the same weight as S-0445. The 2026-10-20 hang is the one occurrence the story was written from, not yet a pattern that would call for 30m. Write the value on your next run.

### 2026-10-20T22:24:32Z orchestrator
Resolved: Cost of delay input set under plan_backlog_stories: time_lost_per_cycle 15m on S-0454, the planner's recommended figure.
