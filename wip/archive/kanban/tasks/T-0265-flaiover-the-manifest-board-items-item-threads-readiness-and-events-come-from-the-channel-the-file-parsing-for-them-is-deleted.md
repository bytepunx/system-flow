---
id: T-0265
type: task
nature: feature
title: "flaiover: the manifest, board, items, item, threads, readiness, and events come from the channel; the file parsing for them is deleted"
status: done
parent: S-0073
owner: alex
created: 2026-09-20T07:52:15Z
updated: 2026-09-20T08:05:46Z
transitions:
  - to: ready
    at: 2026-09-20T07:57:29Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:57:30Z
    by: system-flow
  - to: done
    at: 2026-09-20T08:05:46Z
    by: system-flow
stream: S-0073
tags: []
---
# T-0265 flaiover: the manifest, board, items, item, threads, readiness, and events come from the channel; the file parsing for them is deleted

## Work
repo.ts asks the channel for the manifest, items, and threads and keeps the shapes its callers use; board.ts maps board.get; the hub turns change notifications into the events /api/events, the caches, and the notifier already listen to; /_ready asks the channel; with no flai connected those routes answer 503 with what is missing and the command that starts it, and the pages show it.

## Done when
- Unit tests with a fake hub for each route and for the 503
- The item and thread parsers and board composition are gone from flaiover

## Notes
