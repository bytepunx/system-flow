---
id: T-1512
type: task
nature: remediation
title: The guard's sub-agent record and flai serve's agents state lock through filelock, with no age break
status: done
parent: S-0377
owner: alex
created: 2026-10-13T16:48:18Z
updated: 2026-10-13T17:18:12Z
transitions:
  - to: ready
    at: 2026-10-13T16:53:46Z
    by: agent-S-0377
  - to: in-progress
    at: 2026-10-13T16:53:49Z
    by: agent-S-0377
  - to: done
    at: 2026-10-13T17:06:38Z
    by: agent-S-0377
stream: S-0377
tags: [flai, concurrency, guard, serve]
touches: [flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/serve/restart.go, flai/internal/serve/restart_test.go, flai/cmd/guard_test.go, flai/cmd/promote_test.go]
after: [T-1510]
usage:
  source: log
  seconds: 769
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 10933
      cache_read: 1541652
      cache_write: 73772
      cost: 1.0067
---
# T-1512 The guard's sub-agent record and flai serve's agents state lock through filelock, with no age break

## Work

`lockRecord` in `flai/internal/guard/subagents.go` and `lockFile` in `flai/internal/serve/restart.go` break a stale lock with the same `os.Stat` then `os.Remove` that I-0140 found in `lockActivity`. Two guards, or two flai serve writers, can both take the lock, and one write of the record or the agents state is lost. Move both to `flai/internal/filelock`.

- `lockRecord` keeps its wait, `recordLockWait`, and its error when the wait passes. Remove `recordLockAge`, and drop "remove it if none runs" from the error, since a lock no live process holds is free.
- `lockFile` keeps its 5 s wait and, when the wait passes, writes without the lock, as it does now and for the reason its comment gives. Remove `staleLock`.

Add a test to each package where none covers the lock: writers in goroutines, more than one at a time, each changing the record or the state, lose nothing, under `-race`.

It waits for T-1510, whose package it uses. It shares no file with T-1511, so the two run together.

## Done when

- Neither `flai/internal/guard` nor `flai/internal/serve` breaks a lock by its age or removes a lock file.
- The new tests pass with `-race -count=20`.
- `flai test` on the changed paths passes.

## Notes
