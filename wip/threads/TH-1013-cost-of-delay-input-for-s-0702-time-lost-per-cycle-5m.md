---
id: TH-1013
title: "Cost of delay input for S-0702: time_lost_per_cycle 5m"
anchor:
  path: wip/kanban/stories/S-0702-flai-check-finds-issues-summary-outside-the-story-at-close-out.md
  item: S-0702
status: open
participants: [planner-S-0702]
created: 2026-11-10T10:27:37Z
updated: 2026-11-10T10:27:37Z
---

# TH-1013 Cost of delay input for S-0702: time_lost_per_cycle 5m

On wip/kanban/stories/S-0702-flai-check-finds-issues-summary-outside-the-story-at-close-out.md.

## Entries

### 2026-11-10T10:27:37Z planner-S-0702
S-0702 has no cost of delay inputs and no epic, so `flai cod S-0702` cannot work out a value.

Recommendation: `time_lost_per_cycle: 5m`.

Reason: I-0266 has one instance (S-0576, 2026-11-07). Its cause was an issue file, I-0264, that reached main without its row in `design/issues/summary.md`. Until the summary is regenerated, every close-out records the `issues.summary` note outside its story and commits an issue for it. The close-out does regenerate the summary, so the damage stops after one round. That makes it rarer and cheaper than `agent.image` (S-0582, 10m), which hit every close-out in a container.

Alternatives: 10m if you count each stale summary as hitting several close-outs before it is regenerated; 2m if one instance in three days is the true rate.

Please answer with the input to set. I will work out the value with `flai cod` once it is set.
