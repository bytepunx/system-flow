---
id: S-0650
type: story
nature: remediation
title: "TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers"
status: done
owner: alex
created: 2026-11-10T05:12:13Z
updated: 2026-11-10T07:43:47Z
transitions:
  - to: ready
    at: 2026-11-10T07:02:50Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T07:19:43Z
    by: agent-S-0650
  - to: review
    at: 2026-11-10T07:30:29Z
    by: agent-S-0650
  - to: done
    at: 2026-11-10T07:43:47Z
    by: alex
tags: []
topics: [testing]
touches: [flai/cmd/components_test.go, design/issues/I-0325-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 655
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 4
      test_runs: 1
      hand_edits: 1
      work: 26
  models:
    - model: claude-opus-5-5
      input: 72
      output: 10931
      cache_read: 7032970
      cache_write: 251812
      cost: 3.64
  strategic:
    - kind: orchestrator
      seconds: 50
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 59
          cache_read: 330611
          cache_write: 6510
          cost: 0.0907
cost_of_delay:
  inputs:
    time_lost_per_cycle: 8m
    by: flai
    at: 2026-11-10T05:12:13Z
  value: 20
  by: planner-S-0650
  at: 2026-11-10T06:58:35Z
forecast:
  duration: 35m
  delivery: 2026-11-10T07:56:00Z
  basis: "Its own forecast of 35m; 1st in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647 and S-0649."
  by: flai
  at: 2026-11-10T07:12:18Z
finalized:
  by: orchestrator
  at: 2026-11-10T07:01:33Z
---
# S-0650 TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers

## Goal

This story remediates [I-0325](../../../design/issues/I-0325-testcomponentsserveendswhenviteexits-fails-under-load-the-fake-vite-exits.md), "TestComponentsServeEndsWhenViteExits fails under load: the fake vite exits before it answers". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0325 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0325 is closed with `flai issue close I-0325 --reason` saying what fixed it

## Tasks
- T-2836 The crash-mode fake vite exits only after it has answered, and a slow-starting fake reproduces the early exit
- T-2837 Close I-0325 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0325. time_lost_per_cycle 8m: 8m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T22:51:04Z, 0.3 days before this story; under one cycle counts as one).

### Planning

Cause, from the code. `TestHelperComponentsVite` in `flai/cmd/components_test.go` is the fake vite: the Go test binary, started again by a script that `exec`s it. In `crash` mode it starts a 300 ms timer, then exits 2, whether or not anyone has asked it anything. `StartDevServer` in `flai/internal/components/devserver.go` polls the harness's path every 150 ms (`pollEvery`). Under load, the fake's listener or flai's poll is slow enough that no poll gets a 200 inside those 300 ms. Then `StartDevServer` returns "exited before it answered", not the "exited: exit status 2" that `TestComponentsServeEndsWhenViteExits` wants. The product code did what it should in both paths. The test's timing is wrong: it means to cover a vite that exits after it answered, but it does not wait for the answer.

Proposed fix: in `crash` mode the fake exits only after it has answered the harness's path once. The handler signals the first answer; the crash goroutine waits for that signal, then a short pause so the response is flushed, then prints "the fake vite crashed" and exits 2. No timer races the poll. A start delay in the fake, slept before it listens and longer than the old 300 ms, reproduces I-0325 on every run without load under the old crash timer. With the fix it passes. The path where vite exits before it answers is covered already, by `TestStartDevServerReturnsTheLastLinesOfAViteThatExitsEarly` in `flai/internal/components/devserver_test.go`.

S-0647 fixes the same kind of race (I-0311) in `flai/internal/components/devserver_test.go`, a different file and a different fake. No `after` on it: this story reads none of its output.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/cmd/components_test.go` | design | the failing test and its fake vite, named in I-0325's instance |
| `design/issues/I-0325-…md` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

Left out from `flai touches suggest` on the test file: `flai/cmd/components.go` (5 of 6 commits) and the rest it lists. They changed with the test when the command changed. This fix changes only the test's fake, and the command behaves correctly.

Forecast: `flai forecast` gave 5m from the median of 17 small remediation stories at size 2. It is raised to 35m. A change under `flai/` runs the integration and smoke tiers in the close-out, 10 to 20 minutes under load, as I-0311's and I-0325's instances show. The fix also needs runs of the components serve tests with the start delay, before and after. It is under S-0647's 40m because only one test file changes. The delivery moves 30m past `flai forecast`'s 08:04 for the same reason.

Cost of delay: `flai cod` gives 20.00 USD a week from the 8m time_lost_per_cycle flai set from I-0325, at 150 USD an hour. It is written as it stands.

No `after`: the story reads no other story's output.
