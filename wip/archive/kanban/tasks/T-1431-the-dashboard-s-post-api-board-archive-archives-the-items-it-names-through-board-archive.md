---
id: T-1431
type: task
nature: feature
title: The dashboard's POST /api/board/archive archives the items it names through board.archive
status: done
parent: S-0343
owner: alex
created: 2026-10-09T16:38:55Z
updated: 2026-10-09T16:46:48Z
transitions:
  - to: ready
    at: 2026-10-09T16:39:22Z
    by: agent-S-0343
  - to: in-progress
    at: 2026-10-09T16:45:58Z
    by: agent-S-0343
  - to: done
    at: 2026-10-09T16:46:48Z
    by: agent-S-0343
stream: S-0343
tags: []
touches: [flaiover/src/routes/api/board/archive/+server.ts, flaiover/src/lib/server/writes.test.ts]
after: [T-1430]
usage:
  source: log
  seconds: 50
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 32
      output: 9032
      cache_read: 1756239
      cache_write: 56455
      cost: 0.9139
---

# T-1431 The dashboard's POST /api/board/archive archives the items it names through board.archive

## Work

Add `flaiover/src/routes/api/board/archive/+server.ts`, `POST { ids }`, passing the IDs to flai's `board.archive` as `api/board/limit` passes its body to `board.limit`, and add `board.archive` to the write methods in `flaiover/src/lib/server/agent.ts`. Waits for T-1430: `writes.test.ts` runs the methods through `flai hostapi` built from this tree, so it needs `board.archive` there.

## Done when

- `writes.test.ts` archives a cancelled item of the fixture copy through `board.archive` and sees it leave the board, and sees an open item and an empty list refused.
- `flai test` on the changed paths passes.

## Notes
