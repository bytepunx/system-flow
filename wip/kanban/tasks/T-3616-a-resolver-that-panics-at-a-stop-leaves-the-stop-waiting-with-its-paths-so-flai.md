---
id: T-3616
type: task
nature: remediation
title: A resolver that panics at a stop leaves the stop waiting with its paths, so flai stream sync reports them instead of crashing
status: backlog
parent: S-0848
owner: alex
created: 2026-11-27T23:56:41Z
updated: 2026-11-27T23:56:41Z
transitions: []
stream: S-0848
tags: [flai, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go]
---
# T-3616 A resolver that panics at a stop leaves the stop waiting with its paths, so flai stream sync reports them instead of crashing

## Work

In I-0407 the panic hid the second stop's message, so the agent saw a crash and not the paths to resolve. `ContinueOverResolved` in `flai/internal/storygit/sync.go` calls each stop's resolvers: `files.merge` (the issue merge and the section merge) and each generated file's `Regenerate`. Today an error from one of them leaves the stop waiting for the agent. A panic must do the same.

- Recover a panic in each resolver call in `ContinueOverResolved`, and treat it as that resolver's error. Log it at warn with the path and the panic's value, and return `wait()`, so that the stop's conflicts reach `SyncResult.Conflicts` and `flai stream sync` prints them with how to continue and abort.
- Add a test in `flai/internal/storygit/sync_test.go`. A branch of two commits conflicts with main on both. The first stop is resolved and staged, as the agent did, and the second stop is on a file whose `Shared.Merge` panics. `Sync` must answer `StopConflicts` with that path, `RebaseInProgress` true, and no panic.

It waits for nothing. It shares no path with the section merge's fix, so the two run together.

## Done when

- `ContinueOverResolved` turns a resolver's panic into a stop left waiting with its conflicts, and logs the panic.
- The test in `flai/internal/storygit/sync_test.go` resumes a resolved stop into a second stop whose merge panics, and passes.
- `flai test` on the two paths passes.

## Notes
