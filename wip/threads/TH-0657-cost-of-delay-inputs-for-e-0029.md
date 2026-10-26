---
id: TH-0657
title: Cost of delay inputs for E-0029
anchor:
  path: wip/kanban/epics/E-0029-flai-runs-agents-in-docker-containers.md
  item: E-0029
status: resolved
participants: [planner-E-0029, orchestrator, alex]
created: 2026-10-26T04:40:45Z
updated: 2026-10-26T12:28:19Z
---

# TH-0657 Cost of delay inputs for E-0029

On wip/kanban/epics/E-0029-flai-runs-agents-in-docker-containers.md.

## Entries

### 2026-10-26T04:40:45Z planner-E-0029
E-0029 has no cost of delay inputs, so `flai cod` refuses a value for it and for each story I draft under it (a story without inputs takes a share of its epic's value). The inputs are yours.

Recommended: `penalty_per_week: 150` USD on E-0029. S-0496 and S-0497, the two research stories that led here, each carried 50 USD a week; this epic delivers six stories (the spike, the HTTP MCP server for containers, pilot, the container adapter, the strategic agents in containers, and the project layer), and each week it waits keeps every agent running as the operator with the whole machine in reach.

Alternatives: 50 USD a week, the research stories' figure, if the boundary is not urgent; 300 USD a week if a containerised agent is wanted before the next harness (E-0023) is adopted. `revenue_per_week` and `time_lost_per_cycle` do not fit: the epic brings no revenue and loses no work per cycle.

Set it with `flai edit E-0029 --penalty-per-week 150` or from the epic's page. I am drafting the stories meanwhile and will write each story's value once the input is there.

### 2026-10-26T04:58:23Z orchestrator (recommendation)
**Recommended for the operator:** take the planner's figure, `penalty_per_week: 150` USD on E-0029.

```sh
flai edit E-0029 --penalty-per-week 150
```

- **Why 150.** It is three times the 50 USD/week you gave S-0496 and S-0497, the research behind this epic, for a six-story delivery. Nothing on the thread or the epic argues for 50 (not urgent) or 300 (needed before E-0023).
- **Why it is yours.** An epic's input is money and the operator's, as on TH-0467 (D-0330). Under plan_backlog_stories I give inputs only to a backlog story, never to an epic. Each story's value is the epic's input shared by duration, so setting it once here values all six.
- **What follows.** The planner writes each story's value on its next run. I then judge S-0505 to S-0510 for finalizing under finalize_drafts.

Source: design/digest/D-0330-s-0407-plan-five-tasks-in-four-layers-two-touches-added-needs-a-cost-of-delay.md

### 2026-10-26T12:16:13Z alex
Confirmed the recommendation of 2026-10-26T04:58:23Z orchestrator.

Source: design/digest/D-0330-s-0407-plan-five-tasks-in-four-layers-two-touches-added-needs-a-cost-of-delay.md

### 2026-10-26T12:16:37Z orchestrator
Your confirmation is recorded, but E-0029 still has no cost of delay input. Confirming a recommendation does not write the figure it names, and flai guard does not let me set an epic's input (I-0187). One command does it:

```sh
flai edit E-0029 --penalty-per-week 150
```

After that, E-0029's planner can write each story's value, S-0505 to S-0510. I then judge them for finalizing under finalize_drafts, once TH-0659's plan is settled.

Source: design/conventions/strategic-agents.md § Asking the operator

### 2026-10-26T12:26:17Z alex
Resolved.

### 2026-10-26T12:27:56Z alex
TH-0659 is settled

### 2026-10-26T12:28:19Z alex
Resolved.
