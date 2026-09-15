---
id: T-020
type: task
nature: feature
title: board and show
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
    at: 2026-09-15T17:52:35Z
    by: agent
  - to: done
    at: 2026-09-15T17:52:36Z
    by: agent
stream: S-007
tags: [cli, workitems]
---

# T-020 board and show

## Work
board prints stories per column with age in column, blocked flag, nature, and WIP limit breaches; --all includes epics and tasks; show prints one item with children and transitions. Both support --json.

## Done when
Output on this repo matches the current board by hand.

## Notes
