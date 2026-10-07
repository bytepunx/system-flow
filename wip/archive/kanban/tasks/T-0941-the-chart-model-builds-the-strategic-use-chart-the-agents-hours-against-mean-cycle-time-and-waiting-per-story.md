---
id: T-0941
type: task
nature: feature
title: "The chart model builds the Strategic Use chart: the agents' hours against mean cycle time and waiting per story"
status: done
parent: S-0216
owner: alex
created: 2026-10-05T05:45:33Z
updated: 2026-10-07T09:43:28Z
transitions:
  - to: ready
    at: 2026-10-07T09:32:19Z
    by: agent-S-0216
  - to: in-progress
    at: 2026-10-07T09:32:19Z
    by: agent-S-0216
  - to: done
    at: 2026-10-07T09:43:28Z
    by: agent-S-0216
stream: S-0216
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
after: [T-0939]
usage:
  source: log
  seconds: 669
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 55
      output: 22774
      cache_read: 2670826
      cache_write: 174468
      cost: 1.9801
---
# T-0941 The chart model builds the Strategic Use chart: the agents' hours against mean cycle time and waiting per story

## Work

- Waits for T-0939: it changes the same file and needs its `Report` type and Strategic group.
- Declare `wait_threads_seconds` and `wait_review_seconds` on the item metrics type in `flaiover/src/lib/viz/charts.ts` (`design/system/metrics.md` § Waiting). `/api/stats` already carries them.
- `strategic-use`: one bar per day of the window, the hours of the planner, orchestrator, and analyzer stacked (`strategic_days[].agents[].seconds`).
- Two lines over the stories completed that day: the mean cycle time (`cycle_time_seconds`), and the mean waiting time. The waiting time is the mean of `wait_threads_seconds` plus `wait_review_seconds` over `items[]` completed that day, a missing wait counted as zero. Draw both on the bars' y-axis, in hours: TH-0140 settled one axis.
- Add it to `build()`.

## Done when

- `build('strategic-use', …)` maps a fixture to the stacked bars and both lines on one y-axis. Tests in `charts.test.ts` pin the waiting mean for a day against the fixture's items. A day with no story completed has no point on the lines.
- `scripts/flaiover-test.sh` passes.

## Notes
