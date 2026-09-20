---
id: T-0264
type: task
nature: feature
title: "flai: a watcher that tells the dashboard which files changed"
status: done
parent: S-0073
owner: alex
created: 2026-09-20T07:52:15Z
updated: 2026-09-20T07:57:29Z
transitions:
  - to: ready
    at: 2026-09-20T07:55:54Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:55:54Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:57:29Z
    by: system-flow
stream: S-0073
tags: []
---
# T-0264 flai: a watcher that tells the dashboard which files changed

## Work
A watcher in flai serve on each project's design, docs, and wip folders and its manifest, debounced as chokidar is today, sending a change notification with the repository-relative path over the project's connection. Recorded in design/tech.

## Done when
- Tests for a change, a new file in a new folder, a removal, and the debounce
- Nothing is sent for paths outside the manifest's folders

## Notes
