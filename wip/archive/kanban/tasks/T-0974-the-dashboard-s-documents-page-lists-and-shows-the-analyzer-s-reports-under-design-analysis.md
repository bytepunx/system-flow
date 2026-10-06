---
id: T-0974
type: task
nature: feature
title: The dashboard's documents page lists and shows the analyzer's reports under design/analysis
status: done
parent: S-0223
owner: alex
created: 2026-10-05T05:47:53Z
updated: 2026-10-06T20:35:38Z
transitions:
  - to: ready
    at: 2026-10-06T20:28:08Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:28:08Z
    by: agent-S-0223
  - to: review
    at: 2026-10-06T20:35:38Z
    by: agent-S-0223
  - to: done
    at: 2026-10-06T20:35:38Z
    by: agent-S-0223
stream: S-0223
tags: [dashboard]
touches: ["flaiover/src/routes/docs/[...path]/+page.svelte", flaiover/src/routes/docs/docs.svelte.test.ts, flai/internal/hostapi/docs_test.go, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0960]
usage:
  source: log
  seconds: 450
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 17
      output: 6853
      cache_read: 1095361
      cache_write: 26720
      cost: 0.5038
---
# T-0974 The dashboard's documents page lists and shows the analyzer's reports under design/analysis

## Work

The documents page (`flaiover/src/routes/docs/[...path]/+page.svelte`) reads its tree from the host's design, docs, and wip folders (`flai/internal/hostapi/docs.go`, `layoutDirs`), so `design/analysis/` should appear with no change of its own. Prove it and close what is missing:

- a host test that the docs tree lists `design/analysis/README.md` and a report, and serves the report (`flai/internal/hostapi/docs_test.go`)
- a page test that the tree shows `analysis` under design and opens a report with its front matter, `focus` and the window among it (`flaiover/src/routes/docs/docs.svelte.test.ts`)
- if the page hides or orders the folder differently from the others, make it list the reports as it lists `experiments`

Say in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` that the documents page shows the analyzer's reports.

It waits for T-0960, which adds the folder and its README.

## Done when

- both tests pass, `go test ./internal/hostapi/` and the dashboard's tests and lint (`npm run check`, `npm test` in `flaiover/`)
- the dashboard design and the user guide name the reports on the documents page

## Notes
