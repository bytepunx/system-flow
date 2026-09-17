---
id: T-0084
type: task
nature: feature
title: "API routes /api/manifest, /api/items, /api/items/:id, /api/docs/tree, /api/docs/file, /api/events with tests"
status: done
parent: S-0011
owner: alex
created: 2026-09-17T04:02:08Z
updated: 2026-09-17T04:09:26Z
transitions:
  - to: ready
    at: 2026-09-17T04:09:26Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:09:26Z
    by: agent
  - to: done
    at: 2026-09-17T04:09:26Z
    by: agent
stream: S-0011
tags: [dashboard]
---

# T-0084 API routes /api/manifest, /api/items, /api/items/:id, /api/docs/tree, /api/docs/file, /api/events with tests

## Work
+server.ts routes: /api/manifest, /api/items (query by type and status), /api/items/[id] with children, /api/docs/tree, /api/docs/file?path=, /api/events (SSE from the watcher). JSON shapes documented in design/system/flaiover-dashboard.md; path traversal rejected.

## Done when
Vitest tests hit the handlers; curl against pnpm dev returns this repository's data.

## Notes
