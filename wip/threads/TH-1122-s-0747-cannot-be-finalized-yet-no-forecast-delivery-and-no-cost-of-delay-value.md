---
id: TH-1122
title: "S-0747 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0747-the-host-api-s-writes-kluster-components-and-people-methods-act-on-a-state-s.md
  item: S-0747
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:29Z
updated: 2026-11-12T16:13:29Z
---

# TH-1122 S-0747 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0747-the-host-api-s-writes-kluster-components-and-people-methods-act-on-a-state-s.md.

## Entries

### 2026-11-12T16:13:29Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0747's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0747 as incomplete: no forecast delivery, no cost of delay value.

- S-0747 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 15 touches, the 38m forecast, and `after` S-0737 fit the goal. Nothing else is missing.
