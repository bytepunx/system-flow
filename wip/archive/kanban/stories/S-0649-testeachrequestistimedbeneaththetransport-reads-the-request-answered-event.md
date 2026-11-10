---
id: S-0649
type: story
nature: remediation
title: TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611
status: done
owner: alex
created: 2026-11-10T05:12:08Z
updated: 2026-11-10T07:43:17Z
transitions:
  - to: ready
    at: 2026-11-10T07:02:50Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-10T07:03:53Z
    by: agent-S-0649
  - to: review
    at: 2026-11-10T07:19:49Z
    by: agent-S-0649
  - to: done
    at: 2026-11-10T07:43:17Z
    by: alex
tags: [flai, test]
topics: [testing]
touches: [flai/internal/channel/channel_test.go, design/issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 967
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 2
      test_runs: 1
      hand_edits: 1
      work: 41
  models:
    - model: claude-opus-5-5
      input: 98
      output: 11773
      cache_read: 8439948
      cache_write: 198253
      cost: 3.5099
  strategic:
    - kind: orchestrator
      seconds: 143
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 18
          output: 174
          cache_read: 647418
          cache_write: 74783
          cost: 0.194
cost_of_delay:
  inputs:
    time_lost_per_cycle: 36m
    by: flai
    at: 2026-11-10T05:12:08Z
  value: 90
  by: planner-S-0649
  at: 2026-11-10T06:58:26Z
forecast:
  duration: 35m
  delivery: 2026-11-10T08:21:00Z
  basis: "Its own forecast of 35m; 4th in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0656, S-0241, S-0637 and S-0648."
  by: flai
  at: 2026-11-10T07:00:52Z
finalized:
  by: orchestrator
  at: 2026-11-10T07:01:31Z
---
# S-0649 TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611

## Goal

This story remediates [I-0324](../../../design/issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md), "TestEachRequestIsTimedBeneathTheTransport reads the request answered event before the channel client has logged it, so it fails under load since S-0611". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0324 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0324 is closed with `flai issue close I-0324 --reason` saying what fixed it

## Tasks
- T-2834 The request answered tests poll the log for the event, and a held-back log line reproduces the race
- T-2835 Close I-0324 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0324. time_lost_per_cycle 36m: 12m per occurrence × 3 occurrences ÷ 1 cycle of 168h (first reported 2026-11-09T22:00:56Z, 0.3 days before this story; under one cycle counts as one).

### Planning

Cause, from the code. S-0611 (2465e6c) moved the channel client's "request answered" event after the answer's write, in the request goroutine in `flai/internal/channel/channel.go`, so that `send_ms` counts the write. `TestEachRequestIsTimedBeneathTheTransport` in `flai/internal/channel/channel_test.go` reads the log once, as soon as the fake dashboard has the reply. Under load the goroutine has not logged yet, so the test finds only the connection lines. The product code does what it should. The test reads too early. `TestTheSendTimeCountsTheWaitForTheWriteLock`, added by the same commit, polls the log with a 3 s deadline already, and does not fail.

Proposed fix: one helper in the test file that polls the log for the "request answered" event with a deadline, used by both tests in place of their own loops. To reproduce the race without load, the first test's log destination holds back the "request answered" line for a while before it writes it: the old one-shot read then fails on every run, and the polled read passes.

Out of scope. The other one-shot log reads in the file (`refusals`, `allowedUnsigned`, `gateRefusals`) read connection events that S-0389 orders before the state they wait on. MCP's `TestEachToolCallIsTimedBeneathTheTransport` logs before it answers, with no `send_ms`.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/channel/channel_test.go` | design | the failing test at line 397, named in all five instances |
| `design/issues/I-0324-…md` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

`flai touches suggest` from the test file lists `flai/internal/channel/channel.go` at 12 of 13 commits. It is left out: the fix changes no product code.

Forecast: `flai forecast` gave 5m from the median of 17 small remediation stories at size 2, with delivery at 07:59. It is raised to 35m. The fix needs repeated runs of the channel tests, with the delay and without. The close-out runs the integration tier, 7 to 14 minutes under load in I-0324's instances. The delivery moves 31m past `flai forecast`'s for the same reason.

Cost of delay: `flai cod` gives 90.00 USD a week from the 36m time_lost_per_cycle flai set from I-0324, at 150 USD an hour. It is written as it stands. The issue counts 5 occurrences now, against the 3 the input was worked out from; the input is the operator's, and the plan thread says so.

No `after`: the story reads no other story's output.
