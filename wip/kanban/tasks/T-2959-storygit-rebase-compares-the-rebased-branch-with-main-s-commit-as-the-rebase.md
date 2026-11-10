---
id: T-2959
type: task
nature: remediation
title: storygit.Rebase compares the rebased branch with main's commit as the rebase took it, not with main's name read again after it
status: backlog
parent: S-0709
owner: alex
created: 2026-11-10T10:27:23Z
updated: 2026-11-10T10:27:23Z
transitions: []
stream: S-0709
tags: [cli]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/sync_test.go, flai/internal/storygit/undone_test.go]
---
# T-2959 storygit.Rebase compares the rebased branch with main's commit as the rebase took it, not with main's name read again after it

## Work

First, in `flai/internal/storygit/sync_test.go`, write a real-git test that reproduces I-0335: a story branch and a main branch in a temporary repository, and a runner that wraps `execx.System` and, right after it runs the `git rebase` call, commits to main a change to a path the branch never touched. Run `Rebase` through it and see `UndoesMain` name that path, as the instances on S-0680 and S-0682 did.

Then fix `Rebase` in `flai/internal/storygit/sync.go`: read main's commit once with `git rev-parse <main>` before the rebase, rebase onto that commit, and pass that commit, not the branch name, to `checkUndone`'s `FindUndone` and to `SettleChangelogs`, the two places that compare the rebased branch with main after the rebase. Keep `res.Base` the branch name, so the sync's output and `fix`, `fold`, and `checkSync`, whose comparisons go through the merge base or read main's files by name, do not change. Update `FindUndone`'s comment in `flai/internal/storygit/undone.go` to say its `main` is the commit the rebase went onto, and its scripted tests in `flai/internal/storygit/undone_test.go` only if the argument's change breaks them.

Waits for nothing: it is the first layer.

## Done when

- The new test fails on the code before the fix and passes after it: a commit to main between the rebase and the check leaves `UndoesMain` empty for a path the branch never changed.
- `TestSyncFindsNothingUndoneOnABranchThatKeepsMainsChanges`, `TestFindUndone`, and the rest of `flai test` on the four paths pass.
- The sync's JSON output names main by its branch name in `base`, as before.

## Notes
