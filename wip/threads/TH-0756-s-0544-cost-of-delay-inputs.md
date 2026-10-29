---
id: TH-0756
title: S-0544 cost of delay inputs
anchor:
  path: wip/kanban/stories/S-0544-make-the-dashboard-and-server-versions-easy-to-see.md
  item: S-0544
status: resolved
participants: [planner-S-0544, orchestrator, alex]
created: 2026-10-29T17:41:59Z
updated: 2026-10-29T18:02:08Z
---

# TH-0756 S-0544 cost of delay inputs

On wip/kanban/stories/S-0544-make-the-dashboard-and-server-versions-easy-to-see.md.

## Entries

### 2026-10-29T17:41:59Z planner-S-0544
S-0544 has no cost of delay inputs and no epic, so `flai cod S-0544` refuses it.

Recommended: `time_lost_per_cycle: 10m`, no revenue, no penalty. That is the time you lose per cycle finding out which flai and dashboard versions run, and that an update is out only after a push. At this project's rate it gives a value of 25 USD/week, as S-0520 and S-0543 have.

Please confirm, or give other inputs with `flai edit S-0544 --time-lost-per-cycle`, `--revenue-per-week`, or `--penalty-per-week`. I plan the touches, forecast, and tasks meanwhile.

### 2026-10-29T17:43:08Z orchestrator (recommendation)
**Recommendation:** set `time_lost_per_cycle: 10m` on S-0544 (25 USD a week), as the planner recommends. Give it no revenue and no penalty:

```sh
flai edit S-0544 --time-lost-per-cycle 10m
```

The planner writes the value on its next run.

**Why 10m:** the goal is about time lost. Each cycle you spend a little time finding out which flai and dashboard versions are running, and you learn that an update is out only after a push. Nothing in the story points to revenue or a penalty. The figure matches S-0520 and S-0543, which are small dashboard visibility stories like this one.

**Why you:** S-0544 is in progress. flai refused my edit with "S-0544 is in-progress, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". It refused S-0427 the same way on TH-0474 (D-0316), so this input is yours to set.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-29T18:02:08Z alex
Resolved: S-0544 was accepted
