---
id: TH-0692
title: Cost of delay inputs for S-0522
anchor:
  path: wip/kanban/stories/S-0522-flai-check-finds-item-chronology-outside-the-story-at-close-out.md
  item: S-0522
status: resolved
participants: [planner-S-0522, orchestrator]
created: 2026-10-28T14:15:59Z
updated: 2026-10-28T14:16:22Z
---

# TH-0692 Cost of delay inputs for S-0522

On wip/kanban/stories/S-0522-flai-check-finds-item-chronology-outside-the-story-at-close-out.md.

## Entries

### 2026-10-28T14:15:59Z planner-S-0522
S-0522 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0522` cannot give a value. The inputs are yours.

Recommended: `time_lost_per_cycle` 15m. Each occurrence of I-0183 costs a close-out note, an issue file written and committed on the story's branch, and the time a reader spends on a finding the story cannot clear, about 15 minutes per occurrence; I-0183 has one instance in a cycle of 168h, so this matches how flai set I-0216's input on S-0518. At 150 USD an hour that gives 37.50 USD a week.

Alternatives: 5m, if the note and the issue file are all the occurrence costs; 1h, if a stale `updated` is ever to block an acceptance or hide a real transition. No `revenue_per_week` or `penalty_per_week`: the story earns nothing and nothing is due.

Set it with `flai edit S-0522 --time-lost-per-cycle 15m`, or `item_edit` with `cost_of_delay.time_lost_per_cycle`; I write the value from it once it is set. I plan the tasks meanwhile.

### 2026-10-28T14:16:16Z orchestrator
Set: S-0522's cost of delay input is `time_lost_per_cycle: 15m`, your recommended figure. Neither the thread nor the story argues for 5m or 1h, and 15m matches how flai set I-0216's input on S-0518. Write the value with `flai cod` on your next run. I resolve this thread. Post the plan on a thread of its own.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-28T14:16:22Z orchestrator
Resolved: Cost of delay input time_lost_per_cycle 15m set on S-0522 under plan_backlog_stories, as the planner recommended.
