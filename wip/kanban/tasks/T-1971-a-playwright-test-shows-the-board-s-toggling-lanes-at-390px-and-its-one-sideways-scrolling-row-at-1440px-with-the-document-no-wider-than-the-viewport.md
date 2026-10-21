---
id: T-1971
type: task
nature: improvement
title: A Playwright test shows the board's toggling lanes at 390px and its one sideways-scrolling row at 1440px, with the document no wider than the viewport
status: done
parent: S-0462
owner: alex
created: 2026-10-21T07:05:07Z
updated: 2026-10-21T09:26:35Z
transitions:
  - to: ready
    at: 2026-10-21T09:14:34Z
    by: agent-S-0462
  - to: in-progress
    at: 2026-10-21T09:14:35Z
    by: agent-S-0462
  - to: done
    at: 2026-10-21T09:26:35Z
    by: agent-S-0462
stream: S-0462
tags: [dashboard]
touches: [flaiover/e2e/board-layout.e2e.ts]
after: [T-1970]
usage:
  source: log
  seconds: 720
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 17768
      cache_read: 2441580
      cache_write: 92930
      cost: 1.4591
---
# T-1971 A Playwright test shows the board's toggling lanes at 390px and its one sideways-scrolling row at 1440px, with the document no wider than the viewport

## Work

- Write `flaiover/e2e/board-layout.e2e.ts` beside `live-events.e2e.ts`, with its setup as it is.
- At 390x844:
  - Ready, In progress, and Review are open, and the others are closed;
  - tapping a heading toggles its lane;
  - the choice survives a reload;
  - the document is no wider than the viewport.
- At 1440x900:
  - the six lanes share one row, with every lane's top at the same y;
  - the row scrolls sideways when it is narrowed;
  - the document is no wider than the viewport.
- Waits for T-1970, the layout it drives.

## Done when

The test passes under the flaiover tier's Playwright run, at both widths.

## Notes
