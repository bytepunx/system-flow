---
id: T-0299
type: task
nature: feature
title: The design and the users' documentation say how items are edited, and stop saying they can be retitled from the board
status: done
parent: S-0085
owner: alex
created: 2026-09-20T15:21:43Z
updated: 2026-09-20T15:46:27Z
transitions:
  - to: ready
    at: 2026-09-20T15:44:27Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:44:27Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:46:27Z
    by: system-flow
stream: S-0085
tags: []
---
# T-0299 The design and the users' documentation say how items are edited, and stop saying they can be retitled from the board

## Work
docs/users/flaiover.md and flai.md; design/system/flai-cli.md, flaiover-dashboard.md, agent-narrative.md, work-hierarchy.md where it says who owns which field. The template's copy of anything that changes, first.

## Done when
- The markdown lint and flai check --strict pass after the last edit, in the worktree and in the main checkout

## Notes
