---
id: T-1648
type: task
nature: feature
title: storygit tells which uncommitted files in the main checkout a change between two commits would overwrite
status: backlog
parent: S-0402
owner: alex
created: 2026-10-16T05:01:17Z
updated: 2026-10-16T05:01:17Z
transitions: []
stream: S-0402
tags: [cli]
touches: [flai/internal/storygit/storygit.go, flai/internal/storygit/storygit_test.go]
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
