---
id: T-2616
type: task
nature: improvement
title: The watcher stats a slow path's folders on every walk and its files on a slower interval, and a test sees an archived item arrive and a hand edit of one
status: done
parent: S-0606
owner: alex
created: 2026-11-07T22:49:05Z
updated: 2026-11-08T00:36:22Z
transitions:
  - to: ready
    at: 2026-11-08T00:03:17Z
    by: agent-S-0606
  - to: in-progress
    at: 2026-11-08T00:03:22Z
    by: agent-S-0606
  - to: done
    at: 2026-11-08T00:36:22Z
    by: agent-S-0606
stream: S-0606
tags: [cli]
touches: [flai/internal/watch/watch.go, flai/internal/watch/watch_test.go]
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
# T-2616 The watcher stats a slow path's folders on every walk and its files on a slower interval, and a test sees an archived item arrive and a hand edit of one

## Work

Give `watch.Watcher` in `flai/internal/watch/watch.go` paths on a slow tier, with an interval of their own that defaults to 10 s. A walk of a slow path stats its folders on every tick, so a file added, renamed in, or removed shows as the folder's modification time moving, and the folder is then read again. It stats every file below it only on the slow interval, since an edit in place does not move a folder's modification time. The fast paths keep today's walk every `Every`. Debounce still reports a file once it has looked the same for one tick, keeping the slow tier's file stats as their own snapshot so that a skipped file is not reported as removed. Polling stays: `design/tech/go-libraries.md` records why fsnotify was not taken.

Waits for no task. It shares no path with the store task, so the two run together.

## Done when

- [ ] Tests in `flai/internal/watch/watch_test.go` show that under a slow path a file moved in is reported within a fast tick or two, a hand edit in place is reported within the slow interval, a removal is reported, and a walk between slow ticks stats none of its files.
- [ ] The existing watch tests pass unchanged.

## Notes
