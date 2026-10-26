---
id: TH-0657
title: Cost of delay inputs for E-0029
anchor:
  path: wip/kanban/epics/E-0029-flai-runs-agents-in-docker-containers.md
  item: E-0029
status: open
participants: [planner-E-0029]
created: 2026-10-26T04:40:45Z
updated: 2026-10-26T04:40:45Z
---

# TH-0657 Cost of delay inputs for E-0029

On wip/kanban/epics/E-0029-flai-runs-agents-in-docker-containers.md.

## Entries

### 2026-10-26T04:40:45Z planner-E-0029
E-0029 has no cost of delay inputs, so `flai cod` refuses a value for it and for each story I draft under it (a story without inputs takes a share of its epic's value). The inputs are yours.

Recommended: `penalty_per_week: 150` USD on E-0029. S-0496 and S-0497, the two research stories that led here, each carried 50 USD a week; this epic delivers six stories (the spike, the HTTP MCP server for containers, pilot, the container adapter, the strategic agents in containers, and the project layer), and each week it waits keeps every agent running as the operator with the whole machine in reach.

Alternatives: 50 USD a week, the research stories' figure, if the boundary is not urgent; 300 USD a week if a containerised agent is wanted before the next harness (E-0023) is adopted. `revenue_per_week` and `time_lost_per_cycle` do not fit: the epic brings no revenue and loses no work per cycle.

Set it with `flai edit E-0029 --penalty-per-week 150` or from the epic's page. I am drafting the stories meanwhile and will write each story's value once the input is there.
