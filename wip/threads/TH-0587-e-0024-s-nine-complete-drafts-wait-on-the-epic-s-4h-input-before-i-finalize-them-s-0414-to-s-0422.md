---
id: TH-0587
title: "E-0024's nine complete drafts wait on the epic's 4h input before I finalize them: S-0414 to S-0422"
anchor:
  path: wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md
  item: E-0024
status: open
participants: [orchestrator]
created: 2026-10-21T16:02:39Z
updated: 2026-10-21T16:02:39Z
---

# TH-0587 E-0024's nine complete drafts wait on the epic's 4h input before I finalize them: S-0414 to S-0422

On wip/kanban/epics/E-0024-workflows-as-data-a-project-defines-a-strategic-role-or-a-repeatable-formula-of-steps-in-files-flai-runs-and-the-built-in-roles-are-expressed-the-same-way-so-a-new-role-or-a-release-process-needs-no-flai-release.md.

## Entries

### 2026-10-21T16:02:39Z orchestrator
**Recommended:** run `flai edit E-0024 --time-lost-per-cycle 4h`. On the planner's next run, `flai cod` then values each story, and I finalize S-0414 to S-0422 under finalize_drafts.

**The alternative:** say "finalize as they are", and I will finalize them with the planners' values.

## Why I have not finalized them

`flai promote --drafts` lists S-0414 to S-0422 as complete. Their criteria, touches, and forecasts fit the plans I approved. Their cost of delay values do not come from flai:

- The planners worked each value out by hand, from the 4h you confirmed on TH-0567 and TH-0473, because `flai cod` refuses while E-0024 has no input.
- flai guard refuses to let me set the epic's input (I-0187).
- The new story after S-0419, which you confirmed on TH-0579, will change each story's share of the epic once it exists (I-0188, `flai plan E-0024`).

The order is by throughput, so these values do not change which story goes first. They do feed the release figures.

## Not covered by this thread

- S-0423 waits on your answers on TH-0583.
- S-0424 and S-0425 have no value yet. S-0424's planner is running now.
