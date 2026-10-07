---
id: T-0947
type: task
nature: feature
title: The dashboard design and the user guide describe the Strategic Cost and Strategic Use charts
status: done
parent: S-0216
owner: alex
created: 2026-10-05T05:45:49Z
updated: 2026-10-07T09:43:31Z
transitions:
  - to: ready
    at: 2026-10-07T09:32:19Z
    by: agent-S-0216
  - to: in-progress
    at: 2026-10-07T09:32:20Z
    by: agent-S-0216
  - to: done
    at: 2026-10-07T09:43:31Z
    by: agent-S-0216
stream: S-0216
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-0939]
usage:
  source: log
  seconds: 671
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 13619
      cache_read: 1597119
      cache_write: 104330
      cost: 1.1841
---
# T-0947 The dashboard design and the user guide describe the Strategic Cost and Strategic Use charts

## Work

- Waits for T-0939, which fixes the kinds, the group, and the axis. It shares no path with T-0941 or the page task, so it runs beside them.
- `design/system/flaiover-dashboard.md` § Views: a row each for `/charts/strategic-cost` and `/charts/strategic-use`, read from `strategic_days[]` and `items[]` of `/api/stats`, under the Strategic group.
- `docs/users/flaiover.md` § Charts: what each chart shows, how to read the ratio and the lines, and that the figures match `flai stats --json`.
- Write each in the style `documentation.md` sets.

## Done when

- Both documents name both routes, their source fields, and how to read them, and agree with T-0939 and T-0941.
- The markdown lint passes on both files.

## Notes
