---
id: TH-1297
title: "Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange"
anchor:
  path: wip/kanban/stories/S-0838-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md
  item: S-0838
status: resolved
participants: [planner-S-0838, orchestrator]
created: 2026-11-17T12:32:44Z
updated: 2026-11-17T12:35:12Z
---

# TH-1297 Plan for S-0838: read each lane's kind from the board's flow in the epic filters, the move route, and ProposeChange

On wip/kanban/stories/S-0838-flaiover-still-compares-state-names-in-the-epic-filters-and-the-move-route-so.md.

## Entries

### 2026-11-17T12:32:44Z planner-S-0838
Plan for S-0838. The fix: every site I-0399 names reads the board's flow and asks `$lib/workflow` what a state is, as S-0758's views already do.

## Tasks and layers

| Layer | Task | Changes |
|-------|------|---------|
| 1 | T-3568 | The move route reads the flow from the repo's board view. It uses `isCancelled` for the dry-run preview and `isAccepted` for the 600000 ms acceptance timeout. Also adds a renamed-lanes fixture to `lanes.fixture.ts` and `kinds.test.ts` |
| 2 | T-3571 (after T-3568) | The design page filters epics with `isOpen` and passes the flow to `ProposeChange`. `ProposeChange` moves the new story through `movesTo` to the first pullable start lane, not `ready`. Its examples get the new prop |
| 2 | T-3575 (after T-3568) | `NewItemForm` and `DeploySettings` filter epics with `isOpen` over the flow from `/api/board` |
| 3 | T-3579 (after T-3571, T-3575) | `design/system/flaiover-dashboard.md` § A state's kind lists the new readers. A search confirms only the task plan's and charts' comparisons remain. `flai issue close I-0399` |

T-3568 comes first because its fixture is what the layer-2 tests read. T-3571 and T-3575 share no file, so they run together.

## Figures

- Forecast 34m, as `flai forecast` gives it. The size over-counts flai-written files and under-counts four tests that need a `/api/board` answer, so I made no adjustment.
- Cost of delay value 25 USD a week, from flai's input of 10m per cycle.

## Assumptions

- The repo's board view (`repo().boardView()`) is cached, so the move route can read the flow on every call without a second round trip of note.
- In the move route, a dry run to a lane that is not of kind cancelled stays a plain move, as today.
- `docs/users/flaiover.md` needs no change: it already speaks of open epics, which the fix makes true.
- `$lib/taskplan.ts` and `$lib/viz/charts.ts` stay out of scope. The design section keeps the task plan's states, and S-0754 owns the charts.
- With no pullable start lane, `ProposeChange` leaves the story in the initial lane and says so.

No task to split, merge, or drop.

### 2026-11-17T12:34:54Z orchestrator
Approved. The plan fits S-0838.

- T-3568, T-3571, and T-3575 make each place I-0399 names read a lane's kind from the board's flow: the move route, the design page, `ProposeChange`, `NewItemForm`, and `DeploySettings`. They test it on a renamed-lanes fixture. That meets criterion 1.
- T-3579 updates `design/system/flaiover-dashboard.md` and closes I-0399. That meets criterion 2.
- The touches are 15 files and no folder. They include the `ProposeChange` examples that `flai check` reads. Each task's touches are inside them.
- The layers follow the fixture that T-3568 writes and the layer-2 tests read. T-3571 and T-3575 share no file.
- The forecast of 34m and the value of 25 USD a week are flai's figures.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-17T12:35:12Z orchestrator
Resolved: The orchestrator approved S-0838's plan under plan_backlog_stories.
