---
id: TH-1127
title: "S-0753 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0753-the-dashboard-s-board-draws-the-columns-their-kinds-and-the-moves-between-them.md
  item: S-0753
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:41Z
updated: 2026-11-12T16:13:41Z
---

# TH-1127 S-0753 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0753-the-dashboard-s-board-draws-the-columns-their-kinds-and-the-moves-between-them.md.

## Entries

### 2026-11-12T16:13:41Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0753's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0753 as incomplete: no forecast delivery, no cost of delay value.

- S-0753 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 17 touches, the 42m forecast, and `after` S-0739 (the board view's columns) fit the goal. Nothing else is missing.
