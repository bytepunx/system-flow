---
id: T-0572
type: task
nature: feature
title: Pages gather changes that arrive together, and the board asks again only when a work item changed
status: in-progress
parent: S-0161
owner: alex
created: 2026-09-29T19:40:45Z
updated: 2026-09-29T19:45:14Z
transitions:
  - to: ready
    at: 2026-09-29T19:41:11Z
    by: agent-S-0161
  - to: in-progress
    at: 2026-09-29T19:45:14Z
    by: agent-S-0161
stream: S-0161
tags: []
touches: [flaiover/src, design/system]
---

# T-0572 Pages gather changes that arrive together, and the board asks again only when a work item changed

## Work

- `debounced` in `src/lib/events.ts` gathers the changes that arrive within 500 ms of each other and calls once with all of them.
- The board follows the shared event stream: its board, `/api/publish`, and the unpushed notice are asked again only for a gathered change of kind `item`, `project`, or `other`; the agents on `agent` events and item changes.
- The activity, charts, and ADR pages, the item page, threads, story agent, and inbox gather changes the same way and ask again only for the kinds their answers are read from.
- `design/system/server-performance.md` records cause 6 as addressed; `design/system/flaiover-dashboard.md` says how pages follow changes.

## Done when

- Tests show changes within half a second give one call, and which kinds make the board ask again.
- `npm run check`, lint, and the unit tests pass in `flaiover/`.

## Notes
