---
id: TH-1109
title: "S-0757 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0757-the-dashboard-shows-each-agent-image-job-as-it-runs-its-output-and-a-failure-s.md
  item: S-0757
status: resolved
participants: [orchestrator, agent-S-0757]
created: 2026-11-12T16:11:37Z
updated: 2026-11-13T16:48:46Z
---

# TH-1109 S-0757 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0757-the-dashboard-shows-each-agent-image-job-as-it-runs-its-output-and-a-failure-s.md.

## Entries

### 2026-11-12T16:11:37Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0757's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0757 as incomplete for one reason: no cost of delay value.

- S-0757 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the 17 touches, the 50m forecast, and `after` S-0751 fit the goal. Nothing else is missing.

### 2026-11-13T16:47:51Z agent-S-0757
Nothing for S-0757's agent to do here. The operator finalized S-0757, moved it to ready on 2026-11-13, and its agent is working it now. Its cost of delay still waits on E-0036's inputs, which are the operator's to give on TH-1094. Once they are set, the planner can write the value whenever the story is in progress. This thread can be resolved by its opener.

### 2026-11-13T16:48:46Z orchestrator
Resolved: The concern is gone: the operator finalized S-0757 and moved it to ready on 2026-11-13, and it is in progress now. Its cost of delay still waits on E-0036's inputs, which TH-1094 asks the operator for.
