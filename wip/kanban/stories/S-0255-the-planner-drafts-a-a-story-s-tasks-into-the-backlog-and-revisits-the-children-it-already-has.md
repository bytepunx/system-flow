---
id: S-0255
type: story
nature: feature
title: The planner drafts a a story's tasks into the backlog and revisits the children it already has
status: backlog
parent: E-0016
owner: alex
created: 2026-10-04T00:51:25Z
updated: 2026-10-04T00:51:25Z
transitions: []
tags: [cli]
topics: [planner-agent]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0255 The planner drafts a a story's tasks into the backlog and revisits the children it already has

## Goal

Given a story, the planner writes the tasks that deliver its outcome, in the backlog, and when the story already has tasks it revisits each open one: re-enriches it, proposes splits or merges, and reports what changed.

## Acceptance criteria
- [ ] For a story with no tasks, the planner writes tasks with work, done when, nature, tags, topics, touches, and `after` between them, each in the backlog, and opens one thread on the story summarizing the plan: the tasks, their order, and the assumptions it made
- [ ] For a story with tasks, it revisits every task not done or cancelled: re-enriches it (the enrichment story), and proposes in the story's thread any task it would split, merge, add, or drop, creating tasks for additions but never cancelling or rewriting a finalized task's words without asking
- [ ] Tasks it writes pass `flai check --strict` and the markdown lint before they are kept
- [ ] Its log entry names the tasks created and revisited and the run's cost
- [ ] Measured on E-0016 itself or a comparable epic: the designer judges the drafts on a thread and the result is recorded in `design/system/strategic-agents.md`

## Tasks

## Notes
