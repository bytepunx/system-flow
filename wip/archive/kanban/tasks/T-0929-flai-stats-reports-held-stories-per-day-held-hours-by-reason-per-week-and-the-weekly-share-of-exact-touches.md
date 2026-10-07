---
id: T-0929
type: task
nature: feature
title: flai stats reports held stories per day, held hours by reason per week, and the weekly share of exact touches
status: done
parent: S-0214
owner: alex
created: 2026-10-05T05:45:09Z
updated: 2026-10-07T08:33:29Z
transitions:
  - to: ready
    at: 2026-10-07T08:23:19Z
    by: agent-S-0214
  - to: in-progress
    at: 2026-10-07T08:23:19Z
    by: agent-S-0214
  - to: done
    at: 2026-10-07T08:33:24Z
    by: agent-S-0214
stream: S-0214
tags: [flai]
touches: [flai/internal/metrics, flai/internal/statsread/statsread.go]
after: [T-0924]
usage:
  source: log
  seconds: 605
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 48
      output: 20449
      cache_read: 2364967
      cache_write: 88680
      cost: 1.4006
---
# T-0929 flai stats reports held stories per day, held hours by reason per week, and the weekly share of exact touches

## Work

In `flai/internal/metrics`, compute what T-0924 defines. Extend the replay in `heldSeconds` (`claims.go`) to keep each hold's first reason (`workitem.HoldOverlap`, `HoldAfter`, `HoldNoTouches`). From it, sum the seconds per ISO week of the window and count the stories held at each day's end into `claims.days[].held`, over every story in ready at the time, not only those of the report. From `claims.drift`, give the weekly count and share of stories whose touches were exact. Waits for T-0924, whose definitions this implements.

## Done when

- [ ] `Compute` reports `claims.days[].held`, `claims.weeks[]` with held seconds per reason, and the weekly exact-touches share, as `metrics.md` defines them
- [ ] Tests in `claims_test.go` pin each value on a fixture, to the second, including a hold with two reasons and a week with no drift
- [ ] `scripts/flai-test.sh` passes

## Notes
