---
id: S-0337
type: story
nature: feature
title: flai stats reports the conversations between agents, the conflicts found at sync and at acceptance, and the hold time shares saved, per week
status: backlog
parent: E-0018
owner: alex
created: 2026-10-07T20:11:31Z
updated: 2026-10-07T20:28:28Z
transitions: []
tags: [flai, flaiover]
topics: [cli, dashboard, git]
touches: [design/adrs, design/system/metrics.md, flai/internal/storygit/conflicts.go, flai/internal/storygit/conflicts_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/internal/metrics/coordination.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md]
after: [S-0332, S-0333, S-0334]
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
      seconds: 300
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 45
          output: 818
          cache_read: 3661514
          cache_write: 10355
          cost: 0.9048
cost_of_delay:
  value: 139.53
  by: planner-E-0018
  at: 2026-10-07T20:22:19Z
forecast:
  duration: 54m
  delivery: 2026-10-08T11:02:00Z
  basis: "Its own forecast of 54m; 44th in the pull order with an in-progress limit of 3, behind S-0329, S-0330, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0280, S-0287, S-0288, S-0289, S-0290, S-0291, S-0297, S-0304, S-0305, S-0306, S-0309, S-0310, S-0312, S-0313, S-0314, S-0315, S-0316, S-0317, S-0318, S-0319, S-0320, S-0321, S-0322, S-0323, S-0324, S-0325, S-0326, S-0327, S-0331, S-0332, S-0333, S-0334, S-0335 and S-0336."
  by: flai
  at: 2026-10-07T20:28:28Z
finalized:
  by: orchestrator
  at: 2026-10-07T20:24:19Z
---
# S-0337 flai stats reports the conversations between agents, the conflicts found at sync and at acceptance, and the hold time shares saved, per week

## Goal

Measure whether E-0018 does what it set out to: more stories running at once with less conflict work at the end. Report, per week, the conversations between agents and how they ended, the conflicts the trial merge found while both stories were open against those a rebase met at `flai stream sync` after an acceptance or at `flai accept`, and the hold time that shares cut short, beside the held time `flai stats` already reports (ADR-0113).

## Acceptance criteria

- [ ] An ADR records the new values, and `design/system/metrics.md` defines each: conversations opened, closed by agreement, escalated, and closed with their story, per week; conflicts found by a trial merge, and conflicts met by a rebase at sync or at acceptance, per week; and stories started on a share, with the hold time each would have waited, per week.
- [ ] `flai stats --json` reports them under `coordination`, from the message files, the sync and acceptance logs, and the board, and `flai stats` prints them.
- [ ] A rebase that stops on conflicts at `flai stream sync` or `flai accept` is recorded with its story, the paths, and the time, where `flai stats` reads it.
- [ ] The dashboard's charts page charts the weekly values over the window chosen (ADR-0054).
- [ ] `docs/users/flai.md` and `docs/users/flaiover.md` describe them.

## Tasks

- T-1216 An ADR and metrics.md define the coordination values flai stats reports per week
- T-1217 A rebase that stops on conflicts at flai stream sync or flai accept is logged with its story, paths, and time
- T-1218 flai stats reports the coordination values under coordination, in JSON and in print
- T-1219 The dashboard charts the weekly coordination values over the window chosen
- T-1220 The guides and the dashboard design describe the coordination values and their chart

## Notes

### Planning

Planned by planner-E-0018 on 2026-10-07. It waits for S-0332, S-0333, and S-0334 (`after`), whose conversations, notices, and shares it counts. `metrics.md` is the contract with the dashboard and changes only with an ADR, so T-1216 comes first.

Layers, one task each, since each reads what the one before writes:

1. T-1216, the ADR and `metrics.md`.
2. T-1217, the conflict log.
3. T-1218, `flai stats`.
4. T-1219, the chart.
5. T-1220, the docs.

Touches:

- **Declared:** none before planning.
- **Layout:**
  - `flai/internal/storygit/conflicts.go` and its test: a new log beside `sync.go`; `sync.go`, `flai/cmd/accept.go`, and their conflict tests write it.
  - `flai/internal/statsread/statsread.go` and its test: `Read` loads the threads and activities for `flai stats` and the host API's `stats.get`, so it loads the messages and the conflict log too.
  - `flai/internal/metrics/coordination.go` and its test, as `waiting.go` and `claims.go` are; `metrics.go` gathers them.
  - `flai/cmd/stats.go` and `check_stats_test.go`: the printed and JSON stats.
  - `flaiover/src/lib/viz/charts.ts`, `flaiover/src/routes/charts/[kind]/+page.svelte`, and their tests: the charts.
- **Co-change:** `flai touches suggest` from `waiting.go`, `claims.go`, and `stats.go` gave `check_stats_test.go` and `metrics.go` (57%), `docs/users/flai.md` (48%), `design/system/metrics.md` (43%), `docs/users/flai-reference.md` (38%), and `statsread.go` (14%).
- **Design:** `design/system/metrics.md` § Planning, waiting, and claims and `design/system/flaiover-dashboard.md`.
- **Folder touch kept:** `design/adrs`, for T-1216's ADR. Inside `claims.shared`.
- **Not taken:** `flai/internal/hostapi/reads.go` (14%): `stats.get` returns what `statsread.Read` and `metrics` compute, so it needs no change. `flaiover/src/lib/server/stats.ts` may need the type; T-1219 widens its touches if so.

Forecast 54m, delivery 2026-10-08T12:08Z: `flai forecast` gave it, 114 s per unit over 29 done large-band feature stories, times size 28 (5 criteria, 23 touches). It stands.

Cost of delay 139.53 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 54m of 6h27m. It stands.
