---
id: T-0487
type: task
nature: remediation
title: The story page offers to have an agent commit the worktree's outstanding work
status: done
parent: S-0140
owner: alex
created: 2026-09-26T21:11:25Z
updated: 2026-09-26T21:22:01Z
transitions:
  - to: ready
    at: 2026-09-26T21:11:31Z
    by: agent-S-0140
  - to: in-progress
    at: 2026-09-26T21:19:33Z
    by: agent-S-0140
  - to: done
    at: 2026-09-26T21:22:01Z
    by: agent-S-0140
stream: S-0140
tags: []
touches: [flaiover/src]
---
# T-0487 The story page offers to have an agent commit the worktree's outstanding work

## Work

- When the acceptance preview lists uncommitted paths in the story worktree, the story page shows them and a button that calls the new host method.

## Done when

- [x] flaiover tests cover the endpoint and the button; `npm run check` and lint are clean.

## Notes
