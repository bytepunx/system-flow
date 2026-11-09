---
id: TH-0939
title: S-0626 is a draft with no forecast and no cost of delay value, and flai plan --candidates does not list it
anchor:
  path: wip/kanban/stories/S-0626-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md
  item: S-0626
status: resolved
participants: [orchestrator]
created: 2026-11-09T17:16:28Z
updated: 2026-11-09T17:17:51Z
---

# TH-0939 S-0626 is a draft with no forecast and no cost of delay value, and flai plan --candidates does not list it

On wip/kanban/stories/S-0626-design-documents-describe-a-conventions-page-a-streams-page-a-mounted.md.

## Entries

### 2026-11-09T17:16:28Z orchestrator
**Recommended:** run the planner for S-0627's sibling once more by hand, `flai plan S-0626`, so it writes the forecast and works out the value from the input it has.

**What `flai promote --drafts` says is missing:**

- forecast duration
- forecast delivery
- cost of delay value

**What it has:** five tasks in two layers (T-2755 to T-2758, then T-2759), eleven file touches, and the input `time_lost_per_cycle: 15m`, which I set on TH-0930 as its planner recommended.

**Why I do not start the planner myself:** `flai plan --candidates` does not list S-0626, neither to plan nor as left out with a planner running. Under `plan_backlog_stories` I start the planner only for a story it lists. A planner run ended after the input was set, and the story's touches, tags, and topics changed at 2026-11-09T17:13:28Z, but no forecast or value was written.

**Alternative:** set the forecast yourself with `flai edit S-0626 --forecast-duration <d>` and `flai cod S-0626` for the value.

Until it has them, I cannot finalize S-0626 under `finalize_drafts`.

### 2026-11-09T17:17:51Z orchestrator
Resolved: Its concern is gone: planner-S-0626 wrote S-0626's forecast and cost of delay value at 17:14:37Z, while this thread was being written, and flai promote --drafts now lists it complete. Nothing is owed here.
