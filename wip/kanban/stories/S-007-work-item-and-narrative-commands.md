---
id: S-007
type: story
nature: feature
title: Work item and narrative commands
status: done
parent: E-002
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-15T18:00:57Z
transitions:
  - to: ready
    at: 2026-09-15T17:42:41Z
    by: agent
  - to: in-progress
    at: 2026-09-15T17:42:41Z
    by: agent
  - to: review
    at: 2026-09-15T17:52:37Z
    by: agent
  - to: done
    at: 2026-09-15T18:00:57Z
    by: alex
tags: []
---

# S-007 Work item and narrative commands

## Goal
epic new, story new, task new, move, block, unblock, board, stream open, stream log, and archive, all validated against the schema.

## Acceptance criteria
- [x] IDs allocate correctly across kanban and archive
- [x] move enforces transition rules and parent-child constraints
- [x] stream open and log maintain index.md
- [x] archive moves items and narratives together

## Tasks
- T-017 Work item package: parse, write, list, IDs
- T-018 Transition rules, block, unblock, move
- T-019 epic new, story new, task new
- T-020 board and show
- T-021 stream open and stream log, index maintenance
- T-022 archive
- T-023 Docs and design updates

## Notes
- Also delivered `flai show` and `flai board --all`.
- Refinement: tasks may skip review (in-progress to done). See workflow.md.
- The round-trip test found and fixed an unquoted colon in T-017's title.
