---
id: S-0543
type: story
nature: remediation
title: "TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with \"directory not empty\", as I-0048's serve tests did"
status: done
owner: alex
created: 2026-10-29T17:25:41Z
updated: 2026-11-07T15:15:20Z
transitions:
  - to: ready
    at: 2026-11-07T10:02:26Z
    by: alex
  - to: in-progress
    at: 2026-11-07T10:02:39Z
    by: agent-S-0543
  - to: review
    at: 2026-11-07T10:21:44Z
    by: agent-S-0543
  - to: done
    at: 2026-11-07T15:15:20Z
    by: alex
tags: [flai, serve, test]
touches: [flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, design/issues/I-0303-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1160
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 1
      test_runs: 3
      hand_edits: 2
      work: 23
  models:
    - model: claude-opus-5-5
      input: 66
      output: 12479
      cache_read: 6025263
      cache_write: 254928
      cost: 3.4943
  strategic:
    - kind: planner
      seconds: 141
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 40
          output: 10024
          cache_read: 2412943
          cache_write: 152281
          cost: 1.9015
    - kind: orchestrator
      seconds: 288
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 58
          output: 374
          cache_read: 1999270
          cache_write: 358221
          cost: 0.5712
        - model: claude-sonnet-5-5
          input: 14
          output: 103
          cache_read: 186256
          cache_write: 71169
          cost: 0.1827
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-10-29T17:25:41Z
  value: 25
  by: planner-S-0543
  at: 2026-11-02T18:12:31Z
forecast:
  duration: 30m
  delivery: 2026-11-07T13:03:00Z
  basis: "Its own forecast of 30m; 11th in the pull order with an in-progress limit of 5, behind S-0560, S-0562, S-0586, S-0406, S-0587, S-0588, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412 and S-0241."
  by: flai
  at: 2026-11-07T09:38:57Z
finalized:
  by: alex
  at: 2026-11-07T10:02:03Z
---
# S-0543 TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did

## Goal

This story remediates [I-0233](../../../design/issues/I-0233-testachangedfilereachesthedashboardasanotification-fails-its-temp-folder-s-cleanup-with-directory-not-empty-as-i-0048-s-serve-tests-did.md), "TestAChangedFileReachesTheDashboardAsANotification fails its temp folder's cleanup with "directory not empty", as I-0048's serve tests did". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0233 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0233 is closed with `flai issue close I-0233 --reason` saying what fixed it

## Tasks
- T-2435 Confirm with a race stress run that the serve tests no longer fail their temp folder's cleanup, and fix the writer if one still does

## Notes

Cost of delay inputs set by flai from I-0233. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-29T16:37:28Z, 0 days before this story; under one cycle counts as one).

### Planning

The outcome looks delivered already. I-0233 was closed on 2026-10-31T09:50:17Z with the reason that S-0418's fe589e207 made `Run` wait for each served project's client, loop of looks, and watcher before it returns, and that S-0548's c6c0310db (T-2362) replaced the test's 80 ms sleep with `watcherHasLooked`. Both commits are on main. The close records 10 of 10 and 50 of 50 passing runs under `GOMAXPROCS=2 go test -race`. No open issue under `design/issues` reports "directory not empty" now. Criterion 2 is met by that close; criterion 1 is left to one confirming task. The plan's thread recommends cancelling the story as delivered by S-0418 and S-0548.

Touches:

- `flai/internal/serve/serve.go`: design and layout. `Run`'s halt holds the wait I-0233 cites (the comment naming I-0233), so a fix the confirmation still needs would land here.
- `flai/internal/serve/serve_test.go`: design and layout. It holds the test the issue names and `watcherHasLooked`.
- `flai touches suggest S-0543` from those two files lists co-change with `flai/cmd/serve.go` (30%) and `flai/internal/serve/agents.go` (23%), among others. None of them is a writer the issue's instances name, so none is predicted.

Figures:

- Forecast: flai's figure is 4m, from 119 s per unit of size over 13 small remediation stories, times size 2. It is raised to 30m. A `-race -count` stress run of the serve package, plus a close-out whose go-test tier took 7 to 8 minutes on this package in both instances, does not fit in 4m. The delivery is flai's, moved by the extra 26m.
- Cost of delay: 25.00 USD a week, which is `flai cod` from flai's input of 10m lost per 168h cycle at 150 USD an hour. It stands as the inputs give it. With the cause fixed, that time is no longer being lost.

### Confirmation (T-2435)

On 2026-11-07, on story/S-0543 at 8cb8d4c5e, on a host at load average 54 to 71 on 24 CPUs:

- `GOMAXPROCS=2 go test -race -count=20` of TestAChangedFileReachesTheDashboardAsANotification and TestServeStartsADefinedRoleOnAnEventAndOnASchedule passed 40 of 40 runs.
- `GOMAXPROCS=2 go test -race -count=3 -skip '^TestGateway' ./internal/serve/` passed 3 of 3 whole-package runs.

No code changed. Criterion 1's test is the existing pair: `Run`'s wait in `serve.go` (fe589e207) and `watcherHasLooked` in `serve_test.go` (c6c0310db). Criterion 2 was met when S-0548 closed I-0233 on 2026-10-31T09:50:17Z, with a reason naming both commits.
