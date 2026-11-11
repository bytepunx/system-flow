---
id: T-2965
type: task
nature: remediation
title: Reproduce I-0338 and I-0335 in real-git sync tests where main gains a commit between the rebase and the checks after it
status: in-progress
parent: S-0712
owner: alex
created: 2026-11-10T10:28:03Z
updated: 2026-11-11T22:12:59Z
transitions:
  - to: ready
    at: 2026-11-11T22:12:59Z
    by: agent-S-0712
  - to: in-progress
    at: 2026-11-11T22:12:59Z
    by: agent-S-0712
stream: S-0712
tags: [flai, git, test]
touches: [flai/internal/storygit/sync_test.go]
after: [T-2963]
---
# T-2965 Reproduce I-0338 and I-0335 in real-git sync tests where main gains a commit between the rebase and the checks after it

## Work

Add real-git tests to `flai/internal/storygit/sync_test.go`, next to `TestSyncFindsNothingUndoneOnABranchThatKeepsMainsChanges` and `TestSyncPutsTheBranchsChangelogLinesUnderMainsTopUnreleased`, using the `gitRepo`, `openStory`, `syncOptions`, `commitFiles`, and `releasedChangelog` helpers. They reproduce the instances in I-0338 and I-0335.

- Give the story branch a commit to a path of its own. Give main a commit since the branch's base, so the sync has something to rebase onto.
- Wrap the sync's `execx.Runner`, as `trialCounter` in `trial_test.go` does. Right after the `git rebase` call returns, the wrapper commits in the main checkout to paths the branch never changes. This is main gaining commits while the sync runs.
- Undone case (I-0338, and I-0335's two commits): the wrapper commits twice to a path the branch never changes. Assert that the sync answers `undoes_main` empty and leaves no record under `.flai-cache/undone/`.
- Changelog case: the branch adds a line to `CHANGELOG.md`, and the wrapper's commit adds another story's line to main's `CHANGELOG.md`. Assert that the settled changelog holds the lines of main as the rebase took it and the branch's line, not the line main gained after the rebase.
- Run each test against the code before T-2963's change, by hand or by stashing, and say in this task's Notes whether it fails there. The undone case must fail there, to show it reproduces the cause.

It waits for T-2963: the tests pass only with that fix, and they call `Rebase` as T-2963 leaves it.

## Done when

- The undone case fails without T-2963's change and passes with it, and the changelog case passes with it. Its Notes say how each ran without the change.
- The tests pass without `-short`, where `gitRepo` does not skip them: `go test -run <their names> ./internal/storygit` in `flai/`. The close-out's integration tier runs them too.
- `flai test flai/internal/storygit` passes.

## Notes
