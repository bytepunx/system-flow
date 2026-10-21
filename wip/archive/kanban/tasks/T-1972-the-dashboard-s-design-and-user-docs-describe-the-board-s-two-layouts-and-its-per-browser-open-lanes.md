---
id: T-1972
type: task
nature: improvement
title: The dashboard's design and user docs describe the board's two layouts and its per-browser open lanes
status: done
parent: S-0462
owner: alex
created: 2026-10-21T07:05:13Z
updated: 2026-10-21T09:15:04Z
transitions:
  - to: ready
    at: 2026-10-21T09:14:36Z
    by: agent-S-0462
  - to: in-progress
    at: 2026-10-21T09:14:36Z
    by: agent-S-0462
  - to: done
    at: 2026-10-21T09:15:04Z
    by: agent-S-0462
stream: S-0462
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1970]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 7
      output: 2688
      cache_read: 369432
      cache_write: 14061
      cost: 0.2208
---
# T-1972 The dashboard's design and user docs describe the board's two layouts and its per-browser open lanes

## Work

- In `design/system/flaiover-dashboard.md` § Views, add the `/board` row's layouts: the one row from `sm` that scrolls sideways with sticky headings, and the toggling headings below `sm`, with `$lib/lanecollapse.svelte.ts` and its `localStorage` key.
- In `docs/users/flaiover.md` § Board, say how to move between lanes on a wide screen and how to open and close a lane on a phone.
- Waits for T-1970, so that it describes what was built. It runs beside the Playwright test.

## Done when

Both documents describe both layouts, and the markdown tier passes.

## Notes
