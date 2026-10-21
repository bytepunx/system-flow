---
id: S-0463
type: story
nature: remediation
title: TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it
status: review
owner: alex
created: 2026-10-21T07:19:49Z
updated: 2026-10-21T09:35:58Z
transitions:
  - to: ready
    at: 2026-10-21T07:33:31Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T07:37:05Z
    by: agent-S-0463
  - to: review
    at: 2026-10-21T08:09:15Z
    by: agent-S-0463
  - to: in-progress
    at: 2026-10-21T08:19:09Z
    by: flai
  - to: review
    at: 2026-10-21T08:57:51Z
    by: agent-S-0463
  - to: in-progress
    at: 2026-10-21T09:23:55Z
    by: flai
  - to: review
    at: 2026-10-21T09:35:58Z
    by: agent-S-0463
tags: [flai, tests, issues]
touches: [flai/internal/host/host_test.go, design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/summary.md, design/issues/I-0175-a-story-s-in-progress-transition-from-flai-story-start-was-lost-when-a-concurrent-write-rewrote-the-item-file-leaving-it-ready-while-its-agent-worked-it.md, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/I-0174-the-merge-queue-s-batch-runs-smoke-with-no-close-out-story-so-scripts-check-sh-runs-an-unscoped-flai-check-strict-and-fails-every-batch-on-findings-outside-its-stories.md, design/issues/I-0179-accepting-a-story-that-adds-a-flaiover-dependency-leaves-the-main-checkout-s-flaiover-node-modules-without-it-so-every-close-out-s-integration-tier-fails-testeverylistedflaiovercomponentcompilesintheharness-until-someone-installs.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4306
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 2
      test_runs: 4
      hand_edits: 2
      work: 69
  models:
    - model: claude-opus-5-5
      input: 160
      output: 35847
      cache_read: 13176350
      cache_write: 209561
      cost: 5.0293
  strategic:
    - kind: planner
      seconds: 62
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 78
          output: 5254
          cache_read: 9386889
          cache_write: 80352
          cost: 2.5393
    - kind: orchestrator
      seconds: 2411
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 200
          output: 3123
          cache_read: 30207875
          cache_write: 86518
          cost: 7.5949
        - model: claude-sonnet-5-5
          input: 8
          output: 69
          cache_read: 68156
          cache_write: 49020
          cost: 0.0959
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-21T07:19:49Z
  value: 25
  by: planner-S-0463
  at: 2026-10-21T07:32:51Z
forecast:
  duration: 25m
  delivery: 2026-10-21T08:02:00Z
  basis: "Its own forecast of 25m; 1st in the pull order with an in-progress limit of 5, behind S-0453 and S-0464."
  by: flai
  at: 2026-10-21T07:33:37Z
finalized:
  by: orchestrator
  at: 2026-10-21T07:33:23Z
---
# S-0463 TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it

## Goal

This story remediates [I-0173](../../../design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md), "TestAChildIsStoppedOnlyOnceItsProcessHasEnded fails under load when the stubborn child dies on SIGTERM before it ignores it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0173 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0173 is closed with `flai issue close I-0173 --reason` saying what fixed it

## Tasks
- T-1977 The stubborn test child says when it ignores SIGTERM, and the host tests wait for that instead of sleeping
- T-1978 Close I-0173 saying the stubborn child now signals that it ignores SIGTERM

## Notes

Cost of delay inputs set by flai from I-0173. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-21T06:07:23Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0173's instance: the `stubborn` child of `TestMain` in `flai/internal/host/host_test.go` writes a file once `signal.Ignore(syscall.SIGTERM)` has returned. `TestAChildIsStoppedOnlyOnceItsProcessHasEnded` and `TestAChildThatWillNotStopIsKilled` wait for that file instead of sleeping 100 ms. A delay knob on the child reproduces the race. The fix is in the tests only: the host's stop logic is right, and the test assumed a readiness it never checked.

Touches, all files, no folder touch:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/host/host_test.go` | design (I-0173) | Holds the stub and both tests that sleep |
| `design/issues/I-0173-…md` | layout | `flai issue close` rewrites the issue file |
| `design/issues/summary.md` | layout | `flai issue close` rewrites the summary, as S-0451's close of I-0169 did |

`flai touches suggest S-0463 flai/internal/host/host_test.go` listed `flai/internal/host/host.go` (4 of 4 commits), and `child.go`, `api.go`, `client.go`, and `flai/cmd/host.go` (2 each). None is taken: the fix changes no production code, and those files co-changed only with the stories that built the host.

Figures:

- Forecast: flai gave 11m (131 s per unit over 8 small remediation stories, times size 5). Raised to 25m. The tasks include a race stress run of the host tests under load, and the close-out's integration tier alone took 5m56s in I-0173's instance. Delivery moved by the same 14m, to 2026-10-22T01:16:00Z.
- Cost of delay value: 25 USD a week, as `flai cod` worked it out from the input flai set (10m lost per 168h cycle at 150 USD an hour). It stands: one occurrence in under a day, and a flake in the integration tier stops whichever story's close-out it hits.
- 2026-10-21T08:19:09Z: moved to in-progress: merge queue batch 20261021-081243 fails with S-0463 at e69059a341aa: smoke failed: smoke: smoke: template render and check
- 2026-10-21T09:23:55Z: moved to in-progress: merge queue batch 20261021-091124 fails with S-0463 at 92c4b9f89cbd: integration failed: flai/internal/serve/agents_test.go:1348 TestASettledReplyStartsTheAgentWaitingOnItOnce: started before the reply: [{At:2026-10-21T09:17:59Z Action:agent Method:serve.agent Project:t Root:/tmp/TestASettledReplyStartsTheAgentWaitingOnItOnce69115204/001 By:flai serve RequestID: Outcome:done Detail:started stub-agent (command) again for S-0001 as builder-S-0001, a message on MS-0001 (pid 1889218); log /tmp/TestASettledReplyStartsTheAgentWaitingOnItOnce69115204/003/serve/agents/t-S-0001-20261021T091759Z.log}]
