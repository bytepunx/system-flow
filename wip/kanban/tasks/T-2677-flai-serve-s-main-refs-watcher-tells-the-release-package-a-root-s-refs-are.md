---
id: T-2677
type: task
nature: improvement
title: flai serve's main-refs watcher tells the release package a root's refs are unchanged, so release.pending starts no git process until a ref moves
status: done
parent: S-0604
owner: alex
created: 2026-11-07T23:50:23Z
updated: 2026-11-08T02:31:51Z
transitions:
  - to: ready
    at: 2026-11-08T02:20:45Z
    by: agent-S-0604
  - to: in-progress
    at: 2026-11-08T02:20:46Z
    by: agent-S-0604
  - to: done
    at: 2026-11-08T02:31:51Z
    by: agent-S-0604
stream: S-0604
tags: [cli]
touches: [flai/internal/release/history.go, flai/internal/release/history_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/release/release.go]
after: [T-2659]
usage:
  source: log
  seconds: 665
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 72
      output: 17621
      cache_read: 9851841
      cache_write: 356208
      cost: 4.3394
---
# T-2677 flai serve's main-refs watcher tells the release package a root's refs are unchanged, so release.pending starts no git process until a ref moves

## Work

The kept history runs `git show-ref --head --tags -d` on every call to learn whether anything moved. flai serve already watches each served project's main refs (`mainRefs` and `queuePass.start` in `flai/internal/serve/mergequeue.go`), but only to wake the merge queue's pass. That watcher leaves out `.git/refs/tags`, which the `show-ref` key reads.

Add a per-root signal to `flai/internal/release/history.go` that a watcher can set: the refs moved, or they are known unchanged since the last read. While a root is known unchanged, `PendingBatch` takes the kept key without `show-ref`. Any report of a change marks the root moved before the next read. A root with no watcher, as in `flai mcp`, `flai hostapi`, and a one-shot `flai board`, runs `show-ref` as today. Add `.git/refs/tags` to what serve watches. Have serve's watcher mark the root on every change it reports.

The watcher reports a change after one stable tick, up to about 600 ms after it. Only the board's and `publish.preview`'s reads may take the watcher's word. A publish that tags and pushes reads git as today. Say which callers skip, and why, in this task's Notes.

Waits for T-2659: it changes the per-root entry that task writes in `history.go`.

## Done when

- [ ] A test in `flai/internal/release/history_test.go` counts git processes: none for a second `Pending` on a root marked unchanged, and one `show-ref` after the root is marked moved.
- [ ] A test in `flai/internal/serve/mergequeue_test.go` moves a tag and a branch in a test repository, and shows the watcher marks the root moved.
- [ ] `go test -race ./internal/release/... ./internal/serve/...` passes.

## Notes

- Callers that take the watcher's word: `release.PendingIDs`, through `PendingBatchWatched`, which serves `board.get`'s `release.pending`, the inbox, and the orchestrator's reason figures. They only show what is pending.
- Callers that still read git: `PendingBatch` in the publish that tags and pushes (`flai/cmd/release.go`), `release.EvaluateRepo`, `preview.Publish` (`publish.preview`), and `preview.PushDryRun`. `publish.preview` is left reading git on purpose: the operator reads it right before a publish, and it gains from T-2659's per-root lock already.
- Only flai serve calls `release.Watch`, from each served project's merge queue pass, which starts for every served project. `flai mcp`, `flai hostapi`, and one-shot commands read git as before.
- The pass marks the root moved once, one watcher tick after it starts, so a move between `Watch` and the watcher's first look is not lost. A watcher count keeps a stopped pass from leaving a root trusted.
