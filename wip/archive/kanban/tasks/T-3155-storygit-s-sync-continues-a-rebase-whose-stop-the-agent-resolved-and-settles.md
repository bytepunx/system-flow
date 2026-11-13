---
id: T-3155
type: task
nature: improvement
title: storygit's sync continues a rebase whose stop the agent resolved and settles the stops after it, with a real-git test of I-0351
status: done
parent: S-0726
owner: alex
created: 2026-11-12T15:56:30Z
updated: 2026-11-12T18:57:35Z
transitions:
  - to: ready
    at: 2026-11-12T18:43:08Z
    by: agent-S-0726
  - to: in-progress
    at: 2026-11-12T18:43:09Z
    by: agent-S-0726
  - to: done
    at: 2026-11-12T18:57:35Z
    by: agent-S-0726
stream: S-0726
tags: [flai]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
after: [T-3154]
usage:
  source: log
  seconds: 866
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 17695
      cache_read: 3430379
      cache_write: 165082
      cost: 2.0051
---
# T-3155 storygit's sync continues a rebase whose stop the agent resolved and settles the stops after it, with a real-git test of I-0351

## Work

In `flai/internal/storygit/sync.go`, `Rebase` refuses any rebase in progress with `StopRebaseInProgress`. Change it to do what T-3154's ADR says:

- When the stop has no unmerged path and no staged file holds a conflict marker, continue the rebase through `ContinueOverResolved` with the same `Files`.
  - `resolvable` returns not-ok on a stop with no conflicts, so the first continue is made before the loop, or the loop is taught that case.
- Report what the loop merged and regenerated, as `continueOverGenerated` does.
- On a clean finish, run the steps after a clean rebase: `regenerateAtHead`, `settleChangelogs`, `checkUndone`, `fold`, and `fix`. `checkUndone` needs the branch's old base and head: recover them as the ADR says.
- Stop again with `StopConflicts` when a later stop is not settleable.
- Keep refusing, with the unmerged or marked paths named, a stop the agent has not resolved.
- Change the stop advice in `rebaseStopped` text (around lines 859–864) from "git add it, and run git rebase --continue; then run flai stream sync again" to "git add it, then run flai stream sync again". Keep `git rebase --abort` as the way back.

Add real-git tests in `flai/internal/storygit/sync_test.go` beside `TestSyncStopsOnConflictsAndRefusesARebaseInProgress` and `TestContinueOverResolvedContinuesARebaseInAnyWorktree`:

- I-0351's case: main and the branch both change a file the section merge cannot settle, in the branch's first commit. Both add an instance to one issue file, and so change `design/issues/summary.md`, in its second commit. The sync stops on the first. The test resolves and stages it. The second sync finishes with the issue merged by instances and the summary regenerated.
- A second sync with the stop unresolved is refused, naming the path.
- Show the first test fails on main's `sync.go`, as `git.md` says.

Waits for T-3154, whose ADR decides the behaviour and how the old base is recovered.

## Done when

- `flai stream sync` run again after the agent resolved and staged a stop continues the rebase and settles later issue and summary stops itself.
- The new tests pass with `flai test flai/internal/storygit`, and the I-0351 test fails without the change.

## Notes

Settled by ADR-0294 at the story's review of the plan:

- Only `Sync` resumes. Add `Resume bool` to `SyncOptions`; `Sync` sets it before it calls `Rebase`. `Rebase` called alone, as `flai accept` and the merge queue call it, still refuses a rebase in progress with `StopRebaseInProgress`.
- A stop is resolved when no path is unmerged, `git diff --cached --check` finds no conflict marker, nothing tracked is changed and unstaged, and the rebase state's `head-name` is `refs/heads/story/S-nnnn`. Otherwise refuse as now, naming the unmerged, marked, or unstaged paths.
- Read `onto` and `orig-head` from the rebase state (`rebase-merge/` or `rebase-apply/`, via `git rev-parse --git-path`) before the continue. Old base is `git merge-base onto orig-head`. The steps after a clean rebase run against that `onto`, not main's head now. A state missing either file is refused, naming `git rebase --abort` and a fresh sync.
- Add `Resumed bool` (`json:"resumed"`) to `SyncResult`.
