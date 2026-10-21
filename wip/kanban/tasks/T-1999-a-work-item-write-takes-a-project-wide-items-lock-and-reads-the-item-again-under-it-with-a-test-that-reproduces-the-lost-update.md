---
id: T-1999
type: task
nature: remediation
title: A work item write takes a project-wide items lock and reads the item again under it, with a test that reproduces the lost update
status: backlog
parent: S-0479
owner: alex
created: 2026-10-21T12:10:12Z
updated: 2026-10-21T12:10:12Z
transitions: []
stream: S-0479
tags: [flai, workitem]
touches: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/repo.go]
---
# T-1999 A work item write takes a project-wide items lock and reads the item again under it, with a test that reproduces the lost update

## Work

Add the primitive every item writer will use. `Repo.Save` writes the whole file from a copy the caller read earlier, with no lock, so a writer in `flai serve` (usage, the replanner's forecast) and one in the MCP process (`story_start`) overwrite each other (I-0175, I-0176).

- Add `flai/internal/workitem/itemlock.go` with a `Repo` method, such as `Repo.Update(fn)`, that takes `filelock.Lock` on one lock file under `Repo.CacheDir()`, such as `items.lock`, and runs `fn` while holding it. `flai/internal/mergequeue/state.go` `Update` is the pattern to follow. Use one project-wide lock rather than one per item, because a transition writes the item, its cascade, its epic, and the board in one step. Use a bounded wait, such as 30 s, and an error that names the lock file.
- Inside the lock, `fn` must read each item fresh from disk. The `parsed` store in `flai/internal/workitem/store.go` is validated by mtime and size, so a write in the same mtime tick with the same size could be served stale. Give `fn` a read that bypasses the cache, or check that the store cannot return a stale copy here.
- `filelock` excludes two opens of the same file in one process too, so a nested `Update` would wait on itself until the timeout. Document that `Save` takes no lock and is called inside `Update`, and make a nested call fail fast with a clear error rather than wait.
- In `flai/internal/workitem/itemlock_test.go`, reproduce the lost update. Read a story, move a fresh copy to in-progress and save it, then save the first copy with only its usage changed. Assert that before the fix the transition is gone, and that through `Update` both changes stay. Add a concurrent test, two goroutines or a subprocess, after `TestConcurrentUpdatesLoseNoWrite` in `flai/internal/mergequeue/state_test.go`, that writes different fields of one item many times and loses none.

It waits for nothing: it is the first layer.

## Done when

- `Repo.Update`, or the name chosen, holds the items lock across a fresh read, a change, and the save. It times out with a named error, and a nested call fails fast.
- `flai/internal/workitem/itemlock_test.go` shows the lost update with the plain read-modify-write, and shows that no write is lost through the lock, under `go test -race`.
- `flai test flai/internal/workitem` passes.

## Notes

Drafted by the planner. The writers move onto the lock in the next layer's two tasks.
