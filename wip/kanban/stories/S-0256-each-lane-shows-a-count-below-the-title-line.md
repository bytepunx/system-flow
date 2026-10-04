---
id: S-0256
type: story
nature: improvement
title: Each lane shows a count below the title line
status: backlog
owner: alex
created: 2026-10-04T03:16:14Z
updated: 2026-10-04T03:16:14Z
transitions: []
tags: [dashboard]
topics: [dashboard-board]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0256 Each lane shows a count below the title line

## Goal

Below the title line of every lane, show the count of epics, stories and tasks in that lane like so: `1 | 2 | 4`. A tooltip when hovering over it should show the long form: `1 epic, 2 stories, 4 tasks`.

## Acceptance criteria
- [ ] Every lane has the counts of each work item type below the title line in the format `{epic count} | {story count} | {task count}`
- [ ] When hovering over the lane count, show a longer format in the tool tip: `{epic count} epic, {story count} story, {task count} task` with each work type using the correct pluralization.

## Tasks

## Notes
