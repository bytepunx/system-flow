---
id: S-0607
type: story
nature: improvement
title: stats.get runs its reads at once and its sections at once, and answers identical requests in flight once
status: ready
parent: E-0033
owner: alex
created: 2026-11-07T22:41:34Z
updated: 2026-11-08T01:13:21Z
transitions:
  - to: ready
    at: 2026-11-08T00:11:20Z
    by: orchestrator
tags: [cli]
touches: [flai/internal/statsread/statsread.go, flai/internal/metrics/metrics.go, flai/internal/hostapi/hostapi.go, flai/internal/statsread/statsread_test.go, flai/internal/statsread/report_test.go, flai/internal/statsread/testdata/good-report.json, flai/internal/metrics/metrics_test.go, flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, design/system/server-performance.md]
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
      seconds: 46
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 36
          cache_read: 451609
          cache_write: 17334
          cost: 0.1196
cost_of_delay:
  inputs:
    penalty_per_week: 20
    by: orchestrator
    at: 2026-11-07T23:51:21Z
  value: 20
  by: planner-E-0033
  at: 2026-11-08T00:02:05Z
forecast:
  duration: 1h
  delivery: 2026-11-08T02:34:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 5, behind S-0599, S-0600, S-0601, S-0603 and S-0606."
  by: flai
  at: 2026-11-08T01:13:21Z
finalized:
  by: orchestrator
  at: 2026-11-08T00:09:39Z
---
# S-0607 stats.get runs its reads at once and its sections at once, and answers identical requests in flight once

## Goal

Cause 17 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `stats.get` runs its 11 reads one after another (`flai/internal/statsread/statsread.go` `Read`), 230 to 300 ms warm, then its sections one after another (`flai/internal/metrics/metrics.go` `Compute`). It takes 0.63 to 1.1 s warm and 1.8 to 4 s cold, allocating about 200 MB a warm call. The dashboard asks again on every item change, and of the 120 slow calls in the serve log, 73 came within 3 s of the one before, each computed in full.

## Acceptance criteria

- [ ] `statsread.Read` runs its independent reads at once, and `Compute` runs its sections at once, with `claims` after `replayHolds`, which it reads; the report is the same as before, which a test over `flai/internal/metrics/testdata/good` shows, run under `-race`.
- [ ] Two `stats.get` requests with the same parameters while one is being computed get one computation's answer.
- [ ] A warm 30-day `stats.get` on this repository takes at least 250 ms less than before, and server-performance.md records the before and after per window.

## Tasks
- T-2637 Pin the stats report over testdata/good in a golden test before anything runs at once
- T-2646 statsread.Read runs its independent reads at once
- T-2650 metrics.Compute runs its sections at once, with claims after replayHolds
- T-2656 stats.get answers identical requests in flight with one computation
- T-2663 Measure stats.get before and after on this repository and record it in server-performance.md

## Notes

- Estimated from the profile: reads about 230 to 90 ms, the slowest read; compute about 450 to 280 ms at 30 days.
- The answer and `metrics.md` do not change.

### Planning

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/statsread/statsread.go` | declared | `Read` runs its reads at once |
| `flai/internal/metrics/metrics.go` | declared | `Compute` runs its sections at once |
| `flai/internal/hostapi/hostapi.go` | declared | `MethodsFor` builds the read methods; the in-flight group may live there |
| `flai/internal/hostapi/reads.go` | layout | the `stats.get` handler is here, not in `hostapi.go`; the in-flight sharing goes here |
| `flai/internal/hostapi/reads_test.go` | co-change, layout | already tests `stats.get`; the sharing test goes here |
| `flai/internal/statsread/statsread_test.go` | co-change (7 of 64) | `Read`'s own tests |
| `flai/internal/statsread/report_test.go` | design | new: the golden test criterion 1 asks for, `Read` then `Compute` over `testdata/good` |
| `flai/internal/statsread/testdata/good-report.json` | design | new: the report as the code gives it before the change |
| `flai/internal/metrics/metrics_test.go` | co-change, layout | `Compute`'s tests over `testdata/good` |
| `design/system/server-performance.md` | design | criterion 3: the before and after per window |

Left out of `touches suggest`'s list: `flai/cmd/stats.go` and `docs/users/flai.md`, since `flai stats` calls the same `Read` and `Compute` and nothing it prints changes; `design/system/metrics.md`, which the Notes say does not change. No folder touch. No `after`: no other story writes what this one reads. S-0608 changes `replayHolds` in `flai/internal/metrics/claims.go`, which this story only calls; a trial merge judges any overlap.

Figures:

- Forecast: `flai forecast` gave 8m, 79 s a unit times size 6 (3 criteria, 3 touches). Raised to 1h. The story has three separate concurrent changes, each needing a race check, a golden test, and a before-and-after measurement on this repository. Comparable cause stories took 16 to 53 min of agent time: S-0500 953 s, S-0531 2,191 s, S-0471 3,153 s. Delivery moved by the same 52 min, to 2026-11-08T01:11:00Z.
- Cost of delay: no value yet. `flai cod` refuses, because neither S-0607 nor E-0033 has inputs. TH-0888 on E-0033 recommends `penalty_per_week: 20` for this story. I asked again on the plan thread on S-0607.
