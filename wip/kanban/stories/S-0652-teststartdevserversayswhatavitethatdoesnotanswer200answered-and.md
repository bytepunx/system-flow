---
id: S-0652
type: story
nature: remediation
title: TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load
status: backlog
owner: alex
created: 2026-11-10T05:44:45Z
updated: 2026-11-10T07:00:52Z
transitions: []
tags: []
topics: [testing]
touches: [flai/internal/components/devserver_test.go, flai/internal/mcpserver/server_test.go, design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md, design/issues/summary.md]
after: [S-0647]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 11m
    by: flai
    at: 2026-11-10T05:44:45Z
  value: 27.5
  by: planner-S-0652
  at: 2026-11-10T07:00:40Z
forecast:
  duration: 45m
  delivery: 2026-11-10T09:09:00Z
  basis: "Its own forecast of 45m; 7th in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0656, S-0241, S-0637, S-0648, S-0649, S-0650 and S-0651."
  by: flai
  at: 2026-11-10T07:00:52Z
---
# S-0652 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Goal

This story remediates [I-0329](../../../design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md), "TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0329 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0329 is closed with `flai issue close I-0329 --reason` saying what fixed it

## Tasks

## Notes

Cost of delay inputs set by flai from I-0329. time_lost_per_cycle 11m: 11m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T05:34:34Z, 0 days before this story; under one cycle counts as one).

### Planning

Two causes, both in the tests, read from the code.

1. `TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered` in `flai/internal/components/devserver_test.go` gives `StartDevServer` a 600 ms `Timeout`. The fake vite is the Go test binary started again by a shell script. Under a load of 115 on 24 cores it is not serving within 600 ms. So `await` never sees the 404, and its deadline error lacks "(it answered 404, and the harness answers 200)", which the test asserts. The product code is right. Proposed fix: the test's deadline outlasts the fake's start under load, and a start delay in the fake longer than the old 600 ms reproduces the failure without load.
2. `TestAnOperatorsReplyBecomesADigestEntry` in `flai/internal/mcpserver/server_test.go` asserts that the server logged nothing. Its server, from `ownersServer`, leaves `Options.Slow` zero, so the timing middleware takes `perf.Slow()`, 500 ms by default. A `thread_reply` that took 637 ms or 783 ms under load logs a slow "request answered" event, and the assertion fails. Proposed fix: `ownersServer` sets `Slow` far above any call's time, so the test still asserts that the digest write logs no warning. `FLAI_SLOW_REQUEST=1ns` reproduces the failure on every run before the fix.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flai/internal/components/devserver_test.go` | design | the first failing test, named in both instances |
| `flai/internal/mcpserver/server_test.go` | design | the second failing test and `ownersServer`, named in both instances |
| `design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md` | design | criterion 2 closes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

`flai touches suggest` from the two test files lists `flai/internal/mcpserver/server.go`, `folder.go`, and `folder_test.go` at 66% and 53%, and `flai/cmd/mcp.go` below them. None is written: both fixes change only the tests, and `Options.Slow` exists in `server.go` already.

After: S-0647, in progress, changes the same dev server tests. It makes them check the pid of the process flai started in place of the pid the fake writes, and adds a start delay to the fake vite. This story reads both: its dev server test ends with the same pid check, `waitGone(t, readPID(t, state, "pid"))`, which loses the pid under load just as I-0311 records, and its reproduction uses that start delay.

Forecast: `flai forecast` gave 5m from the median of 17 small remediation stories at size 2. It is raised to 45m. There are two separate fixes, each run before and after its reproduction. The close-out runs the integration tier, 10 to 20 minutes under load, as I-0329's instances show. The delivery moves 40m past `flai forecast`'s 08:12 for the same reason.

Cost of delay: `flai cod` gives 27.50 USD a week from the 11m `time_lost_per_cycle` flai set from I-0329, at 150 USD an hour. It is written as it stands. I-0329 now has two instances, 15m in all, so flai's inputs undercount it; the inputs are the operator's, and they are left as they are.
