---
id: T-1432
type: task
nature: feature
title: The cancelled lane's menu offers Archive All and a cancelled card's menu offers Archive
status: done
parent: S-0343
owner: alex
created: 2026-10-09T16:38:57Z
updated: 2026-10-09T16:45:53Z
transitions:
  - to: ready
    at: 2026-10-09T16:39:21Z
    by: agent-S-0343
  - to: in-progress
    at: 2026-10-09T16:39:22Z
    by: agent-S-0343
  - to: done
    at: 2026-10-09T16:45:53Z
    by: agent-S-0343
stream: S-0343
tags: []
touches: [flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/cardmenu.ts, flaiover/src/lib/cardmenu.test.ts, flaiover/src/lib/components/LaneMenu.svelte, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
usage:
  source: log
  seconds: 391
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 68
      output: 19255
      cache_read: 3743896
      cache_write: 120349
      cost: 1.9481
---

# T-1432 The cancelled lane's menu offers Archive All and a cancelled card's menu offers Archive

## Work

- `lanes.ts`: a lane action `archive`, labelled **Archive All**, offered by the cancelled lane only while it holds cards; `laneEntries` learns how many cards the lane holds. A pure function orders the lane's IDs for archiving one by one: stories first (each brings its tasks), then tasks, then epics (flai refuses an epic while a story of it is on the board).
- `cardmenu.ts`: a card action `archive`, labelled **Archive**, on a cancelled card that is not archived, for a writer.
- `LaneMenu.svelte`, `CardMenu.svelte`, `board/+page.svelte`: pass the lane's count, and on the choice `POST /api/board/archive` with `{ ids: [id] }`, one item per call, skipping an ID an earlier call archived with its story; the notice names what was archived and what flai refused and why, as **Move stories** does, and the board reloads.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`: only the card-menu and lane-menu sections and the board API row for the route (the split on MS-0046 leaves the charts sections to S-0337).

Waits for nothing: its component tests mock the route, so it runs beside T-1430.

## Done when

- `lanes.test.ts`, `cardmenu.test.ts`, `lanemenu.svelte.test.ts`, and `cardmenu.svelte.test.ts` cover both entries: offered where the criteria say, absent elsewhere, and each calling the route with the right IDs.
- `flai test` on the changed paths passes.

## Notes
