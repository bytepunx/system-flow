---
id: TH-1124
title: "S-0749 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0749-the-guard-messages-between-stories-forecasts-prime-packs-the-release-evaluation.md
  item: S-0749
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:34Z
updated: 2026-11-12T16:13:34Z
---

# TH-1124 S-0749 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0749-the-guard-messages-between-stories-forecasts-prime-packs-the-release-evaluation.md.

## Entries

### 2026-11-12T16:13:34Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0749's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0749 as incomplete: no forecast delivery, no cost of delay value.

- S-0749 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 17 touches, the 42m forecast, and `after` S-0737 fit the goal. Nothing else is missing.
