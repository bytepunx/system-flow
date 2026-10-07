---
id: T-0943
type: task
nature: feature
title: The Charts page lists parallelism, hold time, and touches drift under Planning
status: done
parent: S-0214
owner: alex
created: 2026-10-05T05:45:39Z
updated: 2026-10-07T08:41:13Z
transitions:
  - to: ready
    at: 2026-10-07T08:33:40Z
    by: agent-S-0214
  - to: in-progress
    at: 2026-10-07T08:33:40Z
    by: agent-S-0214
  - to: done
    at: 2026-10-07T08:41:13Z
    by: agent-S-0214
stream: S-0214
tags: [dashboard]
touches: [flaiover/src/routes/charts]
after: [T-0929, T-0936]
usage:
  source: log
  seconds: 453
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 80
      output: 458
      cache_read: 3628603
      cache_write: 115800
      cost: 1.682
---
# T-0943 The Charts page lists parallelism, hold time, and touches drift under Planning

## Work

In `flaiover/src/routes/charts/[kind]/+page.svelte`, add a `planning` group from `PLANNING_KINDS` beside `flow` and `usage`, or add to it if S-0212 or S-0213 made it first. Render `/charts/parallelism`, `/charts/hold-time`, and `/charts/touches-drift` with T-0936's builders, from `/api/stats` over the window chosen. Then run the dashboard against this repository with a flai built from the story branch, and compare each chart's points with `flai stats --json` for the same window. Waits for T-0936, whose builders it renders, and for T-0929, whose values the live comparison needs.

## Done when

- [ ] The three charts appear under Planning, and each reads `/api/stats` over the window
- [ ] `charts.svelte.test.ts` covers the group and each chart's series
- [ ] The charts' points match `flai stats --json` for the same window, checked by hand and logged in the narrative
- [ ] `scripts/flaiover-test.sh` passes

## Notes
