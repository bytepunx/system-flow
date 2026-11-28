---
id: T-3628
type: task
nature: improvement
title: Every rebase flai runs settles a stop on a Go file with the Go merge, and the trial merge leaves such a pair out
status: backlog
parent: S-0857
owner: alex
created: 2026-11-27T23:58:18Z
updated: 2026-11-27T23:58:18Z
transitions: []
stream: S-0857
tags: [go, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/storygit/trial_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go]
after: [T-3625]
---
# T-3628 Every rebase flai runs settles a stop on a Go file with the Go merge, and the trial merge leaves such a pair out

## Work

Give `storygit.Files` a Go kind beside `Shared`, taken by `resolves`, the merge at the stop, and `Unresolved`, and build it in `issues.FilesIn` for every `.go` path, so that `flai stream sync`, `flai task done`'s sync, `flai accept`, the merge queue's batch, the rebase and resolve job, and the trial merge all share it. At the stop it reads the three stages, as `mergeStopped` does, merges them with `MergeGo` from T-3625, which it waits for because it calls it, writes and stages the result, and continues the rebase. When `MergeGo` refuses, the stop waits as today. A path that `Merged` or `Shared` takes keeps its own merge.

## Done when

- A real-git test in `sync_test.go` rebases a branch that appended a test function after the same function main's commit appended one to, and the rebase ends clean with both functions in the file.
- A real-git test shows a stop on a Go file where a side edited an existing function still waits.
- A test in `trial_test.go` shows two branches that each append a top-level declaration at one place in a `.go` file are not reported as conflicting, and `generated_test.go` shows `FilesIn` takes a `.go` path.
- `flai test` on the changed packages passes.

## Notes
