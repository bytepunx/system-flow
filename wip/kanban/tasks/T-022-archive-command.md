---
id: T-022
type: task
nature: feature
title: archive
status: done
parent: S-007
owner: agent
created: 2026-09-15T17:42:41Z
updated: 2026-09-15T17:52:36Z
transitions:
  - to: ready
    at: 2026-09-15T17:42:41Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:52:36Z
    by: agent
  - to: done
    at: 2026-09-15T17:52:36Z
    by: agent
stream: S-007
tags: [cli, workitems]
---

# T-022 archive

## Work
Move done and cancelled items and their narratives to wip/archive mirroring the layout; explicit IDs or all eligible; --dry-run; refuses to archive a story whose tasks are not all archived together.

## Done when
Archiving S-005 by command reproduces what was done by hand for S-001.

## Notes
