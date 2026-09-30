---
id: TH-0039
title: "Time in state on a time axis: a bar per item, or per day?"
anchor:
  path: wip/kanban/stories/S-0168-several-charts-don-t-use-window-correctly.md
  item: S-0168
status: resolved
participants: [agent-S-0168, alex]
created: 2026-09-29T23:58:04Z
updated: 2026-09-30T00:18:31Z
---

# TH-0039 Time in state on a time axis: a bar per item, or per day?

On wip/kanban/stories/S-0168-several-charts-don-t-use-window-correctly.md.

## Entries

### 2026-09-29T23:58:04Z agent-S-0168
Time in state draws its x axis as a list of item IDs, so it has no time range to follow the window. Only the set of items changes. To give it a time axis that spans the window (ADR-0054), I see two layouts:

1. (recommended) One stacked bar per day of the window: the mean hours per state of the items completed that day. The tooltip lists the items, and the table keeps one row per item. It reads cleanly at 1d to 365d. On this repo, items often complete minutes apart, so bars per item would overlap.
2. One stacked bar per item, placed at its completion time. This is closest to metrics.md ("Stacked bar per story completed in the window"). But bars completed close together overlap or shrink to hairlines.

Either one changes the chart's definition in metrics.md, so I'd record it in an ADR that refines ADR-0054.

On cycle time, the running dashboard (flaiover 0.31.1) drew the right axis in every path I tried: 30d, 7d, and 1d, opening the chart directly, and arriving from another chart. What I did find is that the charts page forgets the window whenever it remounts. After you go to the Board and come back through Charts, cycle time loads at 30d, whatever you chose. The page also waits for the epic list before it asks for stats. I'm fixing both: the window will be remembered, and the stats load straight away. If you saw something else on cycle time, tell me what you did and what the axis showed.

### 2026-09-30T00:13:55Z alex
1

### 2026-09-30T00:18:26Z agent-S-0168
Done as option 1, recorded in ADR-0056 (it refines ADR-0054). Time in state now draws one stacked bar per day of the window: the mean hours per state of the items completed that day. Its time axis runs from the day that holds the window's start to today. Hover a bar to see the items it averages; the table still lists each item. Commit 96df977 on story/S-0168.

### 2026-09-30T00:18:31Z agent-S-0168
Resolved: Option 1 chosen by alex; implemented in S-0168 and recorded in ADR-0056.
