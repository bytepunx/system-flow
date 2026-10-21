---
id: S-0471
type: story
nature: remediation
title: The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds
status: done
owner: alex
created: 2026-10-21T11:58:35Z
updated: 2026-10-21T13:25:39Z
transitions:
  - to: ready
    at: 2026-10-21T12:10:31Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T12:13:10Z
    by: agent-S-0471
  - to: review
    at: 2026-10-21T13:05:30Z
    by: agent-S-0471
  - to: done
    at: 2026-10-21T13:25:39Z
    by: orchestrator
tags: [flai]
topics: [metrics, performance]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/claims_test.go, flai/internal/metrics/claims_bench_test.go, flai/internal/workitem/hold.go, design/system/server-performance.md, design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md, design/issues/summary.md, flai/internal/workitem/hold_test.go, design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3153
  priced_by: harness
  turns:
    - day: 2026-10-21
      hand_edits: 4
      work: 29
  models:
    - model: claude-opus-5-5
      input: 394
      output: 165524
      cache_read: 18277379
      cache_write: 650781
      cost: 10.9318
  strategic:
    - kind: orchestrator
      seconds: 308
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 274
          cache_read: 5065158
          cache_write: 17771
          cost: 1.1223
        - model: claude-sonnet-5-5
          input: 4
          output: 17
          cache_read: 47873
          cache_write: 15350
          cost: 0.0435
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4s
    by: flai
    at: 2026-10-21T11:58:35Z
  value: 0.17
  by: planner-S-0471
  at: 2026-10-21T12:09:14Z
forecast:
  duration: 1h30m
  delivery: 2026-10-22T06:01:00Z
  basis: "flai forecast's 18m (115 s per unit of size over 27 done medium remediation stories, size 9) raised to 1h30m for a profile, a history generator, a reference replay, an algorithm change, and a measurement on this repository; delivery is flai's 2026-10-22T04:49Z, 35th in the pull order, plus the 72m added"
  by: planner-S-0471
  at: 2026-10-21T12:09:14Z
finalized:
  by: alex
  at: 2026-10-21T12:03:57Z
---
# S-0471 The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds

## Goal

This story remediates [I-0155](../../../design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md), "The chart data takes about 4 s warm, against 10 ms measured on 2026-09-29, almost all of it in metrics.replayHolds". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0155 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0155 is closed with `flai issue close I-0155 --reason` saying what fixed it

## Tasks
- T-1991 A benchmark and a reference test replay the holds over a generated history of this repository's size, and the profile names the fix
- T-1993 replayHolds judges the holds again only at an event that can change one, over only the open stories, giving the reference replay's result
- T-1995 stats.get and flai stats are measured again on this repository, server-performance.md records the regression and its fix, and I-0155 is closed

## Notes

Cost of delay inputs set by flai from I-0155. time_lost_per_cycle 4s: 4s per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-16T06:09:48Z, 5.2 days before this story; under one cycle counts as one).

### Planning

Cause, from the code: `replayHolds` in `flai/internal/metrics/claims.go` came with S-0205 on 2026-10-03, after the 10 ms measurement of 2026-09-29. At every distinct event time since the first item it builds a new `workitem.Holds` over every item created so far, works out every in-progress story's claim again, and judges every ready story. Its cost grows with the square of the history. Nothing is memoized, and no benchmark covers it.

Touches, all files, no folder touch:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/metrics/claims.go` | design | I-0155's profile names `metrics.replayHolds`, defined here |
| `flai/internal/metrics/claims_test.go` | layout | the replay's tests; the equivalence and counting tests go here |
| `flai/internal/metrics/claims_bench_test.go` | layout | a new file for the generator, the benchmarks, and the reference replay |
| `flai/internal/workitem/hold.go` | design | `Holds.Claim`, `Holds.Overlaps`, and `PathsOverlap` are under the profile; the replay may need a way to build `Holds` from claims already worked out |
| `design/system/server-performance.md` | design | holds the 2026-09-29 figures the issue compares against |
| I-0155's file under `design/issues/` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

`flai touches suggest` listed co-change paths, the busiest `design/system/flai-cli.md` at 7 percent, and close-out issue files. None bears on the replay, so none was taken. `docs/` and `design/system/metrics.md` stay untouched because the fix changes no output.

Forecast: flai gives 18m, 115 s per unit of size over 27 done medium remediation stories, size 9. Raised to 1h30m: this story needs a profile, a generator at this repository's scale, a reference replay, an algorithm change proved equal to it, and a measurement on the repository. Delivery moves from flai's 2026-10-22T04:49Z by the same 72m.

Cost of delay: `flai cod` gives 0.17 USD a week from the operator's input of 4s lost per 168h cycle. The value stands as computed. The input counts one occurrence, while every charts page load and every agent's `flai stats` pays the 4 s; the plan's thread says so.

Tasks, in three layers:

1. T-1991: the generator, the benchmarks, and the reference replay; the profile names the fix
2. T-1993, after T-1991: the fix in `replayHolds`, with the equivalence and counting tests
3. T-1995, after T-1993: measure on this repository, record it in `server-performance.md`, and close I-0155

### Accepted by the orchestrator

- Verified: 0eb3945f66f1c4001d7fac2fc7a21f28b3afa5ef
- Batch: 20261021-131344
- Enqueued verified: 27d62cb78b70f53dd1a35c581db5377703a332a7
- At: 2026-10-21T13:25:39Z

Verdict: pass. Both criteria are met, the replay's output is unchanged against the reference, and every changed file is within the story's touches.
- 1: flai/internal/metrics/claims.go, flai/internal/metrics/claims_bench_test.go, flai/internal/metrics/claims_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, design/system/server-performance.md
- 2: design/issues/I-0155-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md, design/issues/summary.md
