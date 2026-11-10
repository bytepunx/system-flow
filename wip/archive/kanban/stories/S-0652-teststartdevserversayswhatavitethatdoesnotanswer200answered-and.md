---
id: S-0652
type: story
nature: remediation
title: TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load
status: done
owner: alex
created: 2026-11-10T05:44:45Z
updated: 2026-11-10T08:32:38Z
transitions:
  - to: ready
    at: 2026-11-10T07:44:19Z
    by: alex
  - to: in-progress
    at: 2026-11-10T08:03:55Z
    by: agent-S-0652
  - to: review
    at: 2026-11-10T08:29:33Z
    by: agent-S-0652
  - to: done
    at: 2026-11-10T08:32:38Z
    by: orchestrator
tags: []
topics: [testing]
touches: [flai/internal/components/devserver_test.go, flai/internal/mcpserver/server_test.go, design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md, design/issues/summary.md, design/issues/I-0331-testadevserverthatexitsonitsownisforgottenwithitserror-fails-the-integration.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md, flai/internal/serve/components_test.go]
after: [S-0647]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1550
  priced_by: harness
  turns:
    - day: 2026-11-10
      test_runs: 2
      hand_edits: 3
      work: 140
  models:
    - model: claude-opus-5-5
      input: 296
      output: 30966
      cache_read: 41456760
      cache_write: 311814
      cost: 11.4064
  strategic:
    - kind: orchestrator
      seconds: 121
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 182
          cache_read: 1123796
          cache_write: 72667
          cost: 0.3228
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
  delivery: 2026-11-10T08:58:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 5, behind S-0636, S-0637, S-0638, S-0651 and S-0654."
  by: flai
  at: 2026-11-10T07:57:16Z
finalized:
  by: orchestrator
  at: 2026-11-10T07:08:12Z
---
# S-0652 TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load

## Goal

This story remediates [I-0329](../../../design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md), "TestStartDevServerSaysWhatAViteThatDoesNotAnswer200Answered and TestAnOperatorsReplyBecomesADigestEntry fail the integration tier under host load". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0329 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0329 is closed with `flai issue close I-0329 --reason` saying what fixed it

## Tasks
- T-2841 The dev server's 404 test allows the fake vite time to start, and a start delay reproduces the missing 404
- T-2843 The owners' MCP test server sets a slow threshold no call reaches, and FLAI_SLOW_REQUEST=1ns reproduces the slow log
- T-2845 Close I-0329 with what fixed it
- T-2859 The crash-later fake vites crash only after they have answered, and I-0331 is closed

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

### Accepted by the orchestrator

- Verified: a7a8b0e32fa6ca939ae851410d11848f9afa2a47
- At: 2026-11-10T08:32:38Z

Verdict: pass at a7a8b0e32fa6ca939ae851410d11848f9afa2a47; flai verify passed every step, integration and smoke included.
- 1: flai/internal/components/devserver_test.go, flai/internal/mcpserver/server_test.go, flai/internal/serve/components_test.go
- 2: design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md, design/issues/summary.md
