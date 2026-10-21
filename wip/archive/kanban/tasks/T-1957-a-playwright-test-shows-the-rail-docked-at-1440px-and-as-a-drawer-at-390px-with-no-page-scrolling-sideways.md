---
id: T-1957
type: task
nature: feature
title: A Playwright test shows the rail docked at 1440px and as a drawer at 390px, with no page scrolling sideways
status: done
parent: S-0458
owner: alex
created: 2026-10-21T07:02:24Z
updated: 2026-10-21T08:24:33Z
transitions:
  - to: ready
    at: 2026-10-21T08:14:09Z
    by: agent-S-0458
  - to: in-progress
    at: 2026-10-21T08:14:10Z
    by: agent-S-0458
  - to: done
    at: 2026-10-21T08:24:33Z
    by: agent-S-0458
stream: S-0458
tags: [dashboard]
touches: [flaiover/e2e/layout.e2e.ts, flaiover/src/routes/+layout.svelte]
after: [T-1956]
usage:
  source: log
  seconds: 623
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 12154
      cache_read: 1520241
      cache_write: 52639
      cost: 0.8727
---
# T-1957 A Playwright test shows the rail docked at 1440px and as a drawer at 390px, with no page scrolling sideways

## Work

- Write `flaiover/e2e/layout.e2e.ts` beside `live-events.e2e.ts`, with its setup and `flaiover/playwright.config.ts` as they are.
- At 1440x900, check these on `/board`, an item page, and `/settings`:
  - the rail is visible;
  - expanding and collapsing a group shows and hides its pages;
  - the document is no wider than the viewport.
- At 390x844, check these on the same pages:
  - the rail is hidden until the menu button opens it;
  - choosing a page closes it;
  - the document is no wider than the viewport.
- Waits for T-1956, the layout it drives.

## Done when

The test passes under the flaiover tier's Playwright run, at both widths.

## Notes
