---
id: T-2965
type: task
nature: remediation
title: Reproduce I-0338 in a real-git sync test where main gains a commit between the rebase and the undone check
status: backlog
parent: S-0712
owner: alex
created: 2026-11-10T10:28:03Z
updated: 2026-11-10T10:28:13Z
transitions: []
stream: S-0712
tags: [flai, git, test]
touches: [flai/internal/storygit/sync_test.go]
after: [T-2963]
---
# T-2965 Reproduce I-0338 in a real-git sync test where main gains a commit between the rebase and the undone check

## Work

Add a real-git test to `flai/internal/storygit/sync_test.go`, next to `TestSyncFindsNothingUndoneOnABranchThatKeepsMainsChanges`, using the `gitRepo`, `openStory`, `syncOptions`, and `commitFiles` helpers. The test reproduces the instance in I-0338.

- Give the story branch a commit to a path of its own. Give main a commit since the branch's base, so the sync has something to rebase onto.
- Wrap the sync's `execx.Runner`, as `trialCounter` in `trial_test.go` does. Right after the `git rebase` call returns, the wrapper commits once, in the main checkout, to a path the branch never changes. This is main gaining a commit while the sync runs.
- Assert that the sync answers `undoes_main` empty and leaves no record under `.flai-cache/undone/`.
- Run the test against the code before T-2963's change, by hand or by stashing, and say in this task's Notes that it fails there. That shows the test reproduces the cause.

It waits for T-2963: the test passes only with that fix, and it calls `Rebase` as T-2963 leaves it.

## Done when

- The new test fails without T-2963's change and passes with it. Its Notes say so.
- The test passes without `-short`, where `gitRepo` does not skip it: `go test -run <its name> ./internal/storygit` in `flai/`. The close-out's integration tier runs it too.
- `flai test flai/internal/storygit` passes.

## Notes
