---
id: TH-1104
title: "S-0734 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0734-each-flai-release-builds-signs-and-publishes-the-bedrock-and-claude-code-agent.md
  item: S-0734
status: open
participants: [orchestrator]
created: 2026-11-12T16:09:23Z
updated: 2026-11-12T16:09:23Z
---

# TH-1104 S-0734 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0734-each-flai-release-builds-signs-and-publishes-the-bedrock-and-claude-code-agent.md.

## Entries

### 2026-11-12T16:09:23Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I set them under `set_cost_of_delay`, the planner writes S-0734's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0734 as incomplete: no cost of delay value.

- S-0734 has no inputs of its own, and its epic E-0036 has none yet, so `flai cod` cannot work out its share.
- The criteria, the touches, the 45m forecast, and `after` S-0733 are consistent with the goal. Nothing else is missing.
- Its notes say it holds only if S-0733's ADR chooses publication; that is settled when S-0733 runs, not before finalizing.
