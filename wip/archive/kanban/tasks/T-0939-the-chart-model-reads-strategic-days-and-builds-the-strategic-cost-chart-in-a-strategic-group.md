---
id: T-0939
type: task
nature: feature
title: The chart model reads strategic_days and builds the Strategic Cost chart in a Strategic group
status: done
parent: S-0216
owner: alex
created: 2026-10-05T05:45:27Z
updated: 2026-10-07T09:31:53Z
transitions:
  - to: ready
    at: 2026-10-07T09:18:43Z
    by: agent-S-0216
  - to: in-progress
    at: 2026-10-07T09:18:44Z
    by: agent-S-0216
  - to: done
    at: 2026-10-07T09:31:53Z
    by: agent-S-0216
stream: S-0216
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
usage:
  source: log
  seconds: 789
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 79
      output: 32738
      cache_read: 3839280
      cache_write: 250795
      cost: 2.8464
---
# T-0939 The chart model reads strategic_days and builds the Strategic Cost chart in a Strategic group

## Work

- Declare `strategic_days[]` on the `Report` type in `flaiover/src/lib/viz/charts.ts` as `design/system/metrics.md` § Strategic use per day defines it: `date`, `agents` per kind (`cost`, `seconds`, `estimated`), `cost`, `seconds`, `completed`, `cost_per_item`, `cycle_time_seconds`. `/api/stats` already carries it.
- Add a `STRATEGIC_KINDS` group with `strategic-cost` and `strategic-use`, their titles, and `strategic-cost` in `build()`.
- `strategic-cost`: one bar per day of the window, the planner, orchestrator, and analyzer stacked with their fixed colours, and an estimated cost marked as an estimate. The mean cost per story completed that day (`cost_per_item`) is a line on the same y-axis, in dollars. TH-0140 settled one axis, as `design/tech/charts.md` rules. State the ratio over the window: strategic cost per story completed, as a share of the mean agent cost per story.
- Days with nothing completed have no point on the line. The axis spans the window (ADR-0054).
- First layer: the other tasks build on the type and the group.

## Done when

- `build('strategic-cost', …)` maps a fixture `strategic_days` to the stacked bars, the line, and the ratio. Tests in `charts.test.ts` pin the figures against the fixture, which is `flai stats --json` output.
- The existing test "every kind builds with one y-axis and a tooltip" passes with the new kinds, unchanged.
- `scripts/flaiover-test.sh` passes.

## Notes
