---
id: S-0527
type: story
nature: remediation
title: "metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does"
status: backlog
owner: alex
created: 2026-10-28T12:54:11Z
updated: 2026-10-28T23:14:35Z
transitions: []
tags: [design, metrics, flaiover]
topics: [dashboard]
touches: [design/adrs, design/adrs/README.md, design/system/metrics.md, flaiover/src/lib/viz/explain.test.ts, design/issues/I-0193-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 517
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 1141
          output: 38959
          cache_read: 6905775
          cache_write: 233476
          cost: 8.3553
        - model: claude-haiku-5-5
          input: 46
          output: 24125
          cache_read: 1638273
          cache_write: 121601
          cost: 0.1134
    - kind: orchestrator
      seconds: 1068
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 42
          output: 486
          cache_read: 12899443
          cache_write: 46830
          cost: 2.8423
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-28T12:54:11Z
  value: 25
  by: planner-S-0527
  at: 2026-10-28T14:59:56Z
forecast:
  duration: 40m
  delivery: 2026-10-29T15:45:00Z
  basis: "Its own forecast of 40m; 30th in the pull order with an in-progress limit of 5, behind S-0507, S-0535, S-0532, S-0537, S-0534, S-0536, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517, S-0518, S-0519, S-0520, S-0523 and S-0524."
  by: flai
  at: 2026-10-28T23:14:35Z
finalized:
  by: alex
  at: 2026-10-28T12:57:57Z
---
# S-0527 metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does

## Goal

This story remediates [I-0193](../../../design/issues/I-0193-metrics-md-charts-says-burn-up-draws-a-forecast-line-from-throughput-and-cycle-time-filters-by-nature-but-neither-chart-does.md), "metrics.md § Charts says Burn-up draws a forecast line from throughput and Cycle Time filters by nature, but neither chart does". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0193 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0193 is closed with `flai issue close I-0193 --reason` saying what fixed it

## Tasks
- T-2255 An ADR refining ADR-0054 and ADR-0057 settles that metrics.md § Charts says what the dashboard draws, with a row for every chart it draws
- T-2256 metrics.md § Charts says what each chart draws: Cycle Time per epic coloured by nature, Burn-up scope and done, the dashed all-models line on every chart per model, and rows for Agent Waiting, Strategic Cost, and Strategic Use
- T-2257 explain.test.ts holds every chart kind to a row in metrics.md § Charts under its catalog title, which reproduces I-0193's missing Strategic Use row
- T-2258 I-0193 is closed with the reason naming the ADR and the corrected § Charts

## Notes

Cost of delay inputs set by flai from I-0193. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-22T21:09:37Z, 5.7 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0527 on 2026-10-28; the plan is on TH-0697. The story declared no touches, so `flai touches suggest S-0527` was run from the paths the issue's instance names, `design/system/metrics.md` and the explainers `flaiover/src/lib/viz/explain.ts` and `explain.test.ts`.

**The solution proposed from I-0193's instance.** The charts, their explainers (`flaiover/src/lib/viz/explain.ts`, S-0485), `docs/users/flaiover.md`, and `design/system/flaiover-dashboard.md` all agree with each other; only `metrics.md` § Charts differs, and no ADR decided a forecast line or a nature filter. Proposed: correct the contract, not the charts. One ADR (T-2255), refining ADR-0054 (metrics.md defines each chart's data and changes only with an ADR) and ADR-0057 (the dashed all-models line, decided there for Avg. Time / Model and Avg. Cost / Model alone), settles that § Charts says what the dashboard draws: Cycle Time is chosen per epic or whole and coloured by nature, with no nature filter; Burn-up draws scope and done, with no forecast line; every chart per model, Tokens / Min and Tokens / $ included, adds a dashed line for all models with two or more; Agent Waiting, Strategic Cost, and Strategic Use, which the table lacks, get rows, Strategic Use's saying the mean waiting is the dashboard's sum of `wait_threads_seconds` and `wait_review_seconds` over the items completed that day, as TH-0140 settled (T-2256); and every chart the dashboard draws has a row in § Charts, which a test holds (T-2257). Not proposed: building the forecast line and the nature filter, which are features no story asked for, and would be a feature story of their own if the operator wants them. The operator may choose otherwise on TH-0697 before T-2255 is worked.

**Where each touch came from.**

- Layout, the ADR (T-2255): `design/adrs` and `design/adrs/README.md`, the index row `flai adr new` writes; co-change put the README first (43%).
- Declared by the issue (T-2256): `design/system/metrics.md`, the § Charts table and § Strategic use per day.
- Layout, the test (T-2257): `flaiover/src/lib/viz/explain.test.ts`, which already reads `metrics.md` from the repository to check each explainer's anchor; the chart titles it needs are in `flaiover/src/lib/viz/catalog.ts`, which it reads and does not change.
- Declared by the criteria (T-2258): `design/issues/I-0193-….md` and `design/issues/summary.md`, which `flai issue close` rewrites.
- Left out, with the reason: `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` (co-change 18% and 20%) already describe the charts as drawn, and no chart changes; `flai/internal/metrics/metrics.go` and `flai/cmd/stats.go` (27%, 18%) change nothing, since no `flai stats` field changes.

**Folder touch kept.** `design/adrs`: the story adds one ADR, whose number `flai adr new` allots when the story's agent runs it, so no task can name the file yet. The folder is in the manifest's shared paths, so it holds no ready story (ADR-0096), and `flai task done` widens the touches with the committed path (ADR-0144).

**Layers.** Four, serial: T-2255; T-2256 after it, since it cites the ADR's number; T-2257 after T-2256, since the test fails until the rows exist and `flai task done` runs the tier; T-2258 after T-2256 and T-2257, since the close reason names what they did.

**Overlaps on 2026-10-28.** `design/system/metrics.md` is touched by S-0502 (in progress), so this story is held while S-0502 is in progress, and by backlog S-0425; `design/adrs/README.md`, `design/adrs`, and `design/issues` are in the manifest's shared paths and hold nothing. No `after`: the story needs no other story's change on main.

**Forecast.** `flai forecast S-0527` gave 4m with no touches (size 2) and gives 16m with them: the median 118 s per unit of size over 33 done remediation stories on claude-opus-5-5 in the medium band, times size 8 (2 criteria, 6 touches), delivery 2026-10-29T05:14Z, 36th in the pull order behind the E-0024, E-0023, and remediation stories. Raised to 40m: the size counts touches, and six touches hide four serial layers, an ADR written and accepted with `flai adr new`, a rewrite of the § Charts table with three new rows, a vitest test with the flaiover tier run, and the issue close, each committed with `flai task done`. The delivery written is flai's play-out start, 05:14Z less 16m times the cycle factor 1.35, plus 40m times that factor: 2026-10-29T05:46Z.

**Cost of delay.** 25.00 USD a week, as `flai cod S-0527` works it out from the input flai set from I-0193: `time_lost_per_cycle: 10m` at 150 USD an hour, one 168h cycle a week. The figure stands: the issue's one instance cost S-0485 the minutes it took to find that the table and the charts differ, and each explainer or chart written from § Charts until it is corrected costs the same again.
