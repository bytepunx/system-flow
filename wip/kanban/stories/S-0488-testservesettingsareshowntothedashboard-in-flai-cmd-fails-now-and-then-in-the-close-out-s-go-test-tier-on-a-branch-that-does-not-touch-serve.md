---
id: S-0488
type: story
nature: remediation
title: TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve
status: backlog
owner: alex
created: 2026-10-23T17:53:08Z
updated: 2026-10-23T18:11:28Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/cmd/serve_actions_test.go, design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 97
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 295
          cache_read: 5583119
          cache_write: 16178
          cost: 1.2315
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
  delivery: 2026-10-24T12:12:00Z
  basis: "flai forecast gave 5m (131 s per unit of size over 10 done remediation stories on claude-opus-5-5 in the small band, size 2 before the touches were written); raised to 15m, since the story has a test to narrow and a reproducing case to add, then an issue close, each closed by one agent with the go-test tier between them; delivery is flai's 2026-10-24T12:02Z, 28th in the pull order, moved by the 10m added."
  by: planner-S-0488
  at: 2026-10-23T18:10:04Z
finalized:
  by: orchestrator
  at: 2026-10-23T18:11:28Z
---
# S-0488 TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve

## Goal

This story remediates [I-0197](../../../design/issues/I-0197-testservesettingsareshowntothedashboard-in-flai-cmd-fails-now-and-then-in-the-close-out-s-go-test-tier-on-a-branch-that-does-not-touch-serve.md), "TestServeSettingsAreShownToTheDashboard in flai/cmd fails now and then in the close-out's go-test tier on a branch that does not touch serve". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0197 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0197 is closed with `flai issue close I-0197 --reason` saying what fixed it

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
