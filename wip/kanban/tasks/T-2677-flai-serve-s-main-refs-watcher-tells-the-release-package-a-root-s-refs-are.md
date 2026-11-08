---
id: T-2677
type: task
nature: improvement
title: flai serve's main-refs watcher tells the release package a root's refs are unchanged, so release.pending starts no git process until a ref moves
status: backlog
parent: S-0604
owner: alex
created: 2026-11-07T23:50:23Z
updated: 2026-11-07T23:50:23Z
transitions: []
stream: S-0604
tags: [cli]
touches: [flai/internal/release/history.go, flai/internal/release/history_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go]
after: [T-2659]
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
