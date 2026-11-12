---
id: TH-1097
title: "E-0034 cost of delay inputs: the epic has none, so no story of it can be valued"
anchor:
  path: wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md
  item: E-0034
status: open
participants: [planner-E-0034, orchestrator]
created: 2026-11-12T15:55:31Z
updated: 2026-11-12T15:59:30Z
---

# TH-1097 E-0034 cost of delay inputs: the epic has none, so no story of it can be valued

On wip/kanban/epics/E-0034-an-operator-defines-the-board-s-states-transitions-rules-and-timers-in-a.md.

## Entries

### 2026-11-12T15:55:31Z planner-E-0034
E-0034 has no cost of delay inputs, and none of its stories has any, so `flai cod` refuses every story I draft for it. The inputs are yours.

Recommendation: `time_lost_per_cycle: 2h` on E-0034. At `planning.hour_rate` 150 and the default 168h cycle, that is 300 USD a week, shared out to the stories by their forecasts. Why: today any change to the board's states, rules, or timers is a flai and flaiover code story touching about 20 packages. The 2h is what that work costs on average each week, counting the board changes the operator asks for (for example a new timer, or a rule between in-progress and review) that wait on this epic.

Alternatives:

- `revenue_per_week` on E-0034, if a project or client is waiting on a custom workflow, as E-0035 carries 1250 USD for kluster.
- `penalty_per_week: 25`, the low figure E-0015 carries, if the epic only adds flexibility and nobody is waiting on it.

Meanwhile I am drafting the stories, their tasks, touches, and forecasts. Once you answer, I will write each story's value from your figure.

### 2026-11-12T15:59:30Z orchestrator (recommendation)
**Recommended:** set the planner's input on E-0034: `time_lost_per_cycle: 2h`. That gives 300 USD a week, which flai shares out to the stories by their forecasts.

```sh
flai edit E-0034 --time-lost-per-cycle 2h
```

**Why this figure:**

- Today every change to the board's states, rules, or timers is a code story across flai and flaiover. 2h a week is the planner's estimate of that work plus the board changes waiting on this epic.
- No project or client is named as waiting on a custom workflow, so no `revenue_per_week` is shown, unlike E-0035's 1250 USD for kluster.

**Alternatives:**

- `penalty_per_week: 25`, as E-0015 carries, if the epic only adds flexibility and nobody waits on it.
- A `revenue_per_week`, if a project or client does wait on a custom workflow.

**Why this needs you:** these are an epic's inputs, not a backlog story's first. ADR-0224 lets me set them only once you decide them on a thread. When you answer or confirm, I set them under `set_cost_of_delay` on my next run, and the planner writes each story's value.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0224-orchestration-permissions-set-cost-of-delay-off-by-default-lets-the.md § Decision
