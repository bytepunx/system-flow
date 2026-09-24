---
id: T-0427
type: task
nature: feature
title: Runbooks for backing up, restoring, and migrating flai and flaiover
status: done
parent: S-0028
owner: alex
created: 2026-09-24T09:28:02Z
updated: 2026-09-24T09:41:05Z
transitions:
  - to: ready
    at: 2026-09-24T09:28:10Z
    by: agent-S-0028
  - to: in-progress
    at: 2026-09-24T09:38:17Z
    by: agent-S-0028
  - to: done
    at: 2026-09-24T09:41:05Z
    by: agent-S-0028
stream: S-0028
tags: []
touches: [docs/operators]
---
# T-0427 Runbooks for backing up, restoring, and migrating flai and flaiover

## Work
- `backup.md`, `restore.md`, and `migrate.md`, each with a section for flai and one for flaiover: what state exists and where (configuration, `serve/` and `host/` state, secrets, per-project `.flai-cache`, story worktrees), what is worth keeping, the commands in order, and how to check the result. `migrate.md` covers moving to another machine, `flai upgrade`, `flai migrate`, and the dashboard's own upgrade notes, linked rather than copied.

## Done when
- The three pages exist and are in the runbooks index, and every file they name is one flai writes.

## Notes
