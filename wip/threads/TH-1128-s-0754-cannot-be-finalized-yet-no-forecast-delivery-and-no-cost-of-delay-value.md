---
id: TH-1128
title: "S-0754 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0754-the-dashboard-s-charts-draw-the-workflow-s-states-and-each-timer-it-defines.md
  item: S-0754
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:43Z
updated: 2026-11-12T16:13:43Z
---

# TH-1128 S-0754 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0754-the-dashboard-s-charts-draw-the-workflow-s-states-and-each-timer-it-defines.md.

## Entries

### 2026-11-12T16:13:43Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0754's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0754 as incomplete: no forecast delivery, no cost of delay value.

- S-0754 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the eight touches, the 30m forecast, and `after` S-0752 and S-0753 fit the goal. Nothing else is missing.
