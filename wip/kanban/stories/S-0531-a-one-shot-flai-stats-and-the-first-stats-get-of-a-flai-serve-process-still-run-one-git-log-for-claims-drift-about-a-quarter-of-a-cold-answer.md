---
id: S-0531
type: story
nature: improvement
title: A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer
status: in-progress
owner: alex
created: 2026-10-28T12:54:38Z
updated: 2026-10-28T18:25:43Z
transitions:
  - to: ready
    at: 2026-10-28T18:25:33Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-28T18:25:43Z
    by: agent-S-0531
tags: [flai]
topics: [metrics, performance, cli]
touches: [flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, flai/internal/statsread/statsread_test.go, design/system/server-performance.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/runbooks/backup.md, design/issues/I-0209-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md, design/issues/summary.md]
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
      seconds: 188
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 400
          cache_read: 2640723
          cache_write: 132791
          cost: 0.6592
    - kind: orchestrator
      seconds: 204
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 438
          cache_read: 2678948
          cache_write: 16822
          cost: 0.5922
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1s
    by: flai
    at: 2026-10-28T12:54:38Z
  value: 0.04
  by: planner-S-0531
  at: 2026-10-28T18:24:03Z
forecast:
  duration: 30m
  delivery: 2026-10-28T19:01:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 5, behind S-0415 and S-0502."
  by: flai
  at: 2026-10-28T18:25:36Z
finalized:
  by: alex
  at: 2026-10-28T12:56:13Z
---
# S-0531 A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer

## Goal

This story remediates [I-0209](../../../design/issues/I-0209-a-one-shot-flai-stats-and-the-first-stats-get-of-a-flai-serve-process-still-run-one-git-log-for-claims-drift-about-a-quarter-of-a-cold-answer.md), "A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0209 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0209 is closed with `flai issue close I-0209 --reason` saying what fixed it

## Tasks
- T-2267 storygit.Committed keeps the stories' commits in a file under .flai-cache with its branch-tip key, so a fresh process whose key matches reads no git log
- T-2268 A cold stats.get and flai stats --json are measured again on this repository with the kept file, and server-performance.md records the figures
- T-2269 flai-cli.md, the flai stats guide, and the backup runbook say the stories' commits are kept in a file under .flai-cache between processes
- T-2270 I-0209 is closed with the reason naming the kept file and its measured figures, and the story's criteria are ticked

## Notes

Cost of delay inputs set by flai from I-0209. time_lost_per_cycle 1s: 650ms per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-26T04:42:46Z, 2.3 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0209's one instance and `design/system/server-performance.md` § The stories' commits kept between calls (S-0500): `storygit.Committed` (`flai/internal/storygit/committed.go`) also keeps its answer in a file under the main checkout's `.flai-cache` (`repo.CacheDir()`), with the key it keeps in memory today: the main branch, the wip folder, and `git show-ref --heads`. On a miss in memory it reads the file; when the file's key matches, it answers from the file and runs no `git log`. When it reads the log, it writes the file atomically. A missing, unreadable, or damaged file, or one with another key, is a miss that reads the log, and a failed write never fails the call. This removes the log from every one-shot `flai stats`, `flai hostapi --timing`, and the first `stats.get` of a `flai serve` process while no branch tip moved. The first call after a tip moved still reads the whole log, as today. Reading only the commits added since the kept tips is left out: rebases by `flai stream sync` move story tips where the kept commits do not reach, and S-0500 chose a whole read on change for that reason.

Touches, file by file, no folders:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/storygit/committed.go` | layout | Where `Committed` keeps its answer |
| `flai/internal/storygit/committed_test.go` | layout | The counting-runner tests S-0500 wrote, extended to a fresh process |
| `flai/internal/statsread/statsread_test.go` | co-change (12%) | `TestReadStartsNoLogWhileNothingChanged` lives here; a fresh-process variant through `Read` |
| `design/system/server-performance.md` | design | The issue's figures; the new measurement goes here |
| `design/system/flai-cli.md` | co-change (38%) | § `storygit.Committed` says a one-shot `flai stats` reads the log once |
| `docs/users/flai.md` | co-change (31%) | The touches drift paragraph says the same |
| `docs/operators/runbooks/backup.md` | layout | Lists what a project's `.flai-cache` holds |
| `design/issues/I-0209-…md`, `design/issues/summary.md` | design | Criterion 2, `flai issue close` |

`flai touches suggest` also listed `design/system/flaiover-dashboard.md` (25%), `docs/users/flai-reference.md` (19%), `docs/operators/index.md`, `docs/operators/settings.md`, and `docs/users/flaiover.md` (12% each). None is predicted: the change adds no command, flag, or setting, and nothing on the dashboard changes.

Forecast: flai gives 3m, size 2 with no touches counted. Raised to 30m: S-0500, the same code, tests, measurement, and documents, took 16m in progress, and a kept file adds its format, an atomic write, and tests for a missing, stale, and damaged file. Delivery is flai's 2026-10-29T08:29Z plus the 27m added.

Cost of delay: 0.04 USD a week stands, flai's value from the operator's input of 1s lost per 168h cycle at 150 USD an hour. Nothing in the story gives a reason to change it.
