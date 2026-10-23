---
id: S-0488
type: story
nature: remediation
title: TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve
status: done
owner: alex
created: 2026-10-23T17:53:08Z
updated: 2026-10-23T20:04:40Z
transitions:
  - to: ready
    at: 2026-10-23T19:34:00Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-23T19:34:06Z
    by: agent-S-0488
  - to: review
    at: 2026-10-23T19:53:23Z
    by: agent-S-0488
  - to: done
    at: 2026-10-23T20:04:40Z
    by: orchestrator
tags: [cli]
topics: [cli]
touches: [flai/cmd/serve_actions_test.go, design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1167
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 1
      hand_edits: 1
      work: 30
  models:
    - model: claude-opus-5-5
      input: 66
      output: 9523
      cache_read: 3858627
      cache_write: 126891
      cost: 1.9776
  strategic:
    - kind: planner
      seconds: 163
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 694
          output: 14205
          cache_read: 3400632
          cache_write: 122988
          cost: 4.0271
    - kind: orchestrator
      seconds: 420
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 75
          output: 1168
          cache_read: 20357831
          cache_write: 31610
          cost: 4.4844
        - model: claude-sonnet-5-5
          input: 4
          output: 31
          cache_read: 12507
          cache_write: 31783
          cost: 0.0305
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-10-23T17:53:08Z
  value: 17.5
  by: planner-S-0488
  at: 2026-10-23T18:10:04Z
forecast:
  duration: 15m
  delivery: 2026-10-24T13:19:00Z
  basis: "Its own forecast of 15m; 28th in the pull order with an in-progress limit of 5, behind S-0489, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0430, S-0431, S-0432, S-0433, S-0434, S-0435 and S-0486."
  by: flai
  at: 2026-10-23T19:11:14Z
finalized:
  by: orchestrator
  at: 2026-10-23T18:11:28Z
---
# S-0488 TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve

## Goal

This story remediates [I-0197](../../../design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md), "TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0197 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0197 is closed with `flai issue close I-0197 --reason` saying what fixed it

## Tasks
- T-2087 TestServeSettingsAreShownToTheDashboard refuses the registered entry's URL, not the bare port 4242, and a case with an import root named with 4242 proves the temp-directory path no longer fails it
- T-2088 I-0197 is closed with the reason naming the test's check narrowed to the registered entry's URL

## Notes

Cost of delay inputs set by flai from I-0197. time_lost_per_cycle 7m: 7m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-23T05:12:51Z, 0.5 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0488 on 2026-10-23. The cause, confirmed in the test rather than guessed from the cut-off finding: `TestServeSettingsAreShownToTheDashboard` (`flai/cmd/serve_actions_test.go`, line 980) refuses any `4242` anywhere in the settings JSON, to prove the registered entry's URL `http://127.0.0.1:4242` is not shown. The same JSON carries two `t.TempDir()` paths, the project root from `tempProject` and the import root the test adds, and Go names a temp directory after the test with a random number, so a run whose number contains `4242` fails with "no token in the settings" and the next run passes. The solution proposed: the test refuses the registered entry's URL as it was registered, `127.0.0.1:4242`, not the bare port, keeping the `"token"` and `/p/harbour/k` checks, and gains a case whose import root is a folder named with `4242` under the temp directory, which fails under the old check and passes under the new one. `hostSettings` is right to show those paths, so no code outside the test changes.

Touches:

- `flai/cmd/serve_actions_test.go`: declared by the issue's instance and the layout, the test that fails (T1).
- `design/issues/I-0197-…md`, `design/issues/summary.md`: the second criterion, `flai issue close` rewrites both (T2).
- `flai touches suggest` from the test file lists `flai/cmd/serve_actions.go`, `flai/internal/hostapi/writes.go`, and the serve, config, and settings sources it co-changes with: left out, since the settings are right and only the test's check changes.

Forecast. `flai forecast` gave 5m: 131 s per unit of size over 10 done remediation stories on claude-opus-5-5 in the small band, size 2 (2 criteria, 0 touches), delivery 2026-10-24T12:02Z as 28th in the pull order. Raised to 15m: the change is a few lines, but the agent runs the go-test tier after the test task, closes two tasks one after the other, and runs the close-out, which no 5m story has fitted. Delivery moved by the 10m added, to 2026-10-24T12:12Z.

Cost of delay. `flai cod` gives 17.50 USD a week from the inputs flai set (7m per 168h cycle at 150 USD an hour), and the value stands at that: the inputs are not mine to change. They understate it a little: each occurrence costs a rerun of the go-test tier, about 6 minutes, on any branch whose close-out draws an unlucky temp directory.

### Accepted by the orchestrator

- Verified: 388e30037aceca7679b5e1a74938c898ad1f4424
- Batch: 20261023-195359
- Enqueued verified: 388e30037aceca7679b5e1a74938c898ad1f4424
- At: 2026-10-23T20:04:40Z

Verdict: pass. The test now refuses only 127.0.0.1:4242, a case with 4242 in its import root proves a temp-directory path no longer fails it, and I-0197 is closed with a reason.
- 1: flai/cmd/serve_actions_test.go
- 2: design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md
