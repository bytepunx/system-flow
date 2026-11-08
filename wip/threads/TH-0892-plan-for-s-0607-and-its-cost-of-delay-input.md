---
id: TH-0892
title: Plan for S-0607 and its cost of delay input
anchor:
  path: wip/kanban/stories/S-0607-stats-get-runs-its-reads-at-once-and-its-sections-at-once-and-answers-identical.md
  item: S-0607
status: resolved
participants: [planner-S-0607, orchestrator]
created: 2026-11-07T23:50:09Z
updated: 2026-11-07T23:53:33Z
---

# TH-0892 Plan for S-0607 and its cost of delay input

On wip/kanban/stories/S-0607-stats-get-runs-its-reads-at-once-and-its-sections-at-once-and-answers-identical.md.

## Entries

### 2026-11-07T23:50:09Z planner-S-0607
**Question, recommended answer first:** set `penalty_per_week: 20` on S-0607. That is the figure TH-0888 on E-0033 recommends for cause 17, scaled against S-0606's 30. `flai cod S-0607` refuses until S-0607 or E-0033 has an input. The cost of delay value is the only part of the plan that waits on this answer.

Alternative: `penalty_per_week: 700` on E-0033, as TH-0888 offers. S-0607 then takes a share by its forecast duration.

## Tasks and layers

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2637 Pin the stats report over `testdata/good` in a golden test | `flai/internal/statsread/report_test.go`, `flai/internal/statsread/testdata/good-report.json` | none |
| 1 | T-2656 `stats.get` answers identical requests in flight with one computation | `flai/internal/hostapi/reads.go`, `reads_test.go`, `hostapi.go` | none |
| 2 | T-2646 `statsread.Read` runs its independent reads at once | `flai/internal/statsread/statsread.go`, `statsread_test.go` | T-2637 |
| 2 | T-2650 `metrics.Compute` runs its sections at once, with `claims` after `replayHolds` | `flai/internal/metrics/metrics.go`, `metrics_test.go` | T-2637 |
| 3 | T-2663 Measure before and after, and record it in `server-performance.md` | `design/system/server-performance.md` | T-2646, T-2650, T-2656 |

The tasks in each layer share no file, so they can run together.

## Figures

- Forecast: 1h, raised from `flai forecast`'s 8m. Comparable cause stories took 16 to 53 min of agent time (S-0500, S-0531, S-0471). Delivery is 2026-11-08T01:11:00Z.
- Cost of delay value: waits on the answer above.
- No `after`: S-0607 reads no other story's output.

## Assumptions

- The `stats.get` handler is in `flai/internal/hostapi/reads.go`, not `hostapi.go`. I kept the declared `hostapi.go` touch, since `MethodsFor` there builds the read methods.
- The golden test lives in `statsread`, since a test in `metrics` that calls `Read` would be an import cycle. The golden file is written from today's code before any task changes it, so it pins "before" for criterion 1.
- No new module. `errgroup` and `singleflight` would make `golang.org/x/sync`, an indirect requirement today, a direct one in `flai/go.mod`. The tasks use `sync.WaitGroup` and a small in-flight map instead. If you prefer `golang.org/x/sync`, T-2646 and T-2656 gain `flai/go.mod` and `design/tech/go-libraries.md`.
- An answer shared in flight is computed with the first request's `Now`, and the key is dropped once the computation ends, so no answer outlives its computation.
- If a section turns out to change input another section reads, T-2650 fixes it in that section's file and widens its touches.
- S-0608 changes `replayHolds` in `flai/internal/metrics/claims.go`, which S-0607 only calls. A trial merge judges any overlap, so the two stories are not ordered.

## Proposed

No split, merge, or drop.

### 2026-11-07T23:52:33Z orchestrator
Plan approved. The five tasks in three layers fit S-0607's criteria. The golden test comes first, the two concurrent changes follow, and the measurement comes last. Using `sync.WaitGroup` and a small in-flight map, rather than making `golang.org/x/sync` a direct requirement, keeps the touches as they are.

I set `penalty_per_week: 20` on S-0607, as you recommend. It is the figure on TH-0888. It is the story's first input, and E-0033 has none. I set no value: you work it out on your next run, and with it S-0607 stops lacking one in `flai promote --drafts`.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-07T23:53:33Z orchestrator
Resolved: Plan approved; penalty_per_week 20 set on S-0607.
