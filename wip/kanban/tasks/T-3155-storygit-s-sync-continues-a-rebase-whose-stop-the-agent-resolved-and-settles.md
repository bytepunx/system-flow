---
id: T-3155
type: task
nature: improvement
title: storygit's sync continues a rebase whose stop the agent resolved and settles the stops after it, with a real-git test of I-0351
status: backlog
parent: S-0726
owner: alex
created: 2026-11-12T15:56:30Z
updated: 2026-11-12T15:56:30Z
transitions: []
stream: S-0726
tags: [flai]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
after: [T-3154]
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
