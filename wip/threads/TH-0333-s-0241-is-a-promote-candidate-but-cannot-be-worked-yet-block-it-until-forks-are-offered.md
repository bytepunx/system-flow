---
id: TH-0333
title: "S-0241 is a promote candidate but cannot be worked yet: block it until forks are offered?"
anchor:
  path: wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md
  item: S-0241
status: resolved
participants: [orchestrator]
created: 2026-10-07T23:16:59Z
updated: 2026-10-08T07:09:00Z
---

# TH-0333 S-0241 is a promote candidate but cannot be worked yet: block it until forks are offered?

On wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md.

## Entries

### 2026-10-07T23:16:59Z orchestrator
Recommendation: block S-0241 with `flai block S-0241 --reason "no headless session offers subagent_type fork yet (Claude Code 2.1.290)"`, and unblock it when one does.

Why I ask:

- S-0241's plan is now complete: 3 tasks, 2 file touches, a 2h forecast, and the `penalty_per_week: 25` input I set on TH-0332.
- So `flai promote --candidates` lists it as the sole candidate (rank 1, 2h). Under `promote_to_ready` I move candidates to ready in order.
- Its Notes say: "Do not move it to ready until a headless session offers forks: it cannot be worked before then." The planner's own headless session, on Claude Code 2.1.290, offered no `fork` type today.
- If I promote it, flai serve starts its agent. T-1267 would find no fork and block the story in progress, which spends a session and holds an in-progress slot.

Until you answer, I leave S-0241 in the backlog and promote nothing in its place. Alternatives:

- Promote it anyway, and let T-1267 confirm and block.
- Make it a draft again, which keeps it out of `flai promote --candidates`.

### 2026-10-08T07:09:00Z orchestrator
Resolved: The operator blocked S-0241 at 07:08Z ("No candidate solutions yet"), as recommended. It stays out of ready until it is unblocked.
