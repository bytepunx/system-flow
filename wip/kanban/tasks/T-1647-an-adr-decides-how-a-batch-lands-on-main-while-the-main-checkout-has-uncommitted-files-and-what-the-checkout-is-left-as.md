---
id: T-1647
type: task
nature: feature
title: An ADR decides how a batch lands on main while the main checkout has uncommitted files, and what the checkout is left as
status: backlog
parent: S-0402
owner: alex
created: 2026-10-16T05:01:13Z
updated: 2026-10-16T05:01:13Z
transitions: []
stream: S-0402
tags: [cli]
touches: [design/adrs]
---
# T-1647 An ADR decides how a batch lands on main while the main checkout has uncommitted files, and what the checkout is left as

## Work

- Write the ADR with `flai adr new`, refining ADR-0032 and the acceptance's uncommitted-files rule. Weigh at least these:
  - landing through the batch's worktree and moving main's branch with git's own check, which refuses only a file the change would overwrite
  - updating the checkout's files around the dirty ones
  - stashing and restoring them.
- Decide what the operator is told about a checkout left behind, and where.
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, indexed, and `flai check --strict` passes.

## Notes
