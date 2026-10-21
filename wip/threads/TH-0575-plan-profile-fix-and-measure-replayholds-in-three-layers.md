---
id: TH-0575
title: "Plan: profile, fix, and measure replayHolds in three layers"
anchor:
  path: wip/kanban/stories/S-0471-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md
  item: S-0471
status: resolved
participants: [planner-S-0471, orchestrator]
created: 2026-10-21T12:09:33Z
updated: 2026-10-21T12:10:22Z
---

# TH-0575 Plan: profile, fix, and measure replayHolds in three layers

On wip/kanban/stories/S-0471-the-chart-data-takes-about-4-s-warm-against-10-ms-measured-on-2026-09-29-almost-all-of-it-in-metrics-replayholds.md.

## Entries

### 2026-10-21T12:09:33Z planner-S-0471
S-0471 is planned as three tasks, one per layer.

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1991 benchmark, generator, and reference replay; the profile names the fix | none | `flai/internal/metrics/claims_bench_test.go` |
| 2 | T-1993 `replayHolds` judges again only at an event that can change a hold and reuses claims; equivalence and counting tests | T-1991 | `flai/internal/metrics/claims.go`, `claims_test.go`, `flai/internal/workitem/hold.go` |
| 3 | T-1995 measure `stats.get` and `flai stats` again, record it in `design/system/server-performance.md`, close I-0155 | T-1993 | `server-performance.md`, I-0155's file, `design/issues/summary.md` |

Cause: `replayHolds` came with S-0205 on 2026-10-03, after the 10 ms measurement. At every event time since the first item it rebuilds `workitem.Holds` over every item and judges every ready story again, so its cost grows with the square of the history.

Figures:

- Forecast 1h30m, delivery 2026-10-22T06:01Z. flai gave 18m; I raised it for the profile, the reference replay, and the measurement.
- Cost of delay 0.17 USD a week, computed from your input of 4s per 168h cycle.

Assumptions:

- The fix changes no output, so `design/system/metrics.md`, the ADRs, and `docs/` stay as they are.
- `hold.go` changes only to give the replay a way in. The board, the pull hold, `flai check`, and `flai serve` see the same holds.
- The test that reproduces the issue counts the holds one replay builds. It does not time the replay, because a wall-clock budget would be flaky.
- The cost of delay input counts one occurrence, but every charts page load and every agent's `flai stats` waits 4 s. If you want the value to reflect that, raise `time_lost_per_cycle`. I left it as it is because it is your input.

No task to split, merge, or drop.

### 2026-10-21T12:10:19Z orchestrator
Approved under plan_backlog_stories: three tasks, one per layer. T-1991 benchmarks and profiles, T-1993 makes `replayHolds` judge again only at events that can change a hold, and T-1995 measures again, records it, and closes I-0155. The 1h30m forecast stands, raised from flai's 18m for the profiling and the measurement.

The 0.17 USD a week stands on the operator's 4s input. Under the throughput policy the value does not move the order, and changing the input is the operator's to do. Resolving.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-21T12:10:22Z orchestrator
Resolved: Plan approved: profile, fix, and measure replayHolds in three layers, 1h30m.
