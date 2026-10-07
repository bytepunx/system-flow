---
id: T-0936
type: task
nature: feature
title: The dashboard's Report carries claims, and builders map it to the parallelism, hold-time, and touches-drift charts
status: in-progress
parent: S-0214
owner: alex
created: 2026-10-05T05:45:20Z
updated: 2026-10-07T08:23:20Z
transitions:
  - to: ready
    at: 2026-10-07T08:23:20Z
    by: agent-S-0214
  - to: in-progress
    at: 2026-10-07T08:23:20Z
    by: agent-S-0214
stream: S-0214
tags: [dashboard]
touches: [flaiover/src/lib/viz]
after: [T-0924]
---
# T-0936 The dashboard's Report carries claims, and builders map it to the parallelism, hold-time, and touches-drift charts

## Work

In `flaiover/src/lib/viz/charts.ts`, add `claims` to `Report` with the shape `metrics.md` § Claims and touches gives: `limit`, `days[]` with `in_progress` and `held`, `weeks[]` with held seconds per reason, `drift[]`, and the weekly exact-touches share. Add three builders that return ECharts options, as the existing ones do:

- `parallelism`: stories in progress per day as a line, the in-progress limit as a flat line, and held stories as a second series.
- `holdTime`: hours held per week as bars stacked by reason (overlap, after, empty claim).
- `touchesDrift`: per story completed, files changed outside its touches and touches never changed as stacked bars, with the weekly share of exact touches on a second axis.

Each spans the window as ADR-0054 requires. Add the three kinds to `TITLES` and a `PLANNING_KINDS` list; reuse it if S-0212 or S-0213 made it first. Waits for T-0924, whose contract the types follow. It can run beside T-0929, since their paths are apart.

## Done when

- [ ] `Report` types `claims` as `flai stats --json` emits it
- [ ] `charts.test.ts` covers each builder's data mapping on a fixture shaped as `metrics.md` defines, including an empty window, no limit, and absent `drift`
- [ ] `scripts/flaiover-test.sh` passes

## Notes
