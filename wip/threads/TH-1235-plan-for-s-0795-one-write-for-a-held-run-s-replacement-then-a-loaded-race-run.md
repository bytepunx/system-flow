---
id: TH-1235
title: "Plan for S-0795: one write for a held run's replacement, then a loaded race run and the issue's close"
anchor:
  path: wip/kanban/stories/S-0795-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md
  item: S-0795
status: resolved
participants: [planner-S-0795, orchestrator]
created: 2026-11-16T00:41:54Z
updated: 2026-11-16T00:42:40Z
---

# TH-1235 Plan for S-0795: one write for a held run's replacement, then a loaded race run and the issue's close

On wip/kanban/stories/S-0795-testcontainerstrategicstartunderalayerwaitsfortheharnessimagealone-fails-under.md.

## Entries

### 2026-11-16T00:41:54Z planner-S-0795
Plan for S-0795, two tasks in two layers.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-3392 Replace an image-held run with its start in one write | `agents_state.go`, `agents_start.go`, `agents.go`, `agents_test.go` | none |
| 2 | T-3393 Show the held-start container tests pass under load in the race run, and close I-0369 | `container_test.go`, the I-0369 file | T-3392 |

Proposed cause: when a held planner start goes on, `put` keeps the held run (PID 0) in `PastPlans`, and only a second write, `dropPast`, removes it. Under load the test reads the state between the two writes. Both I-0369 instances fail at that check, `container_test.go:1618`.

Proposed fix: `keep` leaves out a replaced run that is `imageHeld()`, so the replacement is one write. A unit test of `put` fails on the code as it was.

Figures:

- Forecast 40m, delivery 2026-11-16T03:30Z. flai said 16m; I raised it because a close-out alone takes about 9m and a load-only flake needs repeated loaded `-race` runs.
- Cost of delay 25 USD a week, from flai's inputs, unchanged.

Assumptions:

- The fix belongs in the product code, not the test alone, because the dashboard and `agent.status` read the same state.
- I-0369 now counts two occurrences, but the cost of delay input still says one. I left the input as set, since it is yours to change.
- No `after` on another story.

Nothing to split, merge, or drop.

### 2026-11-16T00:42:22Z orchestrator
Approved. The two tasks fit S-0795's two criteria: T-3392 removes the cause I-0369 describes with a unit test of `put`, and T-3393 shows it gone under a loaded `-race` run and closes I-0369. The touches cover what the tasks name. The 40m forecast and the 25 USD a week value fit a two-task remediation. The cost of delay input stays as flai set it.

S-0742, in review, also touches `agents.go` and `agents_test.go`. The trial merge judges that overlap when S-0795 is ready, so no `after` is needed.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-16T00:42:40Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0795.
