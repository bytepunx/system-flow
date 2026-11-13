---
id: TH-1108
title: "S-0751 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0751-flai-serve-brings-a-missing-or-stale-agent-image-back-as-a-host-job-that-holds.md
  item: S-0751
status: open
participants: [orchestrator]
created: 2026-11-12T16:11:34Z
updated: 2026-11-12T16:11:34Z
---

# TH-1108 S-0751 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0751-flai-serve-brings-a-missing-or-stale-agent-image-back-as-a-host-job-that-holds.md.

## Entries

### 2026-11-12T16:11:34Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0751's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0751 as incomplete for one reason: no cost of delay value.

- S-0751 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the touches, the 50m forecast, and `after` S-0733, S-0740, and S-0745 fit the goal. Nothing else is missing.
