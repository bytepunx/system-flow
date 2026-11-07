---
id: TH-0884
title: "Plan for S-0596: the agent's five tasks stand; cost of delay inputs needed"
anchor:
  path: wip/kanban/stories/S-0596-explore-where-we-can-improve-flai-serve-performance.md
  item: S-0596
status: resolved
participants: [planner-S-0596, orchestrator, alex]
created: 2026-11-07T22:13:59Z
updated: 2026-11-07T23:38:49Z
---

# TH-0884 Plan for S-0596: the agent's five tasks stand; cost of delay inputs needed

On wip/kanban/stories/S-0596-explore-where-we-can-improve-flai-serve-performance.md.

## Entries

### 2026-11-07T22:13:59Z planner-S-0596
One question for you, then the plan.

## Question: the cost of delay inputs

S-0596 has no cost of delay inputs and no epic, so `flai cod` cannot work out a value.

Recommended: `time_lost_per_cycle` 2h. That is the operator's time lost each week to slow board loads, item writes, charts, and close-outs and acceptances, at the 150 USD an hour flai prices lost time at. It gives about 300 USD a week.

Alternatives:

- 1h a cycle, about 150 USD a week, if the slowness is an irritation rather than lost work.
- 4h a cycle, about 600 USD a week, if acceptances of 17 to 42 s and writes of 8 to 32 s (the serve log figures T-2608 cites) stall the operator many times a day.

Set it with `flai edit S-0596 --time-lost-per-cycle 2h`, or answer here. I will then work out the value.

## Plan

The story's agent wrote the tasks after it started the story. I created no tasks. I re-enriched the five it wrote with the tag `cli`.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2607 Explore and measure the dashboard's reads (agent.status, board.get, settings.get, threads, messages, item.get, and whether flai serve answers channel requests one at a time) | none | none |
| 1 | T-2608 Explore and measure the writes (item.new, item.edit, item.finalize, item.move, issue.story, accept.run, the close-out job) | none | none |
| 1 | T-2609 Explore and measure the chart data, stats.get over 7, 30, 90 days and all history | none | none |
| 2 | T-2610 Record the findings in `design/system/server-performance.md` | `design/system/server-performance.md` | T-2607, T-2608, T-2609 |
| 3 | T-2611 Draft the performance epic and one draft story per actionable cause, linked from that document | `design/system/server-performance.md` | T-2610 |

The three explorations write nothing and run at once. Between them they cover the four slow operations the goal names. T-2610 covers criteria 1 and 2. T-2611 covers criterion 3.

Story enrichment:

- Touches: `design/system/server-performance.md`. The agent narrowed your `flai/cmd` to it, since no task changes a command.
- Topics: added `server-side` and `back-end`.
- Forecast: 1h, raised from flai's 36m by comparison with S-0152 (1101 s for the board's and MCP's reads alone). Delivery 2026-11-07T23:10:00Z.
- After: none.

## Assumptions

- The run's start prompt said you had written on a thread since an earlier plan. Neither exists on S-0596, so I treated this as the first plan.
- E-0012 is done and archived, so T-2611 makes a new epic rather than reopening it.
- A cause that "requires a real change to how we track state" is one that changes the files under `wip/`, `design/`, or their schema. A cache or index under `.flai-cache` that is rebuilt from those files is not such a change.

No split, merge, or drop proposed.

### 2026-11-07T22:16:12Z orchestrator (recommendation)
**Recommended:** set `time_lost_per_cycle: 2h` on S-0596, about 300 USD a week, as the planner recommends.

```bash
flai edit S-0596 --time-lost-per-cycle 2h
```

## The figure

- **The goal is time lost.** It names slow board loads, item writes, charts, and close-outs and acceptances. That is the operator's time, not revenue or a penalty.
- **Neither the thread nor the story gives a reason for 1h or 4h.** The 17 to 42 s acceptances are what T-2608 will measure, not a measured count of stalls a day.
- **Put it on the story.** S-0596 has no epic.

## Why this needs you

S-0596 is in progress (agent-S-0596 started it at 2026-11-07T22:06:54Z). ADR-0119 §4 lets me set a cost of delay on a backlog story only, and flai refused me before on an in-progress story (D-0317, D-0441). If you confirm this recommendation, I set the input under `set_cost_of_delay`. If you answer with another figure, I set that one. The planner then works out the value.

## The plan: approved

- **The tasks fit the criteria.** T-2607, T-2608, and T-2609 cover the four slow operations the goal names. T-2607 also asks whether flai serve answers channel requests one at a time. T-2610 records criteria 1 and 2 in `design/system/server-performance.md`. T-2611 drafts the epic and stories for criterion 3.
- **The layers fit.** The three explorations write nothing and run at once. T-2610 waits for them, and T-2611 waits for T-2610.
- **The touches fit.** Narrowing `flai/cmd` to `design/system/server-performance.md` is right: no task changes a command.
- **The forecast fits.** 1h over flai's 36m, with its reason stated (S-0152's reads alone took 1101 s).
- **The assumptions hold.** E-0012 is done, so a new epic is right. Reading "a real change to how we track state" as a change to the files under `wip/` or `design/`, or their schema, keeps a rebuildable cache under `.flai-cache` in scope.

The thread stays open for the cost of delay input.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-11-07T23:38:49Z alex
Resolved: S-0596 was accepted
