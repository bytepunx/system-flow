---
id: TH-1111
title: "S-0760 cannot be finalized yet: no cost of delay value until E-0036 has inputs; it also touches paths the container protects"
anchor:
  path: wip/kanban/stories/S-0760-this-repository-builds-its-agent-images-through-flai-s-general-path-with.md
  item: S-0760
status: open
participants: [orchestrator]
created: 2026-11-12T16:11:42Z
updated: 2026-11-12T16:11:42Z
---

# TH-1111 S-0760 cannot be finalized yet: no cost of delay value until E-0036 has inputs; it also touches paths the container protects

On wip/kanban/stories/S-0760-this-repository-builds-its-agent-images-through-flai-s-general-path-with.md.

## Entries

### 2026-11-12T16:11:42Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0760's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0760 as incomplete for one reason: no cost of delay value.

- S-0760 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the touches, the 32m forecast, and `after` S-0745 and S-0759 fit the goal.

**One thing to know before it runs:** T-3207 and T-3209 change `flai/images/system-flow/Dockerfile` and `scripts/agent-image.sh`. `flai/internal/protected/protected.go` protects both for the `container` harness, so an agent in a container cannot write them. Run this story's agent on the host, or expect to make those two changes yourself. This does not stop the draft being finalized.
