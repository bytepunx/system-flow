---
id: TH-1102
title: "S-0733 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0733-an-adr-decides-where-an-installed-flai-s-agent-images-come-from-how-each-is.md
  item: S-0733
status: open
participants: [orchestrator]
created: 2026-11-12T15:59:24Z
updated: 2026-11-12T15:59:24Z
---

# TH-1102 S-0733 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0733-an-adr-decides-where-an-installed-flai-s-agent-images-come-from-how-each-is.md.

## Entries

### 2026-11-12T15:59:24Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I set them under `set_cost_of_delay`, the planner writes S-0733's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0733 as incomplete: no cost of delay value.

- S-0733 has no inputs of its own, and its epic E-0036 has none yet, so `flai cod` cannot work out its share.
- The criteria, the two touches, and the 25m forecast are consistent with the goal. Nothing else is missing.
