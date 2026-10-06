---
id: S-0292
type: story
nature: remediation
title: flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run
status: done
owner: alex
created: 2026-10-06T11:14:34Z
updated: 2026-10-06T11:44:49Z
transitions:
  - to: ready
    at: 2026-10-06T11:31:51Z
    by: alex
  - to: in-progress
    at: 2026-10-06T11:33:08Z
    by: agent-S-0292
  - to: review
    at: 2026-10-06T11:43:32Z
    by: agent-S-0292
  - to: done
    at: 2026-10-06T11:44:49Z
    by: alex
tags: []
touches: [design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/I-0090-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md, design/issues/summary.md, flai/internal/serve/orchestrate_test.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 658
  models:
    - model: claude-opus-5-5
      input: 94
      output: 20261
      cache_read: 4905328
      cache_write: 137136
      cost: 2.4837
    - model: claude-sonnet-5-5
      input: 12
      output: 2684
      cache_read: 90536
      cache_write: 48628
      cost: 0.1665
  strategic:
    - kind: planner
      seconds: 178
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 122
          output: 19954
          cache_read: 4097914
          cache_write: 148397
          cost: 2.6997
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-06T11:14:34Z
  value: 25
  by: planner-S-0292
  at: 2026-10-06T11:34:01Z
forecast:
  duration: 1h
  delivery: 2026-10-06T12:35:00Z
  basis: "flai forecast's 16m (median 118 s per unit over 11 done medium remediation stories, size 8) raised to 1h: a race that shows only under the full parallel run needs repeated -race stress runs to reproduce and to prove fixed, each minutes long; delivery is 1h from its start at 11:33Z"
  by: planner-S-0292
  at: 2026-10-06T11:34:01Z
finalized:
  by: alex
  at: 2026-10-06T11:31:48Z
---
# S-0292 flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run

## Goal

This story remediates [I-0090](../../../design/issues/I-0090-flai-serve-s-orchestrator-test-that-waits-for-the-stopped-run-s-activity-fails-under-the-close-out-s-full-parallel-test-run.md), "flai serve's orchestrator test that waits for the stopped run's activity fails under the close-out's full parallel test run". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0090 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0090 is closed with `flai issue close I-0090 --reason` saying what fixed it

## Tasks
- T-1016 The orchestrator stub prints its stream before it says what it was given, so a stop cannot come before its usage is in the run's log
- T-1017 Reproduce the lost orchestrator activity entry and name its cause
- T-1018 Fix the cause so the stopped orchestrator's activity is always logged, and the test passes under load
- T-1019 Close I-0090 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0090. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-06T11:05:36Z, 0 days before this story; under one cycle counts as one).

### Planning

Touches, by where each came from:

- `flai/internal/serve/orchestrate_test.go`: design, the failing test `TestTheOrchestratorIsStoppedWhenTheActionIsTurnedOff` that I-0090 names. The story declared no touches; `flai touches suggest` started from this path.
- `flai/internal/serve/orchestrate.go`: layout. `orchestrateEnded` logs the stopped run's activity, and `stop` settles the run when no launcher waits for it.
- `flai/internal/serve/activity.go`: layout. `logRunEnd` returns without an error, and logs nothing, when it finds no newest run, a span outside the run, or no usage. That is the likeliest place for an entry that "never happened" without a warning.
- `flai/internal/serve/agents.go`: co-change (9% with the files above) and layout. `await` reaps the process and calls `orchestrateEnded` from a goroutine, and `settleOrphans` settles the handed-over case.
- `flai/internal/serve/agents_test.go`: layout. It holds `waitFor` and its fixed 5 s deadline, which the test failed at (5.07 s).
- `flai/internal/serve/plan.go`: layout. `planEnded` logs a planner run's activity through the same `logRunEnd`, so it may share the cause and the fix.
- `design/issues/I-0090-…md` and `design/issues/summary.md`: design, written by `flai issue close` for the second criterion.

The rest of `flai touches suggest` (flai-cli.md, docs/users/flai.md, older issues, flaiover routes) is general co-change with `flai/internal/serve` and not this fix's, so I left it out. A fix that changes what flai serve does, not only the test, would add `design/system/strategic-agents.md`.

Figures:

- Forecast: flai forecast gave 16m (median 118 s per unit of size over 11 done remediation stories on claude-opus-5-5 in the medium band, times size 8). I raised it to 1h. The failure shows only under the full parallel run, so reproducing it and proving it fixed takes repeated `-race` stress runs, each several minutes long; most of the median's stories had a known cause. The story was pulled at 11:33Z while I planned it, so the delivery is 1h from then, 12:35Z, not flai's play-out from 9th in the pull order.
- Cost of delay: 25 USD a week, from flai cod on the inputs flai set from I-0090 (10m lost per 168h cycle at 150 USD an hour). It stands, though a flake that stops a close-out costs every story whose run it hits; I-0090's count will raise the inputs if it recurs.
