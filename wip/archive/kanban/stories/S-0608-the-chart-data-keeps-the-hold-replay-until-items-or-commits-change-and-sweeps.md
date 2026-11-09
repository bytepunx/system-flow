---
id: S-0608
type: story
nature: improvement
title: The chart data keeps the hold replay until items or commits change and sweeps the days once, so a 365-day window costs what a 30-day one does
status: done
parent: E-0033
owner: alex
created: 2026-11-07T22:41:41Z
updated: 2026-11-09T10:32:30Z
transitions:
  - to: ready
    at: 2026-11-09T09:10:19Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T09:10:31Z
    by: agent-S-0608
  - to: review
    at: 2026-11-09T10:30:30Z
    by: agent-S-0608
  - to: done
    at: 2026-11-09T10:32:30Z
    by: orchestrator
tags: [cli]
touches: [flai/internal/metrics/claims.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/metrics.go, flai/internal/metrics/costorder.go, flai/internal/metrics/windows_test.go, flai/internal/metrics/testdata/windows.golden.json, flai/internal/metrics/claims_test.go, flai/internal/metrics/claims_bench_test.go, flai/internal/metrics/costofdelay_test.go, flai/internal/metrics/costorder_test.go, flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/workitem/store_test.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/internal/storygit/committed.go, design/system/server-performance.md, design/system/flai-cli.md, flai/internal/metrics/coordination.go, flai/internal/storygit/committed_test.go]
after: [S-0606]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4812
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 3
      work: 65
  models:
    - model: claude-opus-5-5
      input: 430
      output: 150329
      cache_read: 38100858
      cache_write: 1138905
      cost: 17.0112
  strategic:
    - kind: planner
      seconds: 529
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 201
          output: 20836
          cache_read: 19166838
          cache_write: 173506
          cost: 5.5248
    - kind: orchestrator
      seconds: 43
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 32
          cache_read: 445998
          cache_write: 9445
          cost: 0.1162
cost_of_delay:
  inputs:
    penalty_per_week: 20
    by: orchestrator
    at: 2026-11-08T00:04:56Z
  value: 20
  by: planner-S-0608
  at: 2026-11-08T00:05:51Z
forecast:
  duration: 1h
  delivery: 2026-11-09T10:33:00Z
  basis: "Its own forecast of 1h; 1st in the pull order with an in-progress limit of 5, behind S-0610, S-0612 and S-0613."
  by: flai
  at: 2026-11-09T09:10:29Z
finalized:
  by: orchestrator
  at: 2026-11-08T00:10:06Z
---
# S-0608 The chart data keeps the hold replay until items or commits change and sweeps the days once, so a 365-day window costs what a 30-day one does

## Goal

Cause 18 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `replayHolds` replays all history on every call, whatever the window (`flai/internal/metrics/claims.go`), 175 to 310 ms and about 60 MB. `costOfDelay` and `claims` loop day by day over every item, parse every transition's time again each day through `stateAt`, and start at the window's start even before the first item existed: `costOfDelay` takes 2 ms at 1 day and 170 to 217 ms at 365 days. Thread entries are parsed twice, and `costOrder` builds its forecaster three times. A warm 365-day `stats.get` takes 1,067 to 1,118 ms against 672 to 816 ms at 30 days.

## Acceptance criteria

- [x] `replayHolds` and `touchDrift` are kept between requests until the work items or the stories' commits change, and only the slicing to the window is done per request.
- [x] Each item's transition times are parsed once a call, the day-by-day sections sweep the days once, and the days before the first item are answered without a loop; thread entries are parsed once and the forecaster built once.
- [x] The report is the same as before for windows of 1, 7, 30, 90, and 365 days, which a test shows, and a warm 365-day `stats.get` on this repository takes under 1.2 times a warm 30-day one; server-performance.md records the before and after.

## Tasks
- T-2662 A golden test pins the report on testdata/good at windows of 1, 7, 30, 90, and 365 days before anything changes
- T-2668 The item store counts a generation per root, which moves when an item file is added, changed, or removed, in snapshot mode or not
- T-2672 A thread parses its entries once, and waiting and digest read the same parse
- T-2675 costOrder builds its forecaster once and plays the three orders out on it
- T-2679 Each item's transitions are parsed once a call, and costOfDelay and claims sweep the days once from the first item's day
- T-2683 The hold replay and the touch drift are kept between requests until the item generation or the commits' key moves, and each request only slices them to its window
- T-2687 Measure a warm stats.get at 30 and 365 days before and after on this repository, and record them with the kept replay in the design

## Notes

- The answer and `metrics.md` do not change.
- Keeping needs a generation count on the item store (`flai/internal/workitem/store.go`), in memory.

### Planning

