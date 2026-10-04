---
id: TH-0098
title: "Tasks carry no topics: what does criterion 1's \"topics\" mean for a task?"
anchor:
  path: wip/kanban/stories/S-0255-the-planner-drafts-a-story-s-tasks-into-the-backlog-and-revisits-the-children-it-already-has.md
  item: S-0255
status: resolved
participants: [agent-S-0255, alex]
created: 2026-10-04T04:03:31Z
updated: 2026-10-04T04:08:40Z
---

# TH-0098 Tasks carry no topics: what does criterion 1's "topics" mean for a task?

On wip/kanban/stories/S-0255-the-planner-drafts-a-story-s-tasks-into-the-backlog-and-revisits-the-children-it-already-has.md.

## Entries

### 2026-10-04T04:03:31Z agent-S-0255
Criterion 1 asks for tasks with "work, done when, nature, tags, topics, touches, and after". flai keeps topics on stories and epics only: item_new documents topics as story/epic, and item_edit refuses them on a task ("topics belong to stories and epics", flai/internal/itemedit/itemedit.go). A story's topics are what select its prime pack, and a task is worked inside its story's pack, so a task's topics would select nothing.

Recommended: tasks keep no topics. The planner checks the story's topics against the tasks it drafts and adds a topic to the story when a task reaches one the story lacks (item_edit on the story; topics are not among the operator's inputs). I read criterion 1's "topics" that way and record it in Decisions.

Alternative: give tasks topics, which means a schema change (work-hierarchy.md, check, item_new/item_edit) and deciding what prime does with them, which I would rather leave out of this story.

I go on with the recommended reading meanwhile; say if you want the alternative.

### 2026-10-04T04:04:34Z alex
you are correct, that is an error to have included them in the description (you may change it with an edit)

### 2026-10-04T04:08:40Z agent-S-0255
Resolved.
