---
id: T-2618
type: task
nature: improvement
title: flai serve turns on each served project's item snapshot, replaces it from the watcher's changes under the item folders, and watches wip/archive on the slow tier
status: done
parent: S-0606
owner: alex
created: 2026-11-07T22:49:22Z
updated: 2026-11-08T01:11:53Z
transitions:
  - to: ready
    at: 2026-11-08T00:47:18Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T00:47:23Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T01:11:53Z
    by: agent-S-0606
stream: S-0606
tags: [cli]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/formula_schedule.go]
after: [T-2615, T-2616]
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
          cache_read: 62744
          cache_write: 12650
          cost: 0.0192
---
# T-2618 flai serve turns on each served project's item snapshot, replaces it from the watcher's changes under the item folders, and watches wip/archive on the slow tier

## Work

In `flai/internal/serve/serve.go`, when a project is served, turn on the store's snapshot mode for its root, and turn it off when the project stops being served. In the `watcher.Run` callback, tell the store of each changed path under the manifest's kanban and archive folders. In `watchedPaths`, put the wip folder's `archive` on the watcher's slow tier and leave the rest on the fast one. The merge queue's own watcher on `.git` is unchanged.

Waits for T-2615, whose snapshot mode it turns on, and T-2616, whose slow tier it uses.

## Done when

- [ ] A test in `flai/internal/serve/serve_test.go` shows that a served project's item list reflects an item written by another process once the watcher reports it, and an archived item edited by hand once the slow tier reports it.
- [ ] `watchedPaths` puts the archive on the slow tier and the other folders on the fast one.

## Notes

- Reviewed by agent-S-0606: touches widened to `flai/internal/serve/mergequeue.go` and `flai/internal/serve/formula_schedule.go`. In the watcher callback, tell the store of the path before the dashboard is notified and before the looks run, so that what either reads next is fresh. The merge queue's and the formula scheduler's flai subprocesses write items in the main checkout: mark the root stale when each returns, with T-2615's call, as T-2617 does for the host API's. Turn snapshot mode off when a project stops being served or drains.
