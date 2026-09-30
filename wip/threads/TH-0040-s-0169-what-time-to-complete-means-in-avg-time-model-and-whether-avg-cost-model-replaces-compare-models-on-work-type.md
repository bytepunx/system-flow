---
id: TH-0040
title: "S-0169: what \"time to complete\" means in Avg. Time / Model, and whether Avg. Cost / Model replaces \"compare models\" on $ / Work Type"
anchor:
  path: wip/kanban/stories/S-0169-replace-bad-charts-and-improve-titles.md
  item: S-0169
status: resolved
participants: [agent-S-0169, alex]
created: 2026-09-30T00:21:13Z
updated: 2026-09-30T00:38:27Z
---

# TH-0040 S-0169: what "time to complete" means in Avg. Time / Model, and whether Avg. Cost / Model replaces "compare models" on $ / Work Type

On wip/kanban/stories/S-0169-replace-bad-charts-and-improve-titles.md.

## Entries

### 2026-09-30T00:21:13Z agent-S-0169
Two questions before I build the two replacement charts. Both are drawn per bucket (hour, day, week, the "per" control) for the type chosen (story, task, epic), one line per model.

1. Avg. Time / Model: which time?
   - (a) recommended: agent minutes per item. The mean agent time of the items each model worked on, from the items' usage `seconds`, as flai already records it. It measures the model's own work, and leaves out the time a story waits in review for you.
   - (b) cycle time per item. The mean time from in-progress to done of the items each model worked on. It includes the time in review.

2. Avg. Cost / Model: it plots the mean dollars per item for each model. "Cost per item", renamed "$ / Work Type" by this story, already draws exactly that when its "compare" control is set to "models".
   - (a) recommended: build Avg. Cost / Model as its own chart, and drop the "compare" control from $ / Work Type (and from Tokens per item, for consistency), so it always shows one line per item type.
   - (b) build Avg. Cost / Model and leave $ / Work Type's "compare" control as it is. The same view is then reachable two ways.

Meanwhile I am doing the title changes and removing aging and estimates.

### 2026-09-30T00:24:38Z alex
1 - a, 2 - a

### 2026-09-30T00:38:27Z agent-S-0169
Resolved: 1a and 2a built: Avg. Time / Model plots agent minutes per item; Avg. Cost / Model is its own chart and the compare control is gone (ADR-0057).
