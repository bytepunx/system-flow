---
id: TH-1125
title: "S-0750 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0750-flai-stats-computes-lead-cycle-queue-and-review-time-from-the-baseline-workflow.md
  item: S-0750
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:37Z
updated: 2026-11-12T16:13:37Z
---

# TH-1125 S-0750 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0750-flai-stats-computes-lead-cycle-queue-and-review-time-from-the-baseline-workflow.md.

## Entries

### 2026-11-12T16:13:37Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0750's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0750 as incomplete: no forecast delivery, no cost of delay value.

- S-0750 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 20 touches, the 1h forecast, and `after` S-0737 fit the goal. It writes the metrics ADR the project's rule asks for before `metrics.md` changes. Nothing else is missing.
