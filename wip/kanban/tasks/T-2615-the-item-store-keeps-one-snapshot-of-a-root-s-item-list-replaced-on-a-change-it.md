---
id: T-2615
type: task
nature: improvement
title: The item store keeps one snapshot of a root's item list, replaced on a change it is told of or on the process's own save, and lists and gets from it without statting
status: backlog
parent: S-0606
owner: alex
created: 2026-11-07T22:48:59Z
updated: 2026-11-07T22:48:59Z
transitions: []
stream: S-0606
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 12
          cache_read: 62745
          cache_write: 12650
          cost: 0.0192
---
# T-2615 The item store keeps one snapshot of a root's item list, replaced on a change it is told of or on the process's own save, and lists and gets from it without statting

## Work

Give the process-wide store in `flai/internal/workitem/store.go` a snapshot mode per repository root. A caller that watches the root, `flai serve`, turns it on and tells the store of each changed path under the item folders. While it is on, the root keeps one parsed list of every item, the archive included. `Repo.List` and `Repo.Get` in `flai/internal/workitem/repo.go` answer from it with no `os.ReadDir`, glob, or stat. A change the store is told of replaces the snapshot, or the entry for that path. `Repo.Save` in this process replaces it too, so a write is seen at once by the next read. Each caller still gets its own deep copy. No folder lock is held across a stat loop while the snapshot answers. A root with snapshot mode off, as in every one-shot `flai` and in `flai mcp`, keeps today's stat-and-compare path unchanged.

Waits for no task. It shares no path with the watcher task, so the two run together.

## Done when

- [ ] Tests in `flai/internal/workitem/store_test.go` show that with snapshot mode on, a list and a get read no file and stat nothing after the first list, that a change the store is told of and a `Repo.Save` are seen by the next list, and that callers still get their own copies.
- [ ] With snapshot mode off, the existing store tests pass unchanged.

## Notes
