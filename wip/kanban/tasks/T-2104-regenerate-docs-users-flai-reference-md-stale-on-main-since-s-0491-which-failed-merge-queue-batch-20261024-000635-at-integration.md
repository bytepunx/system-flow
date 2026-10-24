---
id: T-2104
type: task
nature: feature
title: Regenerate docs/users/flai-reference.md, stale on main since S-0491, which failed merge queue batch 20261024-000635 at integration
status: done
parent: S-0490
owner: alex
created: 2026-10-24T00:15:36Z
updated: 2026-10-24T00:15:54Z
transitions:
  - to: ready
    at: 2026-10-24T00:15:42Z
    by: agent-S-0490
  - to: in-progress
    at: 2026-10-24T00:15:44Z
    by: agent-S-0490
  - to: done
    at: 2026-10-24T00:15:54Z
    by: agent-S-0490
stream: S-0490
tags: []
touches: [docs/users/flai-reference.md]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 477
      cache_read: 185206
      cache_write: 3837
      cost: 0.0773
---

# T-2104 Regenerate docs/users/flai-reference.md, stale on main since S-0491, which failed merge queue batch 20261024-000635 at integration

## Work

Run `make flai-reference` in the worktree. Main lacks the `orchestration.permissions.edit_subprojects` line in `flai manifest set`'s help, which S-0491 added, so `TestReferenceIsCurrent` failed in merge queue batch 20261024-000635 (TH-0632). Recorded as I-0203. It waits for no other task: the story's other tasks are done.

## Done when

- `docs/users/flai-reference.md` matches the command help, and `flai test` on `flai/cmd/reference_test.go` passes with it.

## Notes
