---
id: TH-1116
title: "S-0739 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md
  item: S-0739
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:14Z
updated: 2026-11-12T16:13:14Z
---

# TH-1116 S-0739 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0739-workitem-s-epic-following-board-view-task-plan-holds-order-and-promotion-read-a.md.

## Entries

### 2026-11-12T16:13:14Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0739's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0739 as incomplete: no forecast delivery, no cost of delay value.

- S-0739 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 19 touches, the 45m forecast, and `after` S-0737 fit the goal. Nothing else is missing.
