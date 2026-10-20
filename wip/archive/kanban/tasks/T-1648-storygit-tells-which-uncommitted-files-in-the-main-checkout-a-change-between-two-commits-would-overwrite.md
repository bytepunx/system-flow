---
id: T-1648
type: task
nature: feature
title: storygit tells which uncommitted files in the main checkout a change between two commits would overwrite
status: done
parent: S-0402
owner: alex
created: 2026-10-16T05:01:17Z
updated: 2026-10-19T23:59:06Z
transitions:
  - to: ready
    at: 2026-10-19T23:57:25Z
    by: agent-S-0402
  - to: in-progress
    at: 2026-10-19T23:57:26Z
    by: agent-S-0402
  - to: done
    at: 2026-10-19T23:59:06Z
    by: agent-S-0402
stream: S-0402
tags: [cli]
touches: [flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go]
usage:
  source: log
  seconds: 100
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 6207
      cache_read: 1701466
      cache_write: 34619
      cost: 0.7203
---
# T-1648 storygit tells which uncommitted files in the main checkout a change between two commits would overwrite

## Work

- Beside `DirtyOutsideWip`, add a function that, given main's head and the commit a landing would move it to, lists the uncommitted and untracked files in the main checkout that the change touches.
- Leave the rest of the dirty files out: they stop nothing.
- Test a dirty file the change touches, one it does not, and an untracked file the change adds.
- First in its layer: it waits for no task.

## Done when

The function lists exactly the dirty files the change would overwrite in each tested case, and the tests pass.

## Notes
