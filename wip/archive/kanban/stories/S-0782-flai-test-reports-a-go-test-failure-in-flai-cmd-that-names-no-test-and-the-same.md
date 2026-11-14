---
id: S-0782
type: story
nature: remediation
title: flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again
status: done
owner: alex
created: 2026-11-13T22:46:52Z
updated: 2026-11-14T07:04:29Z
transitions:
  - to: ready
    at: 2026-11-13T22:54:43Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-13T22:54:53Z
    by: agent-S-0782
  - to: review
    at: 2026-11-13T23:19:24Z
    by: agent-S-0782
  - to: done
    at: 2026-11-14T07:04:29Z
    by: alex
tags: [flai, testing]
touches: [flai/internal/verify/gotest.go, flai/internal/verify/parse_test.go, flai/internal/verify/testdata/go-test-crash.json, design/system/project-manifest.md, design/issues/I-0365-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md, design/issues/summary.md, design/issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1483
  priced_by: harness
  turns:
    - day: 2026-11-13
      ceremony: 4
      work: 30
  models:
    - model: claude-opus-5-5
      input: 108
      output: 29994
      cache_read: 7376213
      cache_write: 361330
      cost: 4.5426
  strategic:
    - kind: planner
      seconds: 210
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 38
          output: 125
          cache_read: 1530154
          cache_write: 75944
          cost: 0.1094
        - model: claude-opus-5-5
          input: 62
          output: 15780
          cache_read: 4128627
          cache_write: 186398
          cost: 2.541
    - kind: orchestrator
      seconds: 114
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 41
          output: 460
          cache_read: 1443752
          cache_write: 196424
          cost: 0.4843
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-11-13T22:46:52Z
  value: 15
  by: planner-S-0782
  at: 2026-11-13T22:52:52Z
forecast:
  duration: 20m
  delivery: 2026-11-14T05:28:00Z
  basis: flai forecast's 6m (median 152 s per size unit over 24 small remediation stories, size 2) raised to the 15m S-0540 took for the same parser's timeout case, plus a second fixture and the issue close.
  by: planner-S-0782
  at: 2026-11-13T22:52:52Z
finalized:
  by: orchestrator
  at: 2026-11-13T22:54:23Z
---
# S-0782 flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again

## Goal

This story remediates [I-0365](../../../design/issues/I-0365-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md), "flai test reports a go-test failure in flai/cmd that names no test, and the same paths pass when run again". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0365 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0365 is closed with `flai issue close I-0365 --reason` saying what fixed it

## Tasks
- T-3344 go-test-json names the test a crashed test binary was running, with its crash and exit, when a package fails with no test failing
- T-3346 Close I-0365 with flai issue close, its reason naming what go-test-json now names for a crashed test binary

## Notes

Cost of delay inputs set by flai from I-0365. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T20:17:49Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed solution. The `go-test-json` format (`flai/internal/verify/gotest.go`, `parseGoTest`) gives a package that fails with no failed test and no timeout panic one finding, named for the package, whose message is `packageMessage` of the package's own output lines only, those `go test -json` sends with no `Test`. When the test binary dies mid-run, the crash is written as output of the test that was running: a panic in a goroutine, a `fatal error:`, a test's `os.Exit`. That test gets no `fail` event, so the package's own lines hold only `FAIL<tab><package>`, which `packageMessage` drops, and the finding names no test, line, or message: the shape I-0365 reports. The S-0540 fix for I-0227 covers the timeout panic alone. The fix: for such a package, name the test that started and never ended, place it by the crash's first frame under the root, give the crash line (`panic:`, `fatal error:`, the `os.Exit` panic) and the binary's exit (`exit status n`, `signal: killed`), and fall back to the last lines of the package's whole output, its tests' included, never an empty message. That makes the next occurrence readable. The flake behind the instance, a crash under load that passed when run again, cannot be named from the instance, which kept no output. Its own remediation is a story once a readable finding names it.

Touches:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/verify/gotest.go` | design, layout | `parseGoTest` and its package-failure branch, as S-0540 (commit 3b0277d8a) changed |
| `flai/internal/verify/parse_test.go` | layout, co-change | the go-test-json tests sit here, S-0540 added its timeout tests here |
| `flai/internal/verify/testdata/go-test-crash.json` | layout, co-change | a new fixture recorded from a real go test run, as S-0540 added `go-test-timeout.json` |
| `design/system/project-manifest.md` | design, co-change | its `format` paragraph says what `go-test-json` names for a package failure, and S-0540 changed it |
| `design/issues/I-0365-flai-test-reports-a-go-test-failure-in-flai-cmd-that-names-no-test-and-the-same.md` | design | criterion 2 closes it with `flai issue close` |

`flai touches suggest S-0782` listed nothing: the story declared no touches. No folder touch.

No `after`: the story reads no other open story's output.

Figures:

- Forecast 20m, raised from flai's 6m. flai's figure is the median of small remediation stories times size 2. S-0540 fixed the same parser for the timeout panic in 908 s of agent time, with one fixture and one doc line. This story records a second fixture and closes the issue in a task of its own. Delivery moves by the same 14 minutes, to 2026-11-14T05:28:00Z.
- Cost of delay 15 USD a week, as `flai cod S-0782` works it out from the input flai set from I-0365: 6m lost per 168h cycle at 150 USD an hour. It stands: one occurrence so far, and each costs a rerun of the go-test tier.
