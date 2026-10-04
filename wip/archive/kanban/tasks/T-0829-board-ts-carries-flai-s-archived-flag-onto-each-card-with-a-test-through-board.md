---
id: T-0829
type: task
nature: remediation
title: board.ts carries flai's archived flag onto each card, with a test through board()
status: done
parent: S-0247
owner: alex
created: 2026-10-04T21:43:06Z
updated: 2026-10-04T21:49:34Z
transitions:
  - to: ready
    at: 2026-10-04T21:43:14Z
    by: agent-S-0247
  - to: in-progress
    at: 2026-10-04T21:43:14Z
    by: agent-S-0247
  - to: done
    at: 2026-10-04T21:49:34Z
    by: agent-S-0247
stream: S-0247
tags: []
touches: [flaiover/src/lib/server/board.ts, flaiover/src/lib/server/repo-channel.test.ts]
usage:
  source: log
  seconds: 380
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 6169
      cache_read: 1293338
      cache_write: 40667
      cost: 0.6761
---

# T-0829 board.ts carries flai's archived flag onto each card, with a test through board()

## Work

Add `archived?: boolean` to `Card` in `flaiover/src/lib/server/board.ts` and carry it in `board()`'s mapping the way `draft` is carried: present and `true` when flai's `board.get` card has `archived: true`, absent otherwise. Extend the board test in `repo-channel.test.ts`, which goes through `board()` rather than a mocked `/api/board`, with a done, archived card, and check that `doneLane()` leaves it out while the clone lags the remote's tags. The only task of S-0247; it waits for nothing.

## Done when

- A card flai marks archived reaches the page with `archived: true`; one it does not mark has no `archived` property.
- The test fails without the change to `board.ts` and passes with it.

## Notes
