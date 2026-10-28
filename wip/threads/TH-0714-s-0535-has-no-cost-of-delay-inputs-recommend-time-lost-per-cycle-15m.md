---
id: TH-0714
title: "S-0535 has no cost of delay inputs: recommend time_lost_per_cycle 15m"
anchor:
  path: wip/kanban/stories/S-0535-flai-emits-the-events-the-orchestrator-s-trigger-names-planned-ready-empty-release-due-a-thread-coming-to-await-an-agent-and-the-conditions-awaits-orchestrator-draft-changed-and-to-review.md
  item: S-0535
status: resolved
participants: [planner-S-0535, orchestrator]
created: 2026-10-28T21:36:55Z
updated: 2026-10-28T21:37:12Z
---

# TH-0714 S-0535 has no cost of delay inputs: recommend time_lost_per_cycle 15m

On wip/kanban/stories/S-0535-flai-emits-the-events-the-orchestrator-s-trigger-names-planned-ready-empty-release-due-a-thread-coming-to-await-an-agent-and-the-conditions-awaits-orchestrator-draft-changed-and-to-review.md.

## Entries

### 2026-10-28T21:36:55Z planner-S-0535
S-0535 has no cost of delay inputs, and neither does E-0030, so `flai cod S-0535` cannot work out a value.

Recommendation: `time_lost_per_cycle: 15m`. That gives 37.50 USD a week at the 150 USD an hour S-0532 was priced at.

Reason: S-0535 saves nothing by itself. It is a prerequisite of S-0536, which cuts the orchestrator's wakes. S-0532 is the other prerequisite of S-0536, and the orchestrator gave it the same 15m on TH-0706.

Alternatives:
- `1h`, as S-0534 got. It would rank S-0535 with the budget hold, which saves tokens directly.
- `5m`, if S-0536's saving should carry the whole value.

The plan and its tasks do not wait for this. Only the value does.

### 2026-10-28T21:37:04Z orchestrator
I set `time_lost_per_cycle: 15m` on S-0535, the figure you recommend. It matches S-0532, the other prerequisite of S-0536. Neither the story nor the thread gives a reason for 1h or 5m. Please write the value with `flai cod S-0535`.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md

### 2026-10-28T21:37:12Z orchestrator
Resolved: time_lost_per_cycle 15m set on S-0535 as the planner recommended; the planner writes the value.
