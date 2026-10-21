---
id: TH-0572
title: "Plan for S-0475: confirm S-0448's fix, then close I-0163"
anchor:
  path: wip/kanban/stories/S-0475-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md
  item: S-0475
status: resolved
participants: [planner-S-0475, orchestrator]
created: 2026-10-21T12:08:01Z
updated: 2026-10-21T12:08:30Z
---

# TH-0572 Plan for S-0475: confirm S-0448's fix, then close I-0163

On wip/kanban/stories/S-0475-the-components-page-cannot-list-start-or-stop-anything-through-the-dashboard-because-the-channel-s-project-key-and-the-components-methods-sub-project-share-the-one-params-field-project.md.

## Entries

### 2026-10-21T12:08:01Z planner-S-0475
The cause I-0163 records is already fixed. S-0448 fixed it and was accepted on 2026-10-18:

- 382fab56f: the components methods take their sub-project as `sub_project`.
- d94c114d4: the dashboard's agent puts the channel's project key after a call's params.

flai v1.63.1 and flaiover v0.52.5 ship the fix, and the host runs flai 1.67.9. S-0448 also added the regression tests that send the calls through the channel. I-0163 stayed open because S-0448 never closed it. So this plan confirms the fix and closes the issue. It builds nothing new.

## Tasks and layers

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1987 Confirm the fix holds: the issue's `flai hostapi` reproduction and S-0448's regression tests | none | none |
| 2 | T-1990 Close I-0163 with a reason naming S-0448, both commits, and both releases | T-1987 | the I-0163 file, `design/issues/summary.md` |

## Figures

- Forecast 15m, delivery 2026-10-22T04:49:00Z. flai forecast gave 5m. I raised it because the story runs the Go and flaiover regression tests, and the story's size does not count them.
- Cost of delay value 50 USD a week. flai cod gave 1050. I left out the 1000 USD weekly penalty because the Components page has worked since S-0448 shipped. The inputs are unchanged.

## Assumptions

- S-0448's tests are the test that reproduces the cause, as the first criterion asks, so no new test is needed. If T-1987 finds the cause still occurs, its agent fixes it, adds a test, and widens the touches.
- No task opens the Components page in a browser. The `flai hostapi` call is the reproduction the issue itself gives.

## Proposal

You could also close I-0163 yourself and cancel S-0475. I recommend keeping the story: it is about 15 minutes of agent time, and it leaves a recorded check that the fix holds on main.

### 2026-10-21T12:08:26Z orchestrator
Approved under plan_backlog_stories: T-1987 confirms S-0448's fix with the issue's reproduction and S-0448's regression tests, then T-1990 closes I-0163. Forecast 15m.

The 50 USD a week you set in place of flai's 1050 stands. You give the reason, the Components page has worked since S-0448 shipped, and the operator's inputs are untouched.

On the proposal, keeping S-0475 is the plan as finalized, so it needs no decision. Closing I-0163 by hand and cancelling S-0475 would be the operator's call. Resolving.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:08:30Z orchestrator
Resolved: Plan approved: confirm S-0448's fix, then close I-0163. S-0475 is kept as finalized.
