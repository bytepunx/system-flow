---
id: TH-1114
title: "S-0737 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0737-flai-reads-the-baseline-workflow-it-carries-and-the-board-s-states-and.md
  item: S-0737
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:10Z
updated: 2026-11-12T16:13:10Z
---

# TH-1114 S-0737 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0737-flai-reads-the-baseline-workflow-it-carries-and-the-board-s-states-and.md.

## Entries

### 2026-11-12T16:13:10Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0737's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0737 as incomplete: no forecast delivery, no cost of delay value.

- S-0737 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 11 touches, the 33m forecast, and `after` S-0736 fit the goal.

**One inconsistency, minor:** its Notes say the template ships the baseline workflow file in "S-0750's work". S-0755 does that, not S-0750 (S-0755's second criterion). The planner can correct the reference on its next run. It does not stop the draft being finalized.
