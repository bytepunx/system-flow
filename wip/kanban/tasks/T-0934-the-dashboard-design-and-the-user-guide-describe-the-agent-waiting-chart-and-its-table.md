---
id: T-0934
type: task
nature: feature
title: The dashboard design and the user guide describe the agent-waiting chart and its table
status: done
parent: S-0215
owner: alex
created: 2026-10-05T05:45:18Z
updated: 2026-10-07T09:05:36Z
transitions:
  - to: ready
    at: 2026-10-07T08:54:33Z
    by: agent-S-0215
  - to: in-progress
    at: 2026-10-07T09:02:04Z
    by: agent-S-0215
  - to: done
    at: 2026-10-07T09:05:36Z
    by: agent-S-0215
stream: S-0215
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0927]
usage:
  source: log
  seconds: 212
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 7284
      cache_read: 973817
      cache_write: 34087
      cost: 0.5464
---
# T-0934 The dashboard design and the user guide describe the agent-waiting chart and its table

## Work

Describe `/charts/agent-waiting` in `design/system/flaiover-dashboard.md` and in `docs/users/flaiover.md`, where the other Flow charts are described: what a bar and the line show, where the figures come from (`waiting` in `/api/stats`, defined in `metrics.md`), and what the table under the chart lists.

It waits for T-0927, which fixes what the chart plots. It touches no path that T-0932's table does, so it runs beside it.

## Done when

- Both documents describe the chart and its table, and link `metrics.md` for the definitions.
- The markdown lint and `flai check --strict` are clean.

## Notes
