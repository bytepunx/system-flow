---
id: T-1479
type: task
nature: remediation
title: flai accept's preview lists each path the branch undoes as a blocker with its lines as data, and accept refuses it unless --keep-removal names it with a reason, which the operator alone may give
status: done
parent: S-0371
owner: alex
created: 2026-10-12T17:15:57Z
updated: 2026-10-14T22:09:01Z
transitions:
  - to: ready
    at: 2026-10-14T21:43:54Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T21:43:54Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T22:09:01Z
    by: agent-S-0371
stream: S-0371
tags: [cli, go]
touches: [flai/internal/preview/accept.go, flai/cmd/accept.go, flai/cmd/branch.go, flai/cmd/accept_undone_test.go, flai/internal/preview/orchestrator.go, flai/cmd/move.go, flai/internal/storygit/undone.go]
after: [T-1477]
usage:
  source: log
  seconds: 1507
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 87
      output: 32476
      cache_read: 5123801
      cache_write: 240161
      cost: 2.9932
---
# T-1479 flai accept's preview lists each path the branch undoes as a blocker with its lines as data, and accept refuses it unless --keep-removal names it with a reason, which the operator alone may give

## Work

It waits for T-1477, whose check and record it reads. It shares no path with T-1574, so the two run together.

- In `flai/internal/preview/accept.go`, `Acceptance` gains `UndoesMain`, JSON `undoes_main`: each path the record still lists, with main's lines the branch removes and, for a kept one, its reason. `AcceptWith` reads the record through T-1477's reader, beside `ConflictMarkers` and `StoryDiff`, and adds each path not kept to `Blockers` as a string naming the path and the two ways out, `flai stream restore S-nnnn --path <path>` and `--keep-removal <path> --reason`. Reading it in `AcceptWith` gives it to `flai accept --dry-run` and to the host read `accept.preview` alike, which calls `preview.Accept` in process, and keeps `flai/cmd/hostapi_reads_test.go`'s match between the two. `AcceptOptions` gains the kept paths and their reasons.
- In `flai/cmd/accept.go`, `addAcceptFlags` gains `--keep-removal <path>` (repeatable) and `--reason <text>`, one reason per path in order or one for all. `move.go` shares these flags, so `flai move S-nnnn done` takes them too. Refuse a `--keep-removal` that names a path the record does not hold, one with no reason, and any with `--by orchestrator`: an undone path waits for the operator. Record each kept path and its reason under `### Removals kept` in the story's `## Notes`, in the acceptance commit, with `withNotesSection` as `orchestratorNotes` does. With `--json`, a refusal on undone paths prints `{refused: {id, reason, undoes_main: [...]}}` and exits 4, as the markdown refusal in `accept_markdown.go` does, so `accept.run` gives it to the dashboard as data. The help says all of this.
- In `flai/cmd/branch.go`, `syncStoryBranch` takes the acceptance's own rebase's `UndoesMain`, and `mergeStoryBranch` stops before `git merge --ff-only` on a path not kept, after `refuseConflictMarkers`, so a path only this rebase finds is refused too, merging nothing.
- In a new `flai/cmd/accept_undone_test.go`, after the pattern of `accept_conflict_test.go`: on I-0124's sequence, the dry run lists the path under `undoes_main` and in the blockers; `flai accept` refuses, `--json` exits 4 with the path as data, main keeps its change, and the story stays in review; with `--keep-removal <path> --reason` the acceptance goes through and the Notes record the reason; with `--by orchestrator` the keep is refused; and once main's file is committed back on the branch, as T-1574's restore does, the acceptance goes through.

## Done when

- `flai accept --dry-run --json` on I-0124's sequence lists the path under `undoes_main` and in `blockers`, and `accept.preview` gives the same.
- The refusals change nothing, the real run's refusal carries the paths as data, and a kept removal is accepted with its reason in the story's Notes.
- `flai test` on the four files passes.

## Notes

`docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md` are regenerated in the docs task, once, after every help changed.
