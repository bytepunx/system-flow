---
id: T-1510
type: task
nature: remediation
title: flai/internal/filelock takes an OS file lock that its holder's death frees, on unix and Windows
status: done
parent: S-0377
owner: alex
created: 2026-10-13T16:48:02Z
updated: 2026-10-13T16:53:39Z
transitions:
  - to: ready
    at: 2026-10-13T16:49:57Z
    by: agent-S-0377
  - to: in-progress
    at: 2026-10-13T16:49:59Z
    by: agent-S-0377
  - to: done
    at: 2026-10-13T16:53:39Z
    by: agent-S-0377
stream: S-0377
tags: [flai, concurrency]
touches: [flai/internal/filelock/filelock.go, flai/internal/filelock/filelock_unix.go, flai/internal/filelock/filelock_windows.go, flai/internal/filelock/filelock_test.go, flai/go.mod, design/tech/go-libraries.md]
usage:
  source: log
  seconds: 220
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 25
      output: 8252
      cache_read: 1163671
      cache_write: 55684
      cost: 0.7599
---
# T-1510 flai/internal/filelock takes an OS file lock that its holder's death frees, on unix and Windows

## Work

Add the package `flai/internal/filelock`. It takes an exclusive lock on a file made for the purpose. It opens the file, creating it and its folder when missing, and locks the open file: `syscall.Flock` with `LOCK_EX` on unix (`filelock_unix.go`, build tag `unix`), and `LockFileEx` with `LOCKFILE_EXCLUSIVE_LOCK` from `golang.org/x/sys/windows` on Windows (`filelock_windows.go`). The returned unlock releases the lock and closes the file, and leaves the file in place. Removing it would let a waiter that already opened it lock an unlinked file while a new writer locks a fresh one.

It waits while another holds the lock, up to a wait the caller gives; no wait means it waits until the lock is free. It never breaks a lock by its age: the kernel frees the lock when the holding process dies. A blocking `flock` in a goroutine, or a poll with `LOCK_NB` every 10 ms, are both fine; say which and why in the package comment.

Move `golang.org/x/sys` from indirect to direct in `flai/go.mod`, and add its row to `design/tech/go-libraries.md`, as the documentation convention asks for a dependency added.

This task waits for none: the call sites in the next layer use the package.

## Done when

- `filelock_test.go` shows that, of many goroutines that each lock, read a counter file, and write it back incremented, none loses an increment, run with `-race` and `-count=20`.
- A test starts a child process (the test binary, re-run with an environment variable) that takes the lock and is killed while holding it, and the parent then takes the lock at once, without any wait for an age.
- A lock not freed within the caller's wait returns an error naming the file.
- `GOOS=windows go build ./...` and `GOOS=darwin go build ./...` in `flai/` succeed.
- `flai test` on the changed paths passes.

## Notes
