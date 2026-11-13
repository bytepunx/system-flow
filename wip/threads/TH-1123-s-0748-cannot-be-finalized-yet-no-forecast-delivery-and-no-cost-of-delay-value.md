---
id: TH-1123
title: "S-0748 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0748-a-role-file-s-to-review-and-to-backlog-conditions-become-one-condition-that.md
  item: S-0748
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:32Z
updated: 2026-11-12T16:13:32Z
---

# TH-1123 S-0748 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0748-a-role-file-s-to-review-and-to-backlog-conditions-become-one-condition-that.md.

## Entries

### 2026-11-12T16:13:32Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0748's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0748 as incomplete: no forecast delivery, no cost of delay value.

- S-0748 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 20 touches, the 53m forecast, and `after` S-0737 fit the goal.

**One thing to know before it runs:** it changes the baseline role files, `flai/internal/roles/baseline/` and `design/workflows/roles/`, which are code you run (ADR-0207 §7). Its Notes say so, and its acceptance is yours to weigh for that. It does not stop the draft being finalized.
