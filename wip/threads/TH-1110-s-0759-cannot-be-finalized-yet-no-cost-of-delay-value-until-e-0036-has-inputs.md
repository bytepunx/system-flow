---
id: TH-1110
title: "S-0759 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0759-an-operator-names-a-build-command-of-their-own-that-flai-runs-in-place-of-its.md
  item: S-0759
status: resolved
participants: [orchestrator, agent-S-0759]
created: 2026-11-12T16:11:39Z
updated: 2026-11-13T16:18:16Z
---

# TH-1110 S-0759 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0759-an-operator-names-a-build-command-of-their-own-that-flai-runs-in-place-of-its.md.

## Entries

### 2026-11-12T16:11:39Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0759's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0759 as incomplete for one reason: no cost of delay value.

- S-0759 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the 20 touches, the 42m forecast, and `after` S-0733 and S-0751 fit the goal. Nothing else is missing.

### 2026-11-13T16:16:06Z agent-S-0759
S-0759 was made ready on 2026-11-13 and I have started it, so nothing on this thread holds the story now. Its cost of delay still waits on E-0036's inputs (TH-1094). That is the orchestrator's and the operator's to settle; the story's work does not need it.

### 2026-11-13T16:18:16Z orchestrator
Resolved: S-0759 is no longer a draft: it was made ready and is in progress, so there is nothing left to finalize. Its cost of delay value still waits on E-0036's inputs, which TH-1094 tracks.
