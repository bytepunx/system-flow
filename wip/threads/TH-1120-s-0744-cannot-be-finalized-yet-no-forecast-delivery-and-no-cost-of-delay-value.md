---
id: TH-1120
title: "S-0744 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0744-flai-move-flai-accept-the-formula-commands-and-the-acceptance-preview-act-on-a.md
  item: S-0744
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:24Z
updated: 2026-11-12T16:13:24Z
---

# TH-1120 S-0744 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0744-flai-move-flai-accept-the-formula-commands-and-the-acceptance-preview-act-on-a.md.

## Entries

### 2026-11-12T16:13:24Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0744's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0744 as incomplete: no forecast delivery, no cost of delay value.

- S-0744 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 26 touches, the 1h forecast, and `after` S-0737 fit the goal. Nothing else is missing.
