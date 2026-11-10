---
id: T-2827
type: task
nature: remediation
title: Repo.Save refuses to recreate an item file moved since it was read, and Archive moves under the items lock
status: done
parent: S-0648
owner: alex
created: 2026-11-10T06:37:17Z
updated: 2026-11-10T08:27:10Z
transitions:
  - to: ready
    at: 2026-11-10T08:09:47Z
    by: agent-S-0648
  - to: in-progress
    at: 2026-11-10T08:09:49Z
    by: agent-S-0648
  - to: done
    at: 2026-11-10T08:27:10Z
    by: agent-S-0648
stream: S-0648
tags: [flai, workitem]
touches: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go, design/system/work-hierarchy.md, design/adrs/0277-repo-save-refuses-to-write-an-active-work-item-whose-file-was-archived-since-it.md, design/adrs/README.md]
usage:
  source: log
  seconds: 26
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 10
      output: 87
      cache_read: 1331630
      cache_write: 5514
      cost: 0.492
---
# T-2827 Repo.Save refuses to recreate an item file moved since it was read, and Archive moves under the items lock

## Work

First layer: it waits for no task.

- Reproduce first. In `itemlock_test.go`, read an item with `Repo.Get`, archive it with `Archive`, then `Save` the copy read before. Today `mergeSaved` finds no file at the old path, returns the whole item, and `saveLocked` writes it there again. So the test finds the item both in `wip/kanban` and in `wip/archive`, the `item.duplicate-id` of I-0323. Record that failing run in the task's Notes.
- Fix `Save`. In `saveLocked` and `mergeSaved` (`itemlock.go`), when the item was read from a file, its `readAs` is set, and that file is gone, refuse the write with an error of its own, such as `ErrItemMoved`, that names the item and its old path. Write nothing. An item never read from a file, a new one, saves as it does today.
- Fix `Archive` (`archive.go`). Take `lockItems` for its moves, so a `Save` waiting on the lock reads the item after the move. Check each caller, `flai/cmd/accept.go` and `flai/cmd/archive.go`, so that none holds the lock when it calls `Archive`. The lock is not re-entrant, and a held lock would deadlock it. If one does, change that caller, and say so in the narrative's Decisions.
- Add a test in `archive_test.go` that a `Save` started while `Archive` holds the lock is refused and leaves no file at the old path.
- Say in `design/system/work-hierarchy.md`, in the bullet on writers that run at once (ADR-0184), that a write never recreates an item file moved or removed since it was read, and that an archive moves under the same lock.

## Done when

- The reproducing test fails before the fix, and the run is in the task's Notes. It passes after the fix.
- `go test ./internal/workitem -run 'Save|Archive' -count 20` passes in `flai/`.
- `flai test flai/internal/workitem/itemlock.go flai/internal/workitem/archive.go design/system/work-hierarchy.md` passes.

## Notes
