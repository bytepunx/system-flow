---
id: TH-1113
title: "S-0736 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0736-an-adr-fixes-the-workflow-file-s-format-and-place-under-design-workflows-its.md
  item: S-0736
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:07Z
updated: 2026-11-12T16:13:07Z
---

# TH-1113 S-0736 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0736-an-adr-fixes-the-workflow-file-s-format-and-place-under-design-workflows-its.md.

## Entries

### 2026-11-12T16:13:07Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0736's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0736 as incomplete: no forecast delivery, no cost of delay value.

- S-0736 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The nine criteria, the six touches, and the 45m forecast fit the goal. It names no `after`, as the first story of E-0034 should. Nothing else is missing.
