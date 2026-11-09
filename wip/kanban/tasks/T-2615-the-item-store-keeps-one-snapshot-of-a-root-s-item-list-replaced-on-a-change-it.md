---
id: T-2615
type: task
nature: improvement
title: The item store keeps one snapshot of a root's item list, replaced on a change it is told of or on the process's own save, and lists and gets from it without statting
status: done
parent: S-0606
owner: alex
created: 2026-11-07T22:48:59Z
updated: 2026-11-08T00:35:11Z
transitions:
  - to: ready
    at: 2026-11-08T00:03:02Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T00:03:10Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T00:35:10Z
    by: agent-S-0606
stream: S-0606
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go, flai/internal/workitem/archive.go, flai/internal/workitem/itemlock.go]
usage:
  source: log
  seconds: 1920
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 12091
      cache_read: 6017965
      cache_write: 251321
      cost: 2.8157
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

- Reviewed by agent-S-0606: keep the snapshot per item folder, not one list per root. A told-of change marks only the folder that holds the path stale, and the next read of a stale folder refreshes it through today's stat-and-compare path, which parses again only the files that changed; so a change under `kanban/stories` never re-stats the archive's 3,800 files. A folder whose refresh read a file within `racyWindow` of its modification time stays stale, so the racy-git rule still holds. Add a call that marks every folder of a root stale (for a write made by a flai subprocess this process ran), and make `Repo.Save` mark its folders stale. Exported: turn snapshot mode on and off for a root, tell the store of a changed path, mark a root stale.
