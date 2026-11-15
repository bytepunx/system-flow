---
id: T-3383
type: task
nature: improvement
title: Close I-0375 with what fixed it
status: ready
parent: S-0790
owner: alex
created: 2026-11-15T21:11:57Z
updated: 2026-11-15T21:18:47Z
transitions:
  - to: ready
    at: 2026-11-15T21:18:47Z
    by: system-flow
stream: S-0790
tags: [issues]
touches: [design/issues/I-0375-flai-check-finds-board-order-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-3382]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 0
          output: 6
          cache_read: 20828
          cache_write: 530
          cost: 0.0065
---
# T-3383 Close I-0375 with what fixed it

## Work

Close the issue with `flai issue close I-0375 --reason "<reason>"`. The reason names the cause: concurrent story starts saved `board.md` without a lock and wrote back a story another start had removed from the order. It names the fix: every write of `board.md` now takes the lock T-3380's ADR names, reloads the board inside it, and drops names that are no longer ready or backlog stories. It names the ADR and the test T-3381 added. `flai issue close` regenerates `design/issues/summary.md`.

It waits for T-3382 so that the issue closes only once every writer of `board.md` takes the lock.

## Done when

- I-0375 is closed, and its reason names the cause, the fix, the ADR, and the test.
- `design/issues/summary.md` no longer lists I-0375 as open.
- `flai check --strict` passes on the issue files.

## Notes
