---
id: TH-1105
title: "S-0735 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0735-flai-image-status-says-which-agent-images-a-project-s-agents-need-and-whether.md
  item: S-0735
status: open
participants: [orchestrator]
created: 2026-11-12T16:11:26Z
updated: 2026-11-12T16:11:26Z
---

# TH-1105 S-0735 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0735-flai-image-status-says-which-agent-images-a-project-s-agents-need-and-whether.md.

## Entries

### 2026-11-12T16:11:26Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0735's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0735 as incomplete for one reason: no cost of delay value.

- S-0735 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the 13 touches, the 33m forecast, and `after` S-0733 fit the goal. Nothing else is missing.
