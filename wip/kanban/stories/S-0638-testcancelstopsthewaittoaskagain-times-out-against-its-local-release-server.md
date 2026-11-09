---
id: S-0638
type: story
nature: remediation
title: TestCancelStopsTheWaitToAskAgain times out against its local release server under load and stops a close-out's integration tier
status: backlog
owner: alex
created: 2026-11-09T09:08:54Z
updated: 2026-11-09T22:09:59Z
transitions: []
tags: [flai, testing, flaky-test]
topics: [testing]
touches: [flai/internal/selfupgrade/selfupgrade.go, flai/internal/selfupgrade/selfupgrade_test.go, design/issues/I-0317-testcancelstopsthewaittoaskagain-times-out-against-its-local-release-server.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-09T09:08:54Z
  value: 37.5
  by: planner-S-0638
  at: 2026-11-09T22:07:40Z
forecast:
  duration: 30m
  delivery: 2026-11-10T00:35:00Z
  basis: "Its own forecast of 30m; 9th in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0634, S-0633, S-0625, S-0630, S-0628, S-0631, S-0241, S-0636 and S-0637."
  by: flai
  at: 2026-11-09T22:09:59Z
---
# S-0638 TestCancelStopsTheWaitToAskAgain times out against its local release server under load and stops a close-out's integration tier

## Goal

This story remediates [I-0317](../../../design/issues/I-0317-testcancelstopsthewaittoaskagain-times-out-against-its-local-release-server.md), "TestCancelStopsTheWaitToAskAgain times out against its local release server under load and stops a close-out's integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0317 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0317 is closed with `flai issue close I-0317 --reason` saying what fixed it

## Tasks
- T-2806 End TestCancelStopsTheWaitToAskAgain's wait by an event, not a 100 ms deadline, and reproduce a slow first answer
- T-2808 Close I-0317 naming the event-driven cancel that fixed it

## Notes

Cost of delay inputs set by flai from I-0317. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T08:39:50Z, 0 days before this story; under one cycle counts as one).

### Planning

The cause, read from the test and `selfupgrade.go`: `TestCancelStopsTheWaitToAskAgain` gives `List` a context with a 100 ms wall-clock deadline and sets `retryPause` to an hour. It expects the first attempt to fail on the dropped answer (`unexpected EOF`) and the deadline to end the wait before the second. On a host with a load average of 106 on 24 cores, the first attempt itself took longer than 100 ms, so the deadline cut it before the drop: the error was `context deadline exceeded` alone. The deadline races the first request.

Proposed fix: end the wait by an event, not by the clock. `get` in `selfupgrade.go` gains a test seam called when the pause between attempts begins, nil in production. The test cancels its context from that seam, so the cancel always comes after the first attempt failed and before a second, whatever the load. A second test reproduces the load: its front handler holds the first answer longer than the old 100 ms deadline before it drops it, and still gets the dropped answer and the cancel, with one request.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/selfupgrade/selfupgrade_test.go` | design (the issue's instance) | the test that times out |
| `flai/internal/selfupgrade/selfupgrade.go` | co-change (7 of 7 commits), layout | the retry loop gains the seam |
| `design/issues/I-0317-testcancelstopsthewaittoaskagain-times-out-against-its-local-release-server.md` | design (criterion 2) | `flai issue close` writes it |
| `design/issues/summary.md` | design (criterion 2) | `flai issue close` regenerates it |

Co-change also lists `flai/cmd/selfupgrade.go` and `signature.go` at 29%: left out, since neither the retry loop nor the test is in them. No folder touch. No `after`: the story reads no other story's output.

Figures:

- Forecast: flai forecast gave 15m (144 s per unit of size, size 6). Set to 30m, since the close-out runs the integration tier for any change under `flai/`, and that tier alone took 9m21s in I-0317's instance. Delivery moved by the same 15m.
- Cost of delay: 37.50 USD a week, as `flai cod` works it from the 15m input at 150 USD an hour, one cycle a week. It stands: the instance cost about 12 minutes of close-out, which the 15m input covers.
