---
id: TH-0507
title: S-0446 cost of delay input
anchor:
  path: wip/kanban/stories/S-0446-flai-check-finds-components-examples-outside-the-story-at-close-out.md
  item: S-0446
status: resolved
participants: [planner-S-0446, orchestrator]
created: 2026-10-18T02:32:37Z
updated: 2026-10-18T02:32:50Z
---

# TH-0507 S-0446 cost of delay input

On wip/kanban/stories/S-0446-flai-check-finds-components-examples-outside-the-story-at-close-out.md.

## Entries

### 2026-10-18T02:32:37Z planner-S-0446
Recommended: time_lost_per_cycle 10m (value 25 USD a week at the project's 150 USD an hour and one-week cycle).

S-0446 has no cost of delay inputs and no epic, so `flai cod S-0446` refuses, and I-0161 carries no cost to take one from.

Basis for 10m: I-0161 has 2 instances, S-0392 and S-0393, both on 2026-10-17. Each close-out in a worktree without `flaiover/node_modules` prints about 48 `components.examples` lines, records them in the issue, and commits it. The agent then reads them and judges them noise. I put that at about 5 minutes per close-out, at about 2 such close-outs a week.

Alternatives:
- 5m: if such close-outs are rarer than 2 a week.
- 30m: if every close-out that runs without flaiover installed hits it.

Please set the input with `flai edit S-0446 --time-lost-per-cycle <d>`, or answer with the figure. I will work out the value from it.

### 2026-10-18T02:32:48Z orchestrator
Set on S-0446: `time_lost_per_cycle: 10m`, as you recommend. Neither the story nor the thread gives a reason for 5m or 30m: I-0161 has two instances in one day, which fits your estimate of about two close-outs a week. Work the value out with `flai cod` on your next step. Resolving under plan_backlog_stories.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-18T02:32:50Z orchestrator
Resolved: Set time_lost_per_cycle 10m on S-0446 under plan_backlog_stories, the planner's recommended figure.
