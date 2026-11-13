---
id: TH-1126
title: "S-0752 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0752-flai-stats-reports-each-timer-a-workflow-defines-beside-the-baseline-s-figures.md
  item: S-0752
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:39Z
updated: 2026-11-12T16:13:39Z
---

# TH-1126 S-0752 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0752-flai-stats-reports-each-timer-a-workflow-defines-beside-the-baseline-s-figures.md.

## Entries

### 2026-11-12T16:13:39Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0752's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0752 as incomplete: no forecast delivery, no cost of delay value.

- S-0752 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 11 touches, the 33m forecast, and `after` S-0750 fit the goal. Nothing else is missing.
