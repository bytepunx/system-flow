---
id: S-0337
type: story
nature: feature
title: flai stats reports the conversations between agents, the conflicts found at sync and at acceptance, and the hold time shares saved, per week
status: in-progress
parent: E-0018
owner: alex
created: 2026-10-07T20:11:31Z
updated: 2026-10-08T11:04:51Z
transitions:
  - to: ready
    at: 2026-10-08T08:51:33Z
    by: alex
  - to: in-progress
    at: 2026-10-08T11:04:51Z
    by: agent-S-0337
tags: [flai, flaiover]
topics: [cli, dashboard, git]
touches: [design/adrs, design/system/metrics.md, flai/internal/storygit/conflicts.go, flai/internal/storygit/conflicts_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/accept.go, flai/cmd/accept_conflict_test.go, flai/internal/metrics/coordination.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [S-0332, S-0333, S-0334]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 364
  turns:
    - day: 2026-10-08
      ceremony: 1
      work: 31
    - day: 2026-10-09
      work: 18
    - day: 2026-10-10
      ceremony: 1
      work: 1
  models:
    - model: claude-opus-5-5
      input: 126
      output: 20607
      cache_read: 4718989
      cache_write: 375366
      cost: 4.3594
  strategic:
    - kind: orchestrator
      seconds: 443
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 91
          output: 1283
          cache_read: 9011006
          cache_write: 77496
          cost: 2.2394
        - model: claude-sonnet-5-5
          input: 6
          output: 36
          cache_read: 74126
          cache_write: 23071
          cost: 0.0838
cost_of_delay:
  value: 277.17
  by: planner-E-0018
  at: 2026-10-08T04:33:22Z
forecast:
  duration: 51m
  delivery: 2026-10-08T11:53:00Z
  basis: "Its own forecast of 51m; 2nd in the pull order with an in-progress limit of 5, behind S-0232, S-0338 and S-0342."
  by: flai
  at: 2026-10-08T10:52:25Z
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

Planned by planner-E-0018 on 2026-10-07, revisited on 2026-10-08. It waits for S-0332, S-0333, and S-0334 (`after`), whose conversations, notices, and shares it counts; S-0332 and S-0333 are done, so S-0334 alone holds it. It does not wait for S-0336 or S-0338: it reads the message files, not the dashboard's views of them. `metrics.md` is the contract with the dashboard and changes only with an ADR, so T-1216 comes first.

Layers, one task each, since each reads what the one before writes:

1. T-1216, the ADR and `metrics.md`.
2. T-1217, the conflict log.
3. T-1218, `flai stats`.
4. T-1219, the chart.
5. T-1220, the docs.

Touches:

- **Declared:** all kept.
- **Layout:**
  - `flai/internal/storygit/conflicts.go` and its test: a new log beside `sync.go`; `sync.go`, `flai/cmd/accept.go`, and their conflict tests write it.
  - `flai/internal/statsread/statsread.go` and its test: `Read` loads the threads and activities for `flai stats` and the host API's `stats.get`, so it loads the messages and the conflict log too.
  - `flai/internal/metrics/coordination.go` and its test, as `waiting.go` and `claims.go` are; `metrics.go` gathers them.
  - `flai/cmd/stats.go` and `check_stats_test.go`: the printed and JSON stats.
  - `flaiover/src/lib/viz/charts.ts`, `flaiover/src/routes/charts/[kind]/+page.svelte`, and their tests: the charts.
- **Co-change:** `flai touches suggest` from `waiting.go`, `claims.go`, and `stats.go` gave `check_stats_test.go` and `metrics.go` (57%), `docs/users/flai.md` (48%), `design/system/metrics.md` (43%), `docs/users/flai-reference.md` (38%), and `statsread.go` (14%). Run again on 2026-10-08 over the whole claim, it gave `design/system/flai-cli.md` (41%), which describes `flai stats`; added to the story and to T-1220.
- **Design:** `design/system/metrics.md` § Planning, waiting, and claims and `design/system/flaiover-dashboard.md`.
- **Folder touch kept:** `design/adrs`, for T-1216's ADR, whose file name no task can know before it is written. Inside `claims.shared`, so it holds nothing.
- **Not taken:** `flai/internal/hostapi/reads.go` and `writes.go`: `stats.get` returns what `statsread.Read` and `metrics` compute, so it needs no change. `design/system/workflow.md` (11%): the values are defined in `metrics.md`. `flaiover/src/lib/server/stats.ts` may need the type; T-1219 widens its touches if so.

Forecast 51m: `flai forecast` gave it on 2026-10-08, 104 s per unit over 33 done large-band feature stories, times size 29 (5 criteria, 24 touches). It stands; the 54m of the first plan came from 114 s per unit over 29 stories.

Cost of delay 277.17 USD a week: `flai cod` gave its share of E-0018's 1000 USD a week, 51m of the 3h4m forecast over the four open stories (S-0334, S-0336, S-0337, S-0338). It stands.
