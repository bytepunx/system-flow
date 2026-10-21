---
id: T-1972
type: task
nature: improvement
title: The dashboard's design and user docs describe the board's two layouts and its per-browser open lanes
status: backlog
parent: S-0462
owner: alex
created: 2026-10-21T07:05:13Z
updated: 2026-10-21T07:05:13Z
transitions: []
stream: S-0462
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1970]
---
# T-1972 The dashboard's design and user docs describe the board's two layouts and its per-browser open lanes

## Work

- In `design/system/flaiover-dashboard.md` § Views, add the `/board` row's layouts: the one row from `sm` that scrolls sideways with sticky headings, and the toggling headings below `sm`, with `$lib/lanecollapse.svelte.ts` and its `localStorage` key.
- In `docs/users/flaiover.md` § Board, say how to move between lanes on a wide screen and how to open and close a lane on a phone.
- Waits for T-1970, so that it describes what was built. It runs beside the Playwright test.

## Done when

Both documents describe both layouts, and the markdown tier passes.

## Notes
