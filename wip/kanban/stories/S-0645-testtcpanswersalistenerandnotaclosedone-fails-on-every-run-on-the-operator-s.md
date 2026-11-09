---
id: S-0645
type: story
nature: remediation
title: TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there
status: ready
owner: alex
created: 2026-11-09T21:50:55Z
updated: 2026-11-09T22:20:09Z
transitions:
  - to: ready
    at: 2026-11-09T22:19:58Z
    by: orchestrator
tags: [flai, testing]
topics: [testing, go]
touches: [flai/cmd/dashboard_porthold.go, flai/cmd/dashboard_porthold_test.go, design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md, design/issues/summary.md]
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
      seconds: 327
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 72
          output: 274
          cache_read: 3439488
          cache_write: 155745
          cost: 0.2411
        - model: claude-opus-5-5
          input: 343
          output: 37243
          cache_read: 20281318
          cache_write: 1006042
          cost: 8.1255
cost_of_delay:
  inputs:
    time_lost_per_cycle: 12m
    by: flai
    at: 2026-11-09T21:50:55Z
  value: 30
  by: system-flow
  at: 2026-11-09T22:11:02Z
forecast:
  duration: 20m
  delivery: 2026-11-09T22:57:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0633 and S-0634."
  by: flai
  at: 2026-11-09T22:20:09Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:19:35Z
---
# S-0645 TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there

## Goal

This story remediates [I-0321](../../../design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md), "TestTCPAnswersAListenerAndNotAClosedOne fails on every run on the operator's WSL2 host, on main as well, because a closed 127.0.0.1 port still accepts a connection there". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0321 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0321 is closed with `flai issue close I-0321 --reason` saying what fixed it

## Tasks
- T-2816 Make TestTCPAnswersAListenerAndNotAClosedOne's closed-port case independent of what answers on the host's loopback
- T-2817 Close I-0321 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0321. time_lost_per_cycle 12m: 6m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-09T17:13:14Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed solution. Fix the test, not the probe. D-0444 records that this host runs WSL2 with mirrored networking. I-0321's second instance says the failure comes and goes. On 2026-11-09 at about 22:08 UTC, `go test -count=10 -run TestTCPAnswersAListenerAndNotAClosedOne ./cmd` passed 10 of 10 in the main checkout. So a loopback port that was free in Linux is answered only some of the time, likely by the mirrored relay to a Windows listener. `tcpAnswers` is a thin wrapper over `net.DialTimeout`, so a real closed-port dial tests the host, not flai. T-2816 gives the probe a dial seam. The closed cases then go through a dialer that refuses, and only a listener the test holds open is dialled for real. Production keeps `net.DialTimeout`. D-0444's advice, to dial a listener opened on `127.0.0.1:0` and closed, still holds for tests that need only a failed request. It does not hold for a test that asserts a bare TCP connect is refused. The probe's production use, the dashboard's port-held state (ADR-0267), may see the same relay. That is outside this story: see the plan's thread.

Touches, all files:

| Touch | Source |
|-------|--------|
| `flai/cmd/dashboard_porthold.go` | layout: `tcpAnswers` and `publishAnswers` live here; T-2816 adds the dial seam |
| `flai/cmd/dashboard_porthold_test.go` | layout: the failing test |
| `design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md` | criterion 2: `flai issue close` writes it |
| `design/issues/summary.md` | criterion 2: flai regenerates it on the close |

`flai touches suggest` from the two code files found nothing: 1 of 3323 commits changed them, S-0615's c97bd655e. No other test in `flai/` asserts that a real port refuses a TCP connect. The `closedURL` helpers in `flai/cmd/dashboard_test.go` and `flai/internal/serve/mcp_test.go` need only a failed HTTP request, which a relay that accepts and then closes still gives. No `after`: the story reads nothing another open story writes.

Forecast 20m. flai forecast gives 15m: a median of 144 s per unit of size over 17 done small remediation stories on claude-opus-5-5, times size 6 (2 criteria, 4 touches). I added 5m because T-2816 first confirms an intermittent cause with repeated test runs on the host. Delivery is moved the same 5m, to 2026-11-10T01:34:00Z.

Cost of delay 30.00 USD a week, as flai cod gives it from the operator's input: 12m of time lost per 168h cycle at 150 USD an hour. It stands. The issue's two instances are one story's. Every story whose tests select `./cmd` can hit the failure, so the real loss may be higher. Raising it is the operator's call on the input.

Tasks, in two layers:

1. T-2816: confirm the cause on the host, then make the test's closed-port case independent of the host's loopback.
2. T-2817, after T-2816: close I-0321 with a reason that quotes T-2816's finding and fix.
