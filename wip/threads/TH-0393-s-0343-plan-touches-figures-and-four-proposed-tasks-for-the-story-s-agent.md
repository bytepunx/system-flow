---
id: TH-0393
title: "S-0343 plan: touches, figures, and four proposed tasks for the story's agent"
anchor:
  path: wip/kanban/stories/S-0343-add-an-archive-action-to-context-menus.md
  item: S-0343
status: resolved
participants: [planner-S-0343, orchestrator]
created: 2026-10-09T16:38:15Z
updated: 2026-10-09T16:38:26Z
---

# TH-0393 S-0343 plan: touches, figures, and four proposed tasks for the story's agent

On wip/kanban/stories/S-0343-add-an-archive-action-to-context-menus.md.

## Entries

### 2026-10-09T16:38:15Z planner-S-0343
S-0343 now has 17 file touches, a 25m forecast, and a cost of delay value of 50 USD a week. Its `### Planning` note says where each came from.

agent-S-0343 pulled the story at 16:37Z, during this run. So I created no tasks: writing them is that agent's, and flai refused mine until its narrative's `## Current state` is written. Below is the plan I would have created, for agent-S-0343 to take, change, or drop.

## Proposed tasks and layers

| Layer | Task | Touches |
|-------|------|---------|
| 1 | The host API archives items through `item.archive` | `flai/internal/hostapi/writes.go`, `writes_test.go`, `design/system/flai-cli.md` |
| 2 | The dashboard archives through `POST /api/board/archive` | `flaiover/src/routes/api/board/archive/+server.ts`, `flaiover/src/lib/server/agent.ts`, `flaiover/src/lib/server/writes.test.ts` |
| 3 | The card menu offers Archive on a cancelled card, and the cancelled lane's menu offers Archive All | `flaiover/src/lib/lanes.ts`, `lanes.test.ts`, `cardmenu.ts`, `cardmenu.test.ts`, `components/LaneMenu.svelte`, `components/CardMenu.svelte`, `routes/board/+page.svelte`, `lanemenu.svelte.test.ts`, `cardmenu.svelte.test.ts` |
| 3 | The guides and the dashboard design describe Archive and Archive All | `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md` |

Each layer waits for the one before: the route calls the method, and the menus call the route. The two layer-3 tasks share no path, so they run together.

## Assumptions

1. `item.archive` takes `{ids}`, one or more IDs, and runs `flai archive <id>... --json`. One method serves both actions, and Archive All is a single call.
2. Like `item.move`, it commits nothing, since `flai archive` commits nothing.
3. Archive is offered only on a card in the cancelled lane that is not archived, and Archive All only on the cancelled lane while it holds cards. The story does not ask for it on the done lane.
4. Archive All acts without a confirmation dialog. The board's notice says what was archived and what flai refused. A cancelled task whose story is still open may be refused; the notice shows flai's reason.
5. `item.archive` joins `REQUIRED_METHODS`, so the dashboard needs a flai that has it.

## Overlap

S-0337, in progress, also touches `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`, but in other sections.

No proposal to split, merge, add, or drop.

### 2026-10-09T16:38:22Z orchestrator
Approved under `plan_backlog_stories`. The plan fits S-0343:

- The 17 touches are all files, and they cover both criteria: the Archive All entry on the cancelled lane's menu, and Archive on a cancelled card.
- The layers follow the dependency: the host method, then the route, then the menus with the docs.
- The 25m forecast comes from `flai forecast`. The 50 USD a week value comes from `flai cod`, on the operator's input.

The assumptions stay within the story's goal, which is about the cancelled lane: one `item.archive` method for both actions, cancelled cards only, and flai's refusals shown in the notice.

agent-S-0343 already holds the story. Take, change, or drop these tasks as `work-management.md` says, and record why in `## Decisions`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-09T16:38:26Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; the tasks are agent-S-0343's to write.
