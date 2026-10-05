---
id: TH-0145
title: "S-0227 plan: five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0227-the-analyzer-s-cost-is-recorded-on-the-issues-it-filed-and-the-stories-made-from-them.md
  item: S-0227
status: open
participants: [planner-S-0227]
created: 2026-10-05T05:46:14Z
updated: 2026-10-05T05:46:14Z
---

# TH-0145 S-0227 plan: five tasks in four layers

On wip/kanban/stories/S-0227-the-analyzer-s-cost-is-recorded-on-the-issues-it-filed-and-the-stories-made-from-them.md.

## Entries

### 2026-10-05T05:46:14Z planner-S-0227
I planned S-0227 with five tasks in four layers:

- Layer 1: T-0935. An ADR extends ADR-0083 to issues and the carry-over, and `design/system/continuous-improvement.md` shows `usage` in the issue file.
- Layer 2: T-0940. `issues.Issue` reads and writes `usage` with `strategic` entries, with a charge function shared by the next two tasks. `front-matter-fields.txt` lists the key.
- Layer 3: T-0942 and T-0944 run together, since they share no path.
  - T-0942: the analyzer's apportioned usage is charged evenly to the issues its activity names, or to S-0226's project total when it names none.
  - T-0944: `flai issue story` and MCP `issue_story` carry the issue's strategic usage onto the story and its epic.
- Layer 4: T-0948. `flai stats` reports usage per issue and counts it once. It waits for T-0944 because the count-once rule depends on how the carry-over marks the story.

Figures: forecast 50m (flai's 18m, raised against S-0225's 52m) and cost of delay 54.15 USD a week (`flai cod`'s share of E-0016's figure). The reasons are under Notes › Planning.

Assumptions:

- The issue keeps its own figures after a story is made from it, rather than moving them to the story. `flai stats` decides which one to count by whether a story links the issue. T-0935's ADR may decide otherwise.
- S-0226's project total exists by the time this story starts, as its `after` on S-0226 says. T-0942 uses it rather than adding its own.
- Only `I-nnnn` IDs among an analyzer activity's items are charged. Other IDs it names, such as stories, get nothing from the analyzer, which never authors stories.
- The dashboard is left out, because the criteria ask only for `flai stats`.
