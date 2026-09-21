---
id: T-0309
type: task
nature: feature
title: The project list shows each project's state at a glance, and one project's watcher or a slow answer does not stall another
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:24Z
updated: 2026-09-20T22:21:39Z
transitions:
  - to: ready
    at: 2026-09-20T21:49:37Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T21:49:38Z
    by: system-flow
  - to: done
    at: 2026-09-20T22:21:39Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0309 The project list shows each project's state at a glance, and one project's watcher or a slow answer does not stall another

## Work
The switcher/list page shows stories in review, threads awaiting the designer, and whether an agent is attending, per project, from what is already asked of flai. Caches, the SSE stream, and the notifier are confirmed per project (they already follow repo()/agent() once those are request-scoped); a test proves a slow or stalled project does not block another's page.

## Done when
- Tests: the at-a-glance summary for two projects; a stalled connection to one project while the other answers normally
- The flaiover tests pass

## Notes
