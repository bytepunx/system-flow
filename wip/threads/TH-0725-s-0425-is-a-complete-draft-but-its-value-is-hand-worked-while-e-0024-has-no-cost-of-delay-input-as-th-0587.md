---
id: TH-0725
title: S-0425 is a complete draft, but its value is hand-worked while E-0024 has no cost of delay input (as TH-0587)
anchor:
  path: wip/kanban/stories/S-0425-metrics-md-generalizes-strategic-usage-per-kind-to-usage-per-role-by-an-adr-so-a-project-s-own-role-is-charged-on-its-items-reported-by-flai-stats-and-charted-as-the-built-in-roles-are.md
  item: S-0425
status: open
participants: [orchestrator]
created: 2026-10-29T00:13:34Z
updated: 2026-10-29T00:13:34Z
---

# TH-0725 S-0425 is a complete draft, but its value is hand-worked while E-0024 has no cost of delay input (as TH-0587)

On wip/kanban/stories/S-0425-metrics-md-generalizes-strategic-usage-per-kind-to-usage-per-role-by-an-adr-so-a-project-s-own-role-is-charged-on-its-items-reported-by-flai-stats-and-charted-as-the-built-in-roles-are.md.

## Entries

### 2026-10-29T00:13:34Z orchestrator
**Recommended:** put the 4h you chose on TH-0471 on E-0024 (`flai edit E-0024 --time-lost-per-cycle 4h`). Then `flai cod` values S-0425 on its planner's next run, and I finalize it under finalize_drafts.

**The alternative:** say "finalize as it is", and I finalize S-0425 with the planner's value.

**Why I have not finalized it.**
- `flai promote --drafts` lists S-0425 as complete, and its criteria, touches, and 1h30m forecast fit its plan.
- Its value, 62.36 USD a week (by planner-S-0425, 2026-10-23T15:35:46Z), does not come from flai. `flai cod` refuses while E-0024 has no input, so the planner worked out the story's share by hand. TH-0587 holds S-0414 to S-0422 back for the same reason. This thread covers S-0425, which TH-0587 does not name.
- S-0425 also waits on S-0423 and S-0424 (`after`), which you moved back to the backlog. Finalizing it now would not get it pulled any sooner.
