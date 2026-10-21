---
id: T-2010
type: task
nature: feature
title: An end-to-end test finds that the review, charts, and documents pages never scroll sideways below md and keep the sidebar beside the content from md
status: done
parent: S-0461
owner: alex
created: 2026-10-21T12:18:59Z
updated: 2026-10-21T12:53:10Z
transitions:
  - to: ready
    at: 2026-10-21T12:19:35Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T12:19:36Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T12:53:10Z
    by: agent-S-0461
stream: S-0461
tags: []
touches: [flaiover/e2e/page-sidebar.e2e.ts]
after: [T-1965, T-1966, T-1967, T-2011]
usage:
  source: log
  seconds: 2014
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 81
      output: 37713
      cache_read: 4817683
      cache_write: 167081
      cost: 2.7075
---
# T-2010 An end-to-end test finds that the review, charts, and documents pages never scroll sideways below md and keep the sidebar beside the content from md

## Work

- Add `flaiover/e2e/page-sidebar.e2e.ts`, after `board-layout.e2e.ts`'s pattern: sign in with the run's token, answer the pages' own API requests in the browser, and open `/review/<id>`, `/charts/<kind>`, and `/docs/<path>` at a phone's width and a desktop's.
- Below `md`, each page reaches no further than the window sideways, and its sidebar is one row above the content. From `md`, the sidebar sits at the left of the content.
- Waits for T-1965, T-1966, and T-1967, done, whose pages it opens. It runs beside T-1968, which touches only documents.

## Done when

`scripts/flaiover-e2e.sh e2e/page-sidebar.e2e.ts` passes, which verifies criterion 4's narrow width in a real browser.

## Notes
