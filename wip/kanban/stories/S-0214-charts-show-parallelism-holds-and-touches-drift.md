---
id: S-0214
type: story
nature: feature
title: Charts show parallelism, holds, and touches drift
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-06T23:17:59Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:54Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:39Z
    by: alex
tags: [dashboard]
topics: [planning, analysis]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts, flai/internal/metrics, design/system/metrics.md, design/adrs, docs/users/flai.md, flai/cmd/stats.go, flai/cmd/check_stats_test.go, docs/users/flai-reference.md]
after: [S-0205]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 106.42
  by: planner-S-0214
  at: 2026-10-05T05:46:07Z
forecast:
  duration: 1h45m
  delivery: 2026-10-07T04:19:00Z
  basis: "Its own forecast of 1h45m; 8th in the pull order with an in-progress limit of 3, behind S-0261, S-0300, S-0301, S-0302, S-0228, S-0269, S-0270, S-0271, S-0212 and S-0213."
  by: flai
  at: 2026-10-06T23:17:59Z
---
# S-0214 Charts show parallelism, holds, and touches drift

## Goal

The planner's touches exist to raise parallelism safely. The operator should see how many stories run at once, how long claims hold work, and how far declared touches are from what was changed.

## Acceptance criteria
- [ ] `/charts/parallelism`: stories in progress per day against the in-progress limit, with held stories as a second series
- [ ] `/charts/hold-time`: per week the hours stories spent held, by reason (overlap, after, empty claim)
- [ ] `/charts/touches-drift`: per story completed, files changed outside its touches and touches never changed, as stacked bars, with the share of stories whose touches were exact per week
- [ ] Charts span the window, read `/api/stats`, match `flai stats --json`; listed under Planning; design and user guide describe them; tests cover the data mapping

## Tasks
- T-0924 metrics.md defines held stories per day, held hours by reason per week, and the weekly share of exact touches, and an ADR records them
- T-0929 flai stats reports held stories per day, held hours by reason per week, and the weekly share of exact touches
- T-0933 flai stats prints held time by reason and the exact-touches share, and the flai guide describes them
- T-0936 The dashboard's Report carries claims, and builders map it to the parallelism, hold-time, and touches-drift charts
- T-0943 The Charts page lists parallelism, hold time, and touches drift under Planning
- T-0946 The dashboard design and user guide describe the parallelism, hold-time, and touches-drift charts

## Notes

### Planning

By planner-S-0214, 2026-10-05.

S-0205's `claims` gives `in_progress` per day, `held_seconds` per story, and `drift[]`. It has no held stories per day, no held time by reason or by week, and no weekly share of exact touches. ADR-0081 says the charts need no figures of their own, and `metrics.md` changes only with an ADR. So the plan adds those figures to `flai stats` first (T-0924, T-0929, T-0933), then builds the charts on them (T-0936, T-0943, T-0946).

Layers: T-0924; then T-0929 and T-0936 together, their paths apart; then T-0933 and T-0943; then T-0946.

Touches, each kept from what the story declared and the rest added:

- Declared: `flaiover/src/routes/charts`, `flaiover/src/lib/charts`, `flaiover/src/lib/viz`, `flaiover/src/lib/sitemenu.ts`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `flai/internal/metrics`, `design/system/metrics.md`, `design/adrs`, `docs/users/flai.md`. The layout shows `flaiover/src/lib/charts` does not exist: the chart model is `flaiover/src/lib/viz/charts.ts`. The Planning group is on the Charts page, not in `sitemenu.ts`. No task changes either path, and both are kept as declared.
- Layout: `flai/cmd/stats.go` and `flai/cmd/check_stats_test.go`, where `printClaims` prints the claims aggregates; `docs/users/flai-reference.md`, which mirrors the `flai stats` help.
- Co-change: `flai touches suggest` puts `design/system/flai-cli.md` first (51%), but its `flai stats` row adds no flags here, so it is left out.

Topics: `planning` (the charts judge the planner's touches) and `analysis` (the `metrics.md` contract), added because T-0924 and T-0929 reach them.

Forecast: 1h45m, kept over flai's 31m. flai's figure is 132 s a unit of size over 14 (4 criteria, 10 touches), which treats this as a chart story alone. The plan adds the Go metrics, an ADR, and the `flai stats` output, and S-0205's comparable metrics work took 70 minutes in progress. The charts and docs add about 35 minutes. Delivery 2026-10-05T15:47Z: flai's start (14:02Z, 11th in the pull order with an in-progress limit of 3) plus 1h45m.

Cost of delay: 106.42 USD a week, from `flai cod`. The story has no inputs of its own, so it takes its forecast's share (1h45m of 24h40m) of E-0016's 1500 USD a week over its open stories without inputs. This replaces E-0016's planner's 85.37, which came from an earlier share of the same epic value.
