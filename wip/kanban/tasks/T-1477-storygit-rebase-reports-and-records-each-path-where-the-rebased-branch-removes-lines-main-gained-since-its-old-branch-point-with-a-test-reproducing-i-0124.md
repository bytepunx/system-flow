---
id: T-1477
type: task
nature: remediation
title: storygit.Rebase reports and records each path where the rebased branch removes lines main gained since its old branch point, with a test reproducing I-0124
status: done
parent: S-0371
owner: alex
created: 2026-10-12T17:15:30Z
updated: 2026-10-14T21:43:43Z
transitions:
  - to: ready
    at: 2026-10-14T21:36:03Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T21:36:04Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T21:43:43Z
    by: agent-S-0371
stream: S-0371
tags: [cli, go]
touches: [flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
after: [T-1476]
usage:
  source: log
  seconds: 459
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 50
      output: 256
      cache_read: 2546753
      cache_write: 125400
      cost: 1.2318
---
# T-1477 storygit.Rebase reports and records each path where the rebased branch removes lines main gained since its old branch point, with a test reproducing I-0124

## Work

It waits for T-1476, so that the code builds the decision the ADR records.

- In `flai/internal/storygit/sync.go`, `Rebase` reads `git merge-base <main> HEAD` in the worktree before `git rebase`, and after a clean rebase, when the old base is not main's head, calls the new check. Add `UndoesMain []Undone` to `SyncResult`, JSON `undoes_main`, initialised empty as the other lists are.
- In a new `flai/internal/storygit/undone.go`, the check: for each path in both `git diff --name-only <oldbase> <main>` and `git diff --name-only <main> HEAD`, read the lines main added (`git diff <oldbase> <main> -- <path>`) and the lines the branch's net change removes (`git diff <main> HEAD -- <path>`), and report the path when the second holds any line of the first, counted as a multiset of line text. A path the branch deletes counts every line main added. Each `Undone` gives `path`, `lines` (main's lines the branch removes, as text, for the dashboard's card), and `old_base`. A git error is returned and logged by `Sync` as the other checks' are.
- In the same file, the record the ADR names under `.flai-cache/`, one per story: `Sync` adds each new finding to it, and a reader returns the recorded paths still undone, re-checking each against main now, so that a restored path drops out. T-1479's preview and T-1574's restore read it.
- In `flai/internal/storygit/undone_test.go`, unit tests of the line comparison and the record: a removal of main's added lines is reported, an unrelated removal is not, a branch that keeps main's change is not, and a recorded path the branch no longer removes drops out.
- In `flai/internal/storygit/sync_test.go`, reproduce I-0124 with the existing helpers (`gitRepo`, `openStory`, `commitFiles`, `syncOptions`): the branch changes a file and then reverts it to main of that time, main then makes the same change, and `Sync` reports the path in `UndoesMain` after a clean rebase and records it. A second test syncs a branch with no such path and gets none.

## Done when

- The I-0124 reproduction passes and fails with the check taken out.
- A second sync after the first still finds the path through the record.
- `flai test` on the four files passes: gofmt, vet, golangci-lint, and the Go tests.

## Notes

Keep the change in `sync.go` to the base read, the call, the record, and the field.