Planned by planner-S-0608; the plan's thread is TH-0900. Replanned on 2026-11-08 after a touches edit logged as by `system-flow` at 2026-11-07T23:52:00Z. That edit was the first plan's own `flai touches` call, which added `coordination.go` and `committed_test.go` without an agent name. The replan keeps the touches below and the seven tasks. It rewrites T-2668 so the generation also moves in S-0606's snapshot mode, and adds `after: S-0606` for that reason.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/metrics/claims.go` | declared; `replayHolds`, `touchDrift`, and the `claims` day loop |
| `flai/internal/metrics/costofdelay.go` | declared; the `costOfDelay` day loop, `stateAt` per day, and `waited`, which parses every transition again |
| `flai/internal/metrics/metrics.go` | declared; `stateAt`, `dayRange`, and `Compute`, which calls `replayHolds` |
| `flai/internal/metrics/costorder.go` | declared; it calls `planning.PullTimes` once per order, three times |
| `flai/internal/metrics/coordination.go` | layout: it calls `stateAt` too, so it reads the parsed timeline |
| `flai/internal/metrics/windows_test.go` | layout: criterion 3's test of the report at 1, 7, 30, 90, and 365 days |
| `flai/internal/metrics/testdata/windows.golden.json` | layout: the report before the change, which that test compares against |
| `flai/internal/metrics/claims_test.go` | co-change, 12%: the claims and replay tests |
| `flai/internal/metrics/claims_bench_test.go` | layout: `BenchmarkReplayHolds` and `replayHoldsReference`, the kept replay's oracle |
| `flai/internal/metrics/costofdelay_test.go` | co-change, 9%: the day sweep's tests |
| `flai/internal/metrics/costorder_test.go` | layout: the one forecaster's test |
| `flai/internal/workitem/store.go` | declared in the Notes: the generation count |
| `flai/internal/workitem/repo.go` | layout: `Repo` is what `statsread` holds, so it answers the store's generation |
| `flai/internal/workitem/store_test.go` | layout: the generation's tests |
| `flai/internal/threads/threads.go` | layout: `Thread.Entries` runs its regex on every call; `waiting.go` and `digest.go` each call it |
| `flai/internal/threads/threads_test.go` | layout: its test |
| `flai/internal/planning/forecast.go` | layout: `PullTimes` builds `newForecaster` on every call |
| `flai/internal/planning/forecast_test.go` | layout: its test |
| `flai/internal/statsread/statsread.go` | co-change, 28%: `Read` hands `Compute` the generation and the commits' key |
| `flai/internal/statsread/statsread_test.go` | co-change, 19% |
| `flai/internal/storygit/committed.go` | layout: `Committed` keys its cache by the main branch and `show-ref --heads`; the kept replay takes that key |
| `flai/internal/storygit/committed_test.go` | layout: the key's test |
| `design/system/server-performance.md` | design: criterion 3 records the before and after under cause 18 |
| `design/system/flai-cli.md` | co-change, 9%; design: its Internal structure describes the item store and the stats caches |

Not taken from `flai touches suggest`: `flai/cmd/stats.go`, `flai/cmd/check_stats_test.go`, and `flai/internal/metrics/strategic.go`, since the CLI runs once and keeps nothing and the strategic section does not change. Also not taken: `design/system/metrics.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, and the other docs the widened set's co-change now lists, such as `docs/operators/index.md` and `design/system/flaiover-dashboard.md`, since the answer and the commands do not change. No folder touch.

`after`: S-0606. Its T-2615 gives the item store a snapshot mode in `flai/internal/workitem/store.go`, which `flai serve` turns on. In that mode a list or get stats nothing, so the only place the items change is where a change the store is told of replaces the snapshot. T-2668 reads that replace path and moves the generation there; without it, the kept replay in `flai serve` would never see an item change. S-0596, which wrote cause 18's section, is done, so the section is on main. S-0607 changes `metrics.go` and `statsread.go`, and S-0600 changes `store.go`, but this story reads no output of either; the trial merge judges the overlap.

Figures:

- Forecast 1h, adjusted from flai's 43m (94 s per unit of size over 124 done improvement stories, times size 27, its 24 touches counted). The story has seven tasks: a golden test over five windows, a generation on the store, two parse fixes, the day sweep, a cache kept between requests, and a measurement on the operator's host. S-0606's comparable plan of five tasks is forecast at 45m.
- Delivery 2026-11-08T03:30:00Z: flai's delivery for its 8th place in the pull order, 03:13, moved by the 17 minutes added to the duration. S-0606 is in progress and ahead of it in the pull order, so the `after` moves no delivery.
- Cost of delay 20 USD a week, as `flai cod` gives it from the input `penalty_per_week: 20`, which the orchestrator set at 2026-11-08T00:04:56Z from TH-0888's table. Kept as flai gives it: no revenue rides on it, and long chart windows are its only cost.

### Accepted by the orchestrator

- Verified: f54809c9b2d0ac2ef09e2c431ddc900667fb91ce
- At: 2026-11-09T10:32:30Z

Verdict: pass
- 1: flai/internal/metrics/claims.go, flai/internal/metrics/metrics.go, flai/internal/workitem/store.go, flai/internal/workitem/repo.go, flai/internal/statsread/statsread.go, flai/internal/storygit/committed.go, flai/internal/metrics/claims_test.go, flai/internal/workitem/store_test.go, flai/internal/statsread/statsread_test.go, flai/internal/storygit/committed_test.go
- 2: flai/internal/metrics/metrics.go, flai/internal/metrics/costofdelay.go, flai/internal/metrics/claims.go, flai/internal/metrics/coordination.go, flai/internal/threads/threads.go, flai/internal/planning/forecast.go, flai/internal/metrics/costorder.go, flai/internal/metrics/costofdelay_test.go, flai/internal/metrics/costorder_test.go, flai/internal/threads/threads_test.go, flai/internal/planning/forecast_test.go
- 3: flai/internal/metrics/windows_test.go, flai/internal/metrics/testdata/windows.golden.json, design/system/server-performance.md, design/system/flai-cli.md
