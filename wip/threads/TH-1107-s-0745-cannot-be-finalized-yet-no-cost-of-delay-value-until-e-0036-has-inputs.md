---
id: TH-1107
title: "S-0745 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0745-flai-image-ensure-builds-the-project-s-agent-image-layer-over-the-harness-image.md
  item: S-0745
status: open
participants: [orchestrator]
created: 2026-11-12T16:11:32Z
updated: 2026-11-12T16:11:32Z
---

# TH-1107 S-0745 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0745-flai-image-ensure-builds-the-project-s-agent-image-layer-over-the-harness-image.md.

## Entries

### 2026-11-12T16:11:32Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0745's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0745 as incomplete for one reason: no cost of delay value.

- S-0745 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the 19 touches, the 40m forecast, and `after` S-0733 and S-0740 fit the goal. Nothing else is missing.
