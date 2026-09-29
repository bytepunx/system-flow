---
id: T-0583
type: task
nature: feature
title: The charts page offers the new charts with a bucket, a grouping, and a table view each
status: done
parent: S-0163
owner: alex
created: 2026-09-29T20:33:48Z
updated: 2026-09-29T20:57:45Z
transitions:
  - to: ready
    at: 2026-09-29T20:34:13Z
    by: agent-S-0163
  - to: in-progress
    at: 2026-09-29T20:50:30Z
    by: agent-S-0163
  - to: done
    at: 2026-09-29T20:57:45Z
    by: agent-S-0163
stream: S-0163
tags: []
touches: [flaiover/src/routes/charts, docs/users, design/system/flaiover-dashboard.md, design/tech/charts.md]
usage:
  source: log
  seconds: 435
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 34
      output: 283
      cache_read: 5305249
      cache_write: 37492
      cost: 8.7091
---
# T-0583 The charts page offers the new charts with a bucket, a grouping, and a table view each

## Work

On `flaiover/src/routes/charts/[kind]`, list the flow charts and the usage charts as two groups, offer the bucket (hour, day, week) on the charts over time and the grouping (by type, by model) on the charts per item, hide the controls a chart does not use, and give each new chart a table view of what it plots. Say so when the host's flai sends no spend series. Update `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`, and `design/tech/charts.md`.

## Done when

- A component test renders the page's table for a new chart from a report and finds the bucket's numbers.
- The page names the flai it needs when the report has no spend series.
- Prettier, eslint, svelte-check, and vitest pass.

## Notes
