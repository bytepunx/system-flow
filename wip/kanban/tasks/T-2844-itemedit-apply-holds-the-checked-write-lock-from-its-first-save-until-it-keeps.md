---
id: T-2844
type: task
nature: improvement
title: itemedit.Apply holds the checked-write lock from its first save until it keeps or restores the edit
status: done
parent: S-0654
owner: alex
created: 2026-11-10T07:01:09Z
updated: 2026-11-10T07:44:03Z
transitions:
  - to: ready
    at: 2026-11-10T07:35:02Z
    by: agent-S-0654
  - to: in-progress
    at: 2026-11-10T07:35:03Z
    by: agent-S-0654
  - to: done
    at: 2026-11-10T07:44:03Z
    by: agent-S-0654
stream: S-0654
tags: [flai, itemedit, locking]
touches: [flai/internal/itemedit/itemedit.go, flai/internal/itemedit/lock_test.go, flai/internal/workitem/itemlock.go]
usage:
  source: log
  seconds: 540
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 41
      output: 269
      cache_read: 4723374
      cache_write: 316593
      cost: 1.8542
---
# T-2844 itemedit.Apply holds the checked-write lock from its first save until it keeps or restores the edit

## Work

In `flai/internal/itemedit/itemedit.go`, `Apply` take `repo.LockCheckedWrites()` after the check before the write (`runCheck(repo, opt.Now, write)`) and just before the first `repo.Save(it)`, as `itemnew.Create` does in `flai/internal/itemnew/itemnew.go` (ADR-0272). Hold it through the parents' saves, the narrative and relink writes, the index write, the check after the write, and any `u.restore()`. Let go once the check has kept the edit or the restore has put the files back, on every return path, `fail` included, with a release that runs once. A lock it cannot take refuses the edit with nothing written, saying to try again once the writer holding it is done.

Keep the lock order: the checked-write lock is taken first, and `Save` and `Restore` take the items lock inside it. No caller of `Apply` (`flai/cmd/edit.go`, `flai/cmd/criteria.go`, `flai/internal/mcpserver/items_write.go`, `flai/internal/mcpserver/criteria.go`) holds the checked-write lock, and nothing inside the window takes it; confirm both, since a second flock on the same file from one process waits on itself until the 30 s timeout.

Update the comment on `LockCheckedWrites` in `flai/internal/workitem/itemlock.go` to say a checked edit holds it too.

Add a test to `flai/internal/itemedit/lock_test.go` that reproduces I-0328: replace `runCheck` so that its run after the write tries `filelock.Lock` on `.flai-cache/items-checked.lock` with no wait and records that it could not take it, then returns an error finding on the item's file. Assert the probe found the lock held, the edit was refused and restored, and the lock is free once `Apply` returns. A second case: an edit the check keeps also frees the lock. `flai/internal/serve/hold_test.go` probes the same lock and is a model.

It waits for no task: no other task of S-0654 changes these files.

## Done when

- `Apply` holds the checked-write lock from its first save until it keeps or restores the edit, and lets go on every return.
- The new tests pass, and `go test ./internal/itemedit/... ./internal/workitem/...` passes from `flai/`.

## Notes
