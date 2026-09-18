---
id: T-0159
type: task
nature: improvement
title: Board API returns parent_title on every card
status: done
parent: S-0048
owner: alex
created: 2026-09-18T20:41:13Z
updated: 2026-09-18T20:42:02Z
transitions:
  - to: ready
    at: 2026-09-18T20:41:14Z
    by: alex
  - to: in-progress
    at: 2026-09-18T20:41:14Z
    by: alex
  - to: done
    at: 2026-09-18T20:42:02Z
    by: alex
stream: S-0048
tags: []
touches: [flaiover/src/lib/server, design/system]
---

# T-0159 Board API returns parent_title on every card

## Work
In `flaiover/src/lib/server/board.ts`, add optional `parent_title` to `Card`, filled from the item whose ID is the card's `parent`, looked up across every item the reader returns, archived ones included, so a parent is found whatever its state. Leave it out when there is no parent or the parent does not exist. Add it to the card contract in `design/system/flaiover-dashboard.md` (`GET /api/board`). Extend the board test in `src/lib/server/flai.test.ts` (or a board test beside it) against the `good` fixture: a story card carries its epic's title, a task card its story's title, an epic card none.

## Done when
- The new assertions fail without the change and pass with it
- The contract line in `flaiover-dashboard.md` lists `parent_title`

## Notes
