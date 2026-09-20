---
id: T-0274
type: task
nature: feature
title: Try every write end to end without a mount, including an acceptance from the board, and document it
status: done
parent: S-0075
owner: alex
created: 2026-09-20T08:49:15Z
updated: 2026-09-20T09:10:29Z
transitions:
  - to: ready
    at: 2026-09-20T09:03:22Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T09:03:22Z
    by: system-flow
  - to: done
    at: 2026-09-20T09:10:29Z
    by: system-flow
stream: S-0075
tags: []
---
# T-0274 Try every write end to end without a mount, including an acceptance from the board, and document it

## Work
A scratch project with a story branch in a worktree, the new image with no project mount, flai serve from this tree: moves, order, block, new items, threads, a document save with a conflict, an ADR, stats, and an acceptance streamed to the board; a connection lost during an acceptance. Design documents, users' and operators' documentation.

## Done when
- What was tried is in the narrative
- All three suites and flai check --strict pass

## Notes
