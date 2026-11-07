---
id: T-2618
type: task
nature: improvement
title: flai serve turns on each served project's item snapshot, replaces it from the watcher's changes under the item folders, and watches wip/archive on the slow tier
status: backlog
parent: S-0606
owner: alex
created: 2026-11-07T22:49:22Z
updated: 2026-11-07T22:49:22Z
transitions: []
stream: S-0606
tags: [cli]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go]
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
