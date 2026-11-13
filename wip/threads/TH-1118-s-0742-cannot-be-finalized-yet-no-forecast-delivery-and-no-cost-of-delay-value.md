---
id: TH-1118
title: "S-0742 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs"
anchor:
  path: wip/kanban/stories/S-0742-flai-serve-starts-restarts-and-stops-agents-plans-again-and-finds-its-reasons.md
  item: S-0742
status: open
participants: [orchestrator]
created: 2026-11-12T16:13:20Z
updated: 2026-11-12T16:13:20Z
---

# TH-1118 S-0742 cannot be finalized yet: no forecast delivery and no cost of delay value until E-0034 has inputs

On wip/kanban/stories/S-0742-flai-serve-starts-restarts-and-stops-agents-plans-again-and-finds-its-reasons.md.

## Entries

### 2026-11-12T16:13:20Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1097, the inputs for E-0034. I then set them under `set_cost_of_delay`, the planner's next run writes S-0742's value and its forecast delivery, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0742 as incomplete: no forecast delivery, no cost of delay value.

- S-0742 has no inputs of its own, and E-0034 has none yet, so `flai cod` cannot work out its share.
- The criteria, the 30 touches, the 1h20m forecast, and `after` S-0737 fit the goal. Its touches share `serve/orchestrate_reasons.go` and `serve/replan.go` with S-0748, which its Notes name; the trial merge judges that overlap. Nothing else is missing.
