---
id: S-0703
type: story
nature: remediation
title: "TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with \"failed to read frame header: EOF\" under host load, on a branch that does not touch flai/internal/serve"
status: backlog
owner: alex
created: 2026-11-10T08:55:40Z
updated: 2026-11-11T01:41:10Z
transitions: []
tags: [flai, tests]
touches: [flai/internal/serve/candidates_test.go, flai/internal/channel/channeltest/dashboard.go, design/issues/I-0267-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md, design/issues/summary.md]
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
      seconds: 56
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 13
          output: 2322
          cache_read: 1435320
          cache_write: 16427
          cost: 0.1019
        - model: claude-opus-5-5
          input: 122
          output: 25698
          cache_read: 9089807
          cache_write: 217495
          cost: 3.982
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 9m
    by: flai
    at: 2026-11-10T08:55:40Z
  value: 22.5
  by: planner-S-0703
  at: 2026-11-10T10:34:08Z
forecast:
  duration: 1h30m
  delivery: 2026-11-11T17:13:00Z
  basis: "Its own forecast of 1h30m; 28th in the pull order with an in-progress limit of 3, behind S-0676, S-0678, S-0674, S-0704, S-0681, S-0241, S-0659, S-0660, S-0661, S-0662, S-0663, S-0683, S-0686, S-0689, S-0690, S-0691, S-0692, S-0693, S-0694, S-0695, S-0696, S-0697, S-0698, S-0699, S-0700, S-0701 and S-0702."
  by: flai
  at: 2026-11-11T01:41:10Z
---
# S-0703 TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve

## Goal

This story remediates [I-0267](../../../design/issues/I-0267-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md), "TestARepositoryIsOfferedImportedAndThenServed fails the close-out's integration tier with "failed to read frame header: EOF" under host load, on a branch that does not touch flai/internal/serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0267 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0267 is closed with `flai issue close I-0267 --reason` saying what fixed it

## Tasks
- T-3042 A serve test reproduces the EOF import.run reads once the client's ping goes unanswered by the test dashboard
- T-3044 The test dashboard reads its connection for its whole life, so it answers the client's pings
- T-3047 The candidate tests pass repeatedly under load, and I-0267 is closed with what fixed it

## Notes

Cost of delay inputs set by flai from I-0267. time_lost_per_cycle 9m: 9m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T02:16:58Z, 3.3 days before this story; under one cycle counts as one).

### Planning

Proposed cause, from I-0267's one instance and the code (planner-S-0703, 2026-11-10):

- The instance ran at 13346d1 (2026-11-07T02:07Z). S-0560's fix, 6920a061c "flai serve drains a candidate it stops offering, so import.run's answer is written" (T-2556), landed at 2026-11-07T10:05Z, after it. That halt is one path to the EOF, and it is closed now.
- A second path remains. The test dashboard, `flai/internal/channel/channeltest/dashboard.go`, reads only inside `Ask` and `Next`; after the proof its handler waits on `<-ctx.Done()`. coder/websocket answers a ping only while a read runs, so the test dashboard answers no ping (`flai/internal/hostapi/components_test.go` says so where it sets `PingEvery: time.Hour`). The serve tests set no `PingEvery`, so the client pings every 4 s with a 4 s timeout (`flai/internal/channel/channel.go`, `defaults` and the ping goroutine in `serveOnce`), cancels the connection when a ping goes unanswered, and `serveOnce`'s `conn.CloseNow` ends the socket with no close frame, which the dashboard reads as "failed to read frame header: EOF". Under a load average of 56, the gap between a connection's proof and the test's next read can pass a ping's tick and timeout.
- Fix: the test dashboard reads its connection for its whole life, so it answers pings, and hands `Ask` and `Next` the messages it reads. A test with a short `PingEvery` and a pause before `import.run` reproduces the EOF first.

Touches:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/candidates_test.go` | design (I-0267) | the failing test; the reproduction goes beside it |
| `flai/internal/channel/channeltest/dashboard.go` | layout | the test dashboard that answers no ping |
| `design/issues/I-0267-testarepositoryisofferedimportedandthenserved-fails-the-close-out-s-integration.md` | declared by the criteria | closed with `flai issue close` |
| `design/issues/summary.md` | co-change | regenerated when the issue closes |

`flai touches suggest S-0703` gave nothing: the story declared no touches. No folder touch.

Figures:

- Forecast: flai's 6m comes from size 2 (2 criteria, 0 touches) at 160 s per unit. Raised to 1h30m. Three tasks, two Go files, repeated runs of the serve package under load with `-race`, and a close-out whose integration tier ran 8 to 20 minutes on recent stories. The delivery is flai's queue estimate, moved by the longer duration.
- Cost of delay: 22.50 USD a week, flai's value from the operator's input of 9m lost per 168h cycle at 150 USD an hour. It stands. The one instance predates S-0560's fix, so the rate may be lower, but the input is the operator's.
- No `after`: the story reads nothing another open story writes.
