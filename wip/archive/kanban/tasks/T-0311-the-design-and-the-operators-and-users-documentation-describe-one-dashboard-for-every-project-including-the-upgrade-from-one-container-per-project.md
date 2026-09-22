---
id: T-0311
type: task
nature: feature
title: The design and the operators' and users' documentation describe one dashboard for every project, including the upgrade from one container per project
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:25Z
updated: 2026-09-20T22:31:46Z
transitions:
  - to: ready
    at: 2026-09-20T22:25:59Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T22:25:59Z
    by: system-flow
  - to: done
    at: 2026-09-20T22:31:46Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0311 The design and the operators' and users' documentation describe one dashboard for every project, including the upgrade from one container per project

## Work
docs/operators/index.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, dashboard-host-channel.md. What an operator running several projects today must do to move to the shared dashboard: stop each project's old container, run flai dashboard once per project against the new one.

## Done when
- The markdown lint and flai check --strict pass after the last edit, in the worktree and the main checkout

## Notes
