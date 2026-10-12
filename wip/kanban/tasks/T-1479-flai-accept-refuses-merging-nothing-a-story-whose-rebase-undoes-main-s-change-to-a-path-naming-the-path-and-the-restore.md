---
id: T-1479
type: task
nature: remediation
title: flai accept refuses, merging nothing, a story whose rebase undoes main's change to a path, naming the path and the restore
status: backlog
parent: S-0371
owner: alex
created: 2026-10-12T17:15:57Z
updated: 2026-10-12T17:15:57Z
transitions: []
stream: S-0371
tags: [cli, go]
touches: [flai/cmd/branch.go, flai/cmd/accept.go, flai/cmd/accept_undone_test.go]
after: [T-1477]
---
# T-1479 flai accept refuses, merging nothing, a story whose rebase undoes main's change to a path, naming the path and the restore

## Work

It waits for T-1477, which adds `SyncResult.UndoesMain` to the result of `storygit.Rebase`. It shares no path with T-1478, so the two run together.

- In `flai/cmd/branch.go`, `syncStoryBranch` returns an error when the rebase's `UndoesMain` is not empty, and `mergeStoryBranch` stops before the fast-forward, as `refuseConflictMarkers` does. The refusal names each path, the lines main gained that the branch removes, and the restore: in the worktree, `git checkout <main> -- <path>`, commit, and accept again. The branch stays rebased in its worktree and nothing reaches main.
- In `flai/cmd/accept.go`, the help says `flai accept` refuses such a branch.
- In a new `flai/cmd/accept_undone_test.go`, after the pattern of `accept_conflict_test.go`: a story in review whose branch made and reverted a change that main then made too is refused by `flai accept`, main keeps its change, and the story stays in review; once the path is restored and committed, the acceptance goes through.

## Done when

- The two tests pass.
- `flai test` on the three files passes.

## Notes

The dry run and the dashboard's Accept preview do not rebase, so they cannot name the path beforehand; the refusal comes at the real acceptance, after the rebase, as the second conflict-marker check does.
