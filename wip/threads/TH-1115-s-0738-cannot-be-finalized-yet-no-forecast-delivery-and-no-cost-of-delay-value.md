---
id: TH-1115
title: "S-0738 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0738-the-rules-repo-move-enforces-are-named-rules-from-a-vocabulary-flai-implements.md
  item: S-0738
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:12Z
updated: 2026-11-12T16:13:12Z
---

# TH-1115 S-0738 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0738-the-rules-repo-move-enforces-are-named-rules-from-a-vocabulary-flai-implements.md.

## Entries

### 2026-11-12T16:13:12Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0738's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0738 as incomplete: no forecast delivery, no cost of delay value.

- S-0738 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 10 touches, the 29m forecast, and `after` S-0737 fit the goal. Nothing else is missing.
