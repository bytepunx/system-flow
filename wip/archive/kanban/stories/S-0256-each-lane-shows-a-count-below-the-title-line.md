---
id: S-0256
type: story
nature: improvement
title: Each lane shows a count below the title line
status: done
owner: alex
created: 2026-10-04T03:16:14Z
updated: 2026-10-04T21:41:59Z
transitions:
  - to: ready
    at: 2026-10-04T04:00:41Z
    by: alex
  - to: in-progress
    at: 2026-10-04T21:23:26Z
    by: agent-S-0256
  - to: review
    at: 2026-10-04T21:36:25Z
    by: agent-S-0256
  - to: done
    at: 2026-10-04T21:41:59Z
    by: alex
tags: [dashboard]
topics: [dashboard-board]
touches: [flaiover/src/lib/lanes.test.ts, flaiover/src/lib/lanes.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/issues/I-0072-flai-s-wip-markdown-lint-does-not-flag-a-space-inside-a-code-span-so-a-task-body-reached-main-and-failed-a-story-s-close-out.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
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
- [x] Every lane has the counts of each work item type below the title line in the format `{epic count} | {story count} | {task count}`
- [x] When hovering over the lane count, show a longer format in the tool tip: `{epic count} epic, {story count} story, {task count} task` with each work type using the correct pluralization.

## Tasks
- T-0827 Each lane shows its epic, story, and task counts below its title
- T-0828 The board's user guide and design describe the lane counts

## Notes
