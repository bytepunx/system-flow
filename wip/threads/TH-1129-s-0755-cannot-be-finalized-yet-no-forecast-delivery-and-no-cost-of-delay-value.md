---
id: TH-1129
title: "S-0755 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0755-flai-reads-a-project-s-own-workflow-file-in-place-of-the-baseline-the-template.md
  item: S-0755
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:46Z
updated: 2026-11-12T16:13:46Z
---

# TH-1129 S-0755 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0755-flai-reads-a-project-s-own-workflow-file-in-place-of-the-baseline-the-template.md.

## Entries

### 2026-11-12T16:13:46Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0755's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0755 as incomplete: no forecast delivery, no cost of delay value.

- S-0755 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 17 touches, the 47m forecast, and `after` S-0741 (the workflow findings it runs on the project's file) fit the goal. Nothing else is missing.
