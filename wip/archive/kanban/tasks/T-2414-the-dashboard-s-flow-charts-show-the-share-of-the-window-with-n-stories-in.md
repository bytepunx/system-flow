---
id: T-2414
type: task
nature: improvement
title: The dashboard's flow charts show the share of the window with n stories in progress and the ready wait at p50 and p85
status: done
parent: S-0551
owner: alex
created: 2026-10-31T07:44:48Z
updated: 2026-10-31T10:17:51Z
transitions:
  - to: ready
    at: 2026-10-31T10:07:16Z
    by: agent-S-0551
  - to: in-progress
    at: 2026-10-31T10:07:18Z
    by: agent-S-0551
  - to: done
    at: 2026-10-31T10:17:51Z
    by: agent-S-0551
stream: S-0551
tags: [flaiover, dashboard, metrics]
touches: [flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/sitemenu.ts, "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-2410]
usage:
  source: log
  seconds: 633
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 120
      output: 37923
      cache_read: 7457257
      cache_write: 273502
      cost: 3.7744
---
# T-2414 The dashboard's flow charts show the share of the window with n stories in progress and the ready wait at p50 and p85

## Work

Add a flow chart that draws the share of the stats window with n stories in progress at once, from the key `flai stats --json` gives, and show the ready wait at p50 and p85 with it. Add the kind to `FLOW_KINDS` in `catalog.ts` and to the flow group in `sitemenu.ts`, its builder in `charts.ts` beside `parallelism`, and its explanation in `explain.ts`. Document the chart in `docs/users/flaiover.md`'s chart table and in `design/system/flaiover-dashboard.md`.

It waits for the stats task, which gives the data.

## Done when

- Tests in `charts.test.ts` and `catalog.test.ts` show the chart built from a report and listed in the flow group, and `charts.svelte.test.ts` shows its page rendering.
- `flai test` passes on the changed flaiover files.

## Notes

Parallelism and Hold Time sit in the planning group today; leave them there.
