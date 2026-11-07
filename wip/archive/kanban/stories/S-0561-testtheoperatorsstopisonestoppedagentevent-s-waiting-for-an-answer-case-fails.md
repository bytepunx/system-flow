---
id: S-0561
type: story
nature: remediation
title: "TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped"
status: done
owner: alex
created: 2026-11-02T03:56:56Z
updated: 2026-11-07T08:36:58Z
transitions:
  - to: ready
    at: 2026-11-07T08:22:45Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T08:22:54Z
    by: agent-S-0561
  - to: review
    at: 2026-11-07T08:34:35Z
    by: agent-S-0561
  - to: done
    at: 2026-11-07T08:36:58Z
    by: orchestrator
tags: []
touches: [flai/internal/serve/stop_test.go, design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 710
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 4
      test_runs: 2
      hand_edits: 1
      work: 71
  models:
    - model: claude-opus-5-5
      input: 164
      output: 13276
      cache_read: 19476162
      cache_write: 261737
      cost: 6.2553
  strategic:
    - kind: planner
      seconds: 206
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 84
          output: 20321
          cache_read: 6017857
          cache_write: 338487
          cost: 3.0536
    - kind: orchestrator
      seconds: 188
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 273
          cache_read: 1282651
          cache_write: 176342
          cost: 0.3614
        - model: claude-sonnet-5-5
          input: 14
          output: 89
          cache_read: 129626
          cache_write: 61011
          cost: 0.1374
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-02T03:56:56Z
  value: 12.5
  by: planner-S-0561
  at: 2026-11-07T08:20:07Z
forecast:
  duration: 20m
  delivery: 2026-11-07T13:15:00Z
  basis: "flai forecast's 11m (131 s per unit of size over 14 small remediation stories, size 5), raised to S-0488's 1167 s, the same shape: one flaky test file, its issue, and the summary, closed out with the integration and smoke tiers a flai path selects; delivery is flai's 13:06 plus the 9m added."
  by: planner-S-0561
  at: 2026-11-07T08:20:07Z
finalized:
  by: orchestrator
  at: 2026-11-07T08:22:27Z
---
# S-0561 TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped

## Goal

This story remediates [I-0248](../../../design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md), "TestTheOperatorsStopIsOneStoppedAgentEvent's waiting-for-an-answer case fails under a loaded close-out: the ended event is missing before stopped". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0248 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0248 is closed with `flai issue close I-0248 --reason` saying what fixed it

## Tasks
- T-2550 The waiting-for-an-answer case waits for the ended event before it stops the agent
- T-2552 Close I-0248 saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0248. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T02:50:06Z, 0 days before this story; under one cycle counts as one).

### Planning

Cause, read from the code: `endRun` in `flai/internal/serve/agents.go` records the run's outcome with `l.ended(run)`, then calls `appendAgentEnd`. The case at `flai/internal/serve/stop_test.go` waits for the outcome and stops the agent at once, so on a loaded host the stopped event can be appended first. Proposed fix: the case also waits for the ended event before it stops. No code outside the test changes, since the event log is news, not state (ADR-0211 §4), and a stop landing in that gap costs a reader only the order of two events it still gets.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/serve/stop_test.go` | design | I-0248 names it at line 321, and the fix is there |
| `design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md` | design | The second criterion closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

The story declared no touches. `flai touches suggest` from the test file listed `flai/internal/serve/agents.go`, `stop.go`, and `agents_test.go`. None is taken, since no code outside the test changes and the temporary sleep is not committed. Its docs and `proc_*` entries are co-change noise.

Forecast: 20m, adjusted up from `flai forecast`'s 11m. S-0488 took 1167 s with the same shape, one flaky test file, its issue, and the summary, and a path under `flai/` selects the integration and smoke tiers at close-out. Delivery is flai's 2026-11-07T13:06:00Z plus the 9m added.

Cost of delay: 12.50 USD a week stands, as `flai cod` works it out from the inputs: 5m per 168h cycle at 150 USD an hour. Those inputs date from the first occurrence; I-0248 now counts 3 within one cycle. They are the operator's and are unchanged.

### Accepted by the orchestrator

- Verified: 0178c4c23e9f756e3c0a6c120ff2f11b7e3953d0
- At: 2026-11-07T08:36:58Z

Verdict: pass. Both criteria are met at commit 0178c4c23e9f756e3c0a6c120ff2f11b7e3953d0; the close-out passed every step there.
- 1: flai/internal/serve/stop_test.go
- 2: design/issues/I-0248-testtheoperatorsstopisonestoppedagentevent-s-waiting-for-an-answer-case-fails.md, design/issues/summary.md
