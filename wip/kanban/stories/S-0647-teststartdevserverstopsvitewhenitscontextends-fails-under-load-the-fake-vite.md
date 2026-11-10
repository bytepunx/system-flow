---
id: S-0647
type: story
nature: remediation
title: "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid"
status: ready
owner: alex
created: 2026-11-10T05:12:04Z
updated: 2026-11-10T06:27:20Z
transitions:
  - to: ready
    at: 2026-11-10T06:27:20Z
    by: orchestrator
tags: []
topics: [testing]
touches: [flai/internal/components/devserver_test.go, flai/internal/components/devserver.go, design/issues/I-0311-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md, design/issues/summary.md]
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
      seconds: 239
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 225
          cache_read: 1159864
          cache_write: 131948
          cost: 0.3462
cost_of_delay:
  inputs:
    time_lost_per_cycle: 51m
    by: flai
    at: 2026-11-10T05:12:04Z
  value: 127.5
  by: planner-S-0647
  at: 2026-11-10T06:23:46Z
forecast:
  duration: 40m
  delivery: 2026-11-10T08:02:00Z
  basis: flai forecast's 5m (median of 17 small remediation stories at size 2) raised to 40m for the integration tier's close-out, 10 to 20 minutes under load, and repeated runs of the dev server tests with a start delay
  by: planner-S-0647
  at: 2026-11-10T06:23:46Z
finalized:
  by: orchestrator
  at: 2026-11-10T06:26:11Z
---
# S-0647 TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid

## Goal

This story remediates [I-0311](../../../design/issues/I-0311-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md), "TestStartDevServerStopsViteWhenItsContextEnds fails under load: the fake vite wrote no pid". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0311 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0311 is closed with `flai issue close I-0311 --reason` saying what fixed it

## Tasks
- T-2825 The dev server tests check the pid flai started, and a slow-starting fake vite reproduces the missing pid
- T-2826 Close I-0311 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0311. time_lost_per_cycle 51m: 17m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-08T01:43:19Z, 2.1 days before this story; under one cycle counts as one).

### Planning

Cause, from the code. The fake vite in `flai/internal/components/devserver_test.go` is the Go test binary started again by a shell script that `exec`s it. It writes its pid to the state folder only once that binary is running. `TestStartDevServerStopsViteWhenItsContextEnds` ends its context after 300 ms. Under a load of 100 or more on 24 cores, the test binary is not running by then, so `StartDevServer` stops a process that never wrote its pid, and `readPID` fails. The product code did what it should: it stopped vite. The test lost the pid it checks. `TestStartDevServerStopsAViteThatDoesNotAnswerInTime`, at 500 ms, has the same window.

Proposed fix: the tests take the pid of the process flai started, not the pid the fake writes. The script's `exec` keeps the pid, so `cmd.Process.Pid` is the fake's pid. A package seam in `devserver.go`, beside `pickPort`, hands that pid to the tests. A start delay in the fake, longer than the test's context, reproduces the failure on every run without load.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/components/devserver_test.go` | design | the failing test and its helpers, named in all three instances |
| `flai/internal/components/devserver.go` | co-change | `flai touches suggest` from the test: 2 of 2 commits; the seam that records the started pid goes beside `pickPort` |
| `design/issues/I-0311-…md` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Forecast: `flai forecast` gave 5m from the median of 17 small remediation stories at size 2. It is raised to 40m. The close-out runs the integration tier, 10 to 20 minutes under load as I-0311's instances show. The fix needs runs of the dev server tests with the start delay, before and after. The delivery moves 35m past `flai forecast`'s 07:27 for the same reason.

Cost of delay: `flai cod` gives 127.50 USD a week from the 51m time_lost_per_cycle flai set from I-0311, at 150 USD an hour. It is written as it stands.

No `after`: the story reads no other story's output.
