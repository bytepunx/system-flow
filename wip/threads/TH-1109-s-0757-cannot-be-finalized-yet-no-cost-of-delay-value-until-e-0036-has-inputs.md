---
id: TH-1109
title: "S-0757 cannot be finalized yet: no cost of delay value until E-0036 has inputs"
anchor:
  path: wip/kanban/stories/S-0757-the-dashboard-shows-each-agent-image-job-as-it-runs-its-output-and-a-failure-s.md
  item: S-0757
status: open
participants: [orchestrator]
created: 2026-11-12T16:11:37Z
updated: 2026-11-12T16:11:37Z
---

# TH-1109 S-0757 cannot be finalized yet: no cost of delay value until E-0036 has inputs

On wip/kanban/stories/S-0757-the-dashboard-shows-each-agent-image-job-as-it-runs-its-output-and-a-failure-s.md.

## Entries

### 2026-11-12T16:11:37Z orchestrator
**Recommended:** answer or confirm my recommendation on TH-1094, the inputs for E-0036. I then set them under `set_cost_of_delay`, the planner writes S-0757's value, and I finalize the draft once `flai promote --drafts` lists it complete.

`flai promote --drafts` lists S-0757 as incomplete for one reason: no cost of delay value.

- S-0757 has no inputs of its own. Its epic, E-0036, has none yet either, so `flai cod` cannot work out its share.
- The criteria, the 17 touches, the 50m forecast, and `after` S-0751 fit the goal. Nothing else is missing.
