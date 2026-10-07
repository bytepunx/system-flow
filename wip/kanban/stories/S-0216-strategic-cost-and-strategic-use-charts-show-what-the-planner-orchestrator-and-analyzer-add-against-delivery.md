---
id: S-0216
type: story
nature: feature
title: Strategic Cost and Strategic Use charts show what the planner, orchestrator, and analyzer add against delivery
status: in-progress
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:15Z
updated: 2026-10-07T09:18:27Z
transitions:
  - to: ready
    at: 2026-10-03T20:33:57Z
    by: alex
  - to: backlog
    at: 2026-10-04T00:41:16Z
    by: alex
  - to: ready
    at: 2026-10-07T03:29:17Z
    by: alex
  - to: in-progress
    at: 2026-10-07T09:18:27Z
    by: agent-S-0216
tags: [dashboard]
touches: [flaiover/src/routes/charts, flaiover/src/lib/charts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/viz, flaiover/src/lib/sitemenu.ts]
after: [S-0205, S-0225, S-0226, S-0227]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 10
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 28
          cache_read: 386507
          cache_write: 5926
          cost: 0.1024
cost_of_delay:
  value: 46.84
  by: planner-S-0216
  at: 2026-10-05T05:50:17Z
forecast:
  duration: 40m
  delivery: 2026-10-07T12:42:00Z
  basis: "Its own forecast of 40m; 1st in the pull order with an in-progress limit of 3, behind S-0215, S-0246 and S-0275."
  by: flai
  at: 2026-10-07T08:59:42Z
---
# S-0216 Strategic Cost and Strategic Use charts show what the planner, orchestrator, and analyzer add against delivery

## Goal

The operator should be able to tell whether the strategic agents pay for themselves. Two charts contrast what they cost and the time they take with the average cost and delivery time per story over the same period.

## Acceptance criteria
- [ ] `/charts/strategic-cost`: per bucket, the cost of the planner, orchestrator, and analyzer (stacked, from their activity logs) as bars, with the mean cost per story completed in the bucket as a line on the same axis, and the ratio stated
- [ ] `/charts/strategic-use`: per bucket, the agent seconds of the three as bars, with the mean cycle time and the mean agent waiting time per story completed as lines, so a fall in waiting or cycle time can be read against the agents' time
- [ ] Both span the window, read `/api/stats`, match `flai stats --json`; listed under a Strategic group; a note explains how to read them; design and user guide describe them; tests cover the mapping

## Tasks
- T-0939 The chart model reads strategic_days and builds the Strategic Cost chart in a Strategic group
- T-0941 The chart model builds the Strategic Use chart: the agents' hours against mean cycle time and waiting per story
- T-0945 The charts page lists a Strategic group and explains how to read the two charts
- T-0947 The dashboard design and the user guide describe the Strategic Cost and Strategic Use charts

## Notes

### Planning

- The first criterion said "on a second axis". TH-0140 settled one shared y-axis, as `design/tech/charts.md` rules and `charts.test.ts` enforces. The bars and the line share a unit, dollars here and hours on Strategic Use, so the criterion now says "on the same axis".
- Assumptions in TH-0140, which the operator accepted:
  - Buckets are the days of the window, as `strategic_days[]` gives them.
  - The mean waiting per story completed is computed in the dashboard from `items[].wait_threads_seconds` and `wait_review_seconds`, so the metrics contract and flai are unchanged.
  - The ratio is the strategic cost over the window, per story completed, as a share of the mean agent cost per story.
- Touches, all declared and kept: `flaiover/src/routes/charts`, `flaiover/src/lib/viz`, `flaiover/src/lib/sitemenu.ts`, `flaiover/src/lib/charts`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.
  - Layout: the work falls in `flaiover/src/lib/viz/charts.ts` and its test, and in `flaiover/src/routes/charts/[kind]/+page.svelte` and `charts.svelte.test.ts`. The declared folders cover all four, so nothing is added.
  - `flaiover/src/lib/charts` does not exist, and `sitemenu.ts` holds only the Charts menu entry, since the chart groups live in the page. Both are kept as declared, though no task changes them.
- Touches left out:
  - Co-change: `design/system/flai-cli.md` (44%), `docs/users/flai.md` (40%), and `docs/operators/index.md` (34%) are left out, because the story changes neither flai nor operator settings. `strategic_days` already reaches `/api/stats` (`flai/internal/metrics/strategic.go`).
  - `design/tech/charts.md` is left out: one axis keeps its rule as it stands.
- Forecast 40m, raised from flai's 12m (74 s a unit of size over only 3 medium feature stories, times size 9). Comparable chart stories took 12 to 33 agent minutes (S-0163, S-0166, S-0168, S-0169), and this one builds two charts, the waiting mean, a page group, notes, and two guides, in four tasks.
  - Delivery 2026-10-05T20:44Z. flai placed it 12th in the pull order at 16:04Z, but S-0216 waits for S-0227, forecast to be delivered at 20:04Z, and takes 40m after that. It assumes the operator finalizes and promotes it then.
- Cost of delay 46.84 USD a week, from `flai cod` after the forecast change. The story has no inputs of its own, so this is its share of E-0016's 1500 USD a week: 40m of 21h21m over the epic's 17 open stories without inputs. It stands as flai worked it out. It replaces planner-E-0016's 48.78, worked out on the earlier 1h forecast.
- Its chart siblings S-0212, S-0213, S-0214, and S-0215 change `charts.ts` and the chart page too, so the overlap hold orders them. S-0212 and S-0213 add a Planning group; the Strategic group goes after the others.
