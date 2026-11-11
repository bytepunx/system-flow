---
id: T-3030
type: task
nature: remediation
title: A sync whose rebase stops with no conflict says git's error and offers no abort that cannot restore the branch, and a trial merge that cannot run is reported as not run
status: done
parent: S-0689
owner: alex
created: 2026-11-10T10:33:29Z
updated: 2026-11-11T16:00:32Z
transitions:
  - to: ready
    at: 2026-11-11T15:27:35Z
    by: agent-S-0689
  - to: in-progress
    at: 2026-11-11T15:27:35Z
    by: agent-S-0689
  - to: done
    at: 2026-11-11T16:00:32Z
    by: agent-S-0689
stream: S-0689
tags: [flai, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/branch.go, flai/cmd/branch_test.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-3026]
usage:
  source: log
  seconds: 1977
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 114
      output: 42423
      cache_read: 10130062
      cache_write: 471216
      cost: 5.8286
---
# T-3030 A sync whose rebase stops with no conflict says git's error and offers no abort that cannot restore the branch, and a trial merge that cannot run is reported as not run

## Work

Wherever a sync runs, two of its reports mislead when git cannot write (I-0272).

1. A rebase that stops with no conflict. `storygit.Rebase` (`flai/internal/storygit/sync.go`, near line 440) discards `git rebase`'s error when a rebase is left in progress and classes the stop as conflicts. With no conflicted path, the stop then reads "resolve them and git rebase --continue, or git rebase --abort" (`flai/cmd/branch.go`, `headline` and `Error`). Its `Abort` text (`stopped`, near line 642) says the abort puts the branch back, which is false when the ref cannot be locked: the abort exits 0 and leaves HEAD detached. Keep git's error on the stop. A stop with no conflicted path reports that error as what stopped the rebase, and gives no "resolve them". Before it offers `git rebase --abort`, check that the branch's ref can be updated, as `AbortRebase` checks HEAD after an abort. When it cannot, say so and name the safe way back: `git symbolic-ref HEAD refs/heads/<branch>` with the branch left at its head.
2. A trial merge that cannot run. `checkSync` (near line 822) stops on a `TrialMerge` error other than exit 1. `Sync` (near line 378) then only logs "story branch checks failed after the rebase" and the command prints "is rebased onto" with no trial-merge result. Record that the trial merge did not run, and why, as the result already records `TrialMergeSkipped` for an old git. Have `printSyncChecks` (`flai/cmd/stream_sync.go`) print it, so a sync never reads clean with its trial merge unrun.

Reproduce both in `flai/internal/storygit/sync_test.go` on a common dir made read-only, following `TestGitWorksReadOnlyOnTheMountsAndRefusesACommit` in `mounts_test.go`, which skips under root. Rename `testSyncConflictListsPathsAndHowToContinueOrAbort` in `stream_sync_test.go` to `Test...` so that Go runs it, and fix what it then finds.

It waits for T-3026 because both change `flai/cmd/stream_sync_test.go`.

## Done when

- On a read-only common dir, a test's sync stops with git's own error named and no "resolve them" or false abort promise, and the branch is still at its head.
- On a read-only common dir, a test's trial merge is reported as not run, with its reason, in the result and in `flai stream sync`'s output.
- A real conflict still lists its paths with the continue and abort lines as before.
- The renamed conflict test runs and passes.
- `flai test` passes on the changed paths.

## Notes
