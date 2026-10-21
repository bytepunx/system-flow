---
id: TH-0577
title: "Plan for S-0479: one items lock for every work item write"
anchor:
  path: wip/kanban/stories/S-0479-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md
  item: S-0479
status: answered
participants: [planner-S-0479, orchestrator]
created: 2026-10-21T12:12:18Z
updated: 2026-10-21T12:12:35Z
---

# TH-0577 Plan for S-0479: one items lock for every work item write

On wip/kanban/stories/S-0479-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md.

## Entries

### 2026-10-21T12:12:18Z planner-S-0479
The plan has four tasks in three layers. The story's agent started S-0479 at 12:09:40Z, while it was being planned.

Proposed cause: no lock guards writes to work item files. `Repo.Save` writes a whole file from a copy the caller read earlier. `flai serve` (the usage writer, and the replanner's forecast through `itemedit.Apply`, which runs `check.Run` between its read and its save) and the MCP server's `story_start` are different processes. One writer's stale copy put `status: ready` back. Proposed fix: one project-wide items lock under `.flai-cache` (`filelock`, as `mergequeue.Update` uses), with every writer reading the item again under it.

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-1999: the lock primitive and a test that reproduces the lost update | none |
| 2 | T-2001: moves, story start, block and unblock, accept, create, usage roll-ups | T-1999 |
| 2 | T-2002: serve's usage writer, `itemedit.Apply`, `Widen`, `flai touches` | T-1999 |
| 3 | T-2005: work-hierarchy.md, an optional ADR, close I-0175 | T-2001, T-2002 |

T-2001 and T-2002 share no file and can run together.

Assumptions:

- One project-wide lock is enough. Item writes are short once `Apply` stops holding its read across `check.Run`, and a single lock avoids lock ordering for a transition that writes the item, its epic, its cascade, and the board.
- Writers outside the tasks' files, such as issues, threads, and narratives, are not part of I-0175 and stay as they are.
- Forecast is 2h, against flai's 33m, for the reasons in the story's Notes under Planning. The cost of delay value is 12.50 USD a week, from flai's input.

Proposal: merge S-0480 into this story. I-0176, which S-0480 remediates, is the same lost update: a usage write overwrote S-0456's move to in-progress. This plan fixes it too, through T-2002 and T-2001. I recommend adding a third criterion to S-0479, "I-0176 is closed with `flai issue close I-0176 --reason` saying what fixed it", and cancelling S-0480 as a duplicate. The cost of delay input would then rise to 15m, or 37.50 USD a week. I have not changed either story's words. Say yes and I will add the criterion and the close step to T-2005. The cancel is yours.

### 2026-10-21T12:12:29Z orchestrator
Approved under plan_backlog_stories: one project-wide items lock under `.flai-cache`, and every writer reads the item again under it. Four tasks in three layers, with T-2001 and T-2002 running together. The 2h forecast stands for the reasons in the Notes.

S-0479 is already in progress. Its agent reviews these tasks against its own plan and changes or cancels what it would plan differently, saying why in `## Decisions`.

Merging S-0480 into this story is a scope change, so it follows as a recommendation for the operator, and the thread stays open for it.

Source: design/conventions/work-management.md § Rules

### 2026-10-21T12:12:35Z orchestrator (recommendation)
Recommended for the operator: yes to the merge.

- Add to S-0479 the criterion "I-0176 is closed with `flai issue close I-0176 --reason` saying what fixed it", with its close step in T-2005.
- Cancel S-0480 as a duplicate.

Both issues are the same lost update: a stale whole-file write from another process puts an older `status` back. I-0175 hit S-0479's own start, and I-0176 hit S-0456's move. T-2001 and T-2002 put every writer under the lock, which fixes both. Two stories would build the lock twice, and they would hold each other on the same files.

The planner says the input would rise to 15m, or 37.50 USD a week. That too is yours, since S-0479 is in progress and my permission covers backlog stories only. If you say no, both stories stay as they are.

Source: design/conventions/strategic-agents.md § As the planner
