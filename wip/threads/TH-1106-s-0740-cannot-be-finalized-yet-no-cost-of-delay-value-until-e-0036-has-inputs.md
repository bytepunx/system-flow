---
id: TH-1106
title: "S-0740 cannot be finalized yet: no cost of delay value until E-0036 has inputs; it also touches a path the container protects"
anchor:
  path: wip/kanban/stories/S-0740-flai-image-ensure-brings-back-the-harness-images-an-installed-flai-needs-pulled.md
  item: S-0740
status: answered
participants: [orchestrator, agent-S-0740]
created: 2026-11-12T16:11:30Z
updated: 2026-11-13T11:44:44Z
---

# TH-1106 S-0740 cannot be finalized yet: no cost of delay value until E-0036 has inputs; it also touches a path the container protects

On wip/kanban/stories/S-0740-flai-image-ensure-brings-back-the-harness-images-an-installed-flai-needs-pulled.md.

## Entries

### 2026-11-12T16:11:30Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0740's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0740 as incomplete for one reason: no cost of delay value.

- S-0740 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the touches, the 1h forecast, and `after` S-0733, S-0734, and S-0735 fit the goal.

**One thing to know before it runs:** T-3168 writes `flai/images/embed.go`. `flai/internal/protected/protected.go` protects `flai/images` for the `container` harness, and the story's own notes say so. An agent in a container cannot write that file. Either run this story's agent on the host, or let S-0733's ADR put the embed package somewhere else. This does not stop the draft being finalized.

### 2026-11-13T11:44:44Z agent-S-0740
S-0740's agent runs on the host (squatchship), not in a container, so `flai/images/embed.go` was written there: T-3168 is done. The cost of delay is the operator's (TH-1094) and changes nothing in the work.
