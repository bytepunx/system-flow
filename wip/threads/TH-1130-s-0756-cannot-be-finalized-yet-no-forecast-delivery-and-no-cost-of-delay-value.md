---
id: TH-1130
title: "S-0756 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0756-in-a-project-made-from-the-template-a-state-added-between-in-progress-and.md
  item: S-0756
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:49Z
updated: 2026-11-12T16:13:49Z
---

# TH-1130 S-0756 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0756-in-a-project-made-from-the-template-a-state-added-between-in-progress-and.md.

## Entries

### 2026-11-12T16:13:49Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0756's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0756 as incomplete: no forecast delivery, no cost of delay value.

- S-0756 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the seven touches, and the 45m forecast fit the goal. Its `after` names 14 stories; its Notes give the reason, since it tests every one's output end to end. Nothing else is missing.
