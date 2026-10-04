---
id: TH-0099
title: S-0255 relies on T-0792 and T-0793 covering a story's tasks too
anchor:
  path: wip/kanban/stories/S-0209-the-planner-drafts-an-epic-s-stories-into-the-backlog-and-revisits-the-children-it-already-has.md
  item: S-0209
status: resolved
participants: [agent-S-0255, alex]
created: 2026-10-04T04:05:28Z
updated: 2026-10-04T04:06:36Z
---

# TH-0099 S-0255 relies on T-0792 and T-0793 covering a story's tasks too

On wip/kanban/stories/S-0209-the-planner-drafts-an-epic-s-stories-into-the-backlog-and-revisits-the-children-it-already-has.md.

## Entries

### 2026-10-04T04:05:28Z agent-S-0255
From agent-S-0255 (the planner drafts a story's tasks, in progress beside S-0209). S-0255 has the same two criteria for tasks: tasks the planner writes pass flai check --strict and the lint before they are kept, and its log entry names the tasks created and revisited with the run's cost. T-0792 (item_new with a body checked through itemnew.Create) and T-0793 (planEnded's entry names the planned item and the items under it created or changed during the run) deliver both for a story's tasks too, provided they are not limited to stories under an epic: item_new type task with a body, and a planned story's tasks as "the items under it". I am not duplicating them on story/S-0255; I will add task cases to their tests after S-0209 is accepted and I sync. Please keep both generic over the item's type, or say here if you will not.

S-0255 also edits planPrompt's story text, the guard's cliPlans (adding task new), and the convention's "As the planner": expect small conflicts there on whichever of us syncs second.

### 2026-10-04T04:06:36Z alex
Resolved.
