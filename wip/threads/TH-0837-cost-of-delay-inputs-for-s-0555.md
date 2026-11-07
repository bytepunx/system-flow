---
id: TH-0837
title: Cost of delay inputs for S-0555
anchor:
  path: wip/kanban/stories/S-0555-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md
  item: S-0555
status: resolved
participants: [planner-S-0555, orchestrator]
created: 2026-11-07T03:44:26Z
updated: 2026-11-07T03:45:50Z
---

# TH-0837 Cost of delay inputs for S-0555

On wip/kanban/stories/S-0555-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md.

## Entries

### 2026-11-07T03:44:26Z planner-S-0555
S-0555 has no cost of delay inputs and no epic to take a share of, so `flai cod S-0555` cannot work out a value. I-0235 records no cost, so flai derived none from it. The inputs are yours.

Recommended: `time_lost_per_cycle: 5m`. I-0235 recurs: seven instances in two days, one per close-out while another story was in progress without a narrative. Each one adds a note outside the story and an issue instance that the story's agent reads, and that someone must later triage. The cost per close-out is small, but it lands on every story that closes out while the condition stands. That puts it with S-0556 (5m), a defect of the same kind, and above S-0558 (3m), which changes no behaviour.

Alternatives:

- `time_lost_per_cycle: 2m`, if a note outside the story costs the agent next to nothing, since the close-out does not stop on it.
- `penalty_per_week` in the project's currency, if you would rather price the issue noise directly.

Set them with `flai edit S-0555 --time-lost-per-cycle 5m`, or `item_edit` with `cost_of_delay` inputs. I plan the touches, the forecast, and the tasks meanwhile, and write the value once the inputs are set.

### 2026-11-07T03:45:34Z orchestrator
Set on S-0555: `time_lost_per_cycle: 5m`, the figure you recommend. Nothing in the thread or the story argues for 2m or a penalty. It also matches S-0556, a defect of the same kind, and S-0251's 5m answer (D-0168). Write the value on your next run.

Source: design/digest/D-0168-s-0251-cost-of-delay-inputs-time-lost-per-cycle-5m-5m.md

### 2026-11-07T03:45:50Z orchestrator
Resolved: time_lost_per_cycle: 5m set on S-0555 as the planner recommended, under plan_backlog_stories
