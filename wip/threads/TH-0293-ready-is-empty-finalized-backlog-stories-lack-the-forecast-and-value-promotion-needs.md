---
id: TH-0293
title: "Ready is empty: finalized backlog stories lack the forecast and value promotion needs"
anchor:
  path: wip/kanban/board.md
status: resolved
participants: [orchestrator, alex]
created: 2026-10-07T14:58:21Z
updated: 2026-10-08T04:32:20Z
---

# TH-0293 Ready is empty: finalized backlog stories lack the forecast and value promotion needs

On wip/kanban/board.md.

## Entries

### 2026-10-07T14:58:21Z orchestrator
Recommendation: ask the planner to plan the finalized backlog stories below (`flai plan S-nnnn`, or Plan on the dashboard), starting with S-0297, which lacks only a forecast. I promote each one to ready as soon as `flai promote --candidates` lists it.

Ready and review are empty, and in progress has room for one more (S-0293 and S-0298 run). `flai promote --candidates` lists none, and `flai plan --candidates` lists no epic. My `plan` permission covers only those epics, so I cannot start the planner for these stories.

What each lacks:

| Story | Lacks | Also held by |
|-------|-------|--------------|
| S-0297 | forecast | |
| S-0305 | forecast, value | S-0311's claim, now cleared |
| S-0232 | forecast, value | S-0298 (docs/operators) |
| S-0288, S-0289, S-0290, S-0291, S-0304, S-0306 | touches, forecast, value | |
| S-0241 | touches, forecast, value | |
| S-0233 to S-0239 | forecast, value | `after` S-0232 |

The drafts S-0280, S-0287, S-0309, S-0310, S-0312, and S-0313 have their own threads (TH-0252 to TH-0255, TH-0288, TH-0289).

### 2026-10-08T04:32:20Z alex
Resolved.
