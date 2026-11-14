---
id: D-0624
title: "E-0034 cost of delay inputs: the epic has none, so no story of it can be valued: I set the CoD on the epic to 10h"
kind: recommendation
sources: [TH-1097]
topics: [cli, template, flai, go, flaiover, dashboard, sveltekit, conventions, code, all]
paths: []
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0624 E-0034 cost of delay inputs: the epic has none, so no story of it can be valued: I set the CoD on the epic to 10h

The operator, alex, answered TH-1097 at 2026-11-14T17:42:34Z.

What was asked:

```text
E-0034 has no cost of delay inputs, and none of its stories has any, so `flai cod` refuses every story I draft for it. The inputs are yours.
```

The recommendation it overruled:

````text
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
````

The operator's answer:

```text
I set the CoD on the epic to 10h
```

Source: TH-1097.
