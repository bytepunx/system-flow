---
id: S-0447
type: story
nature: remediation
title: The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier
status: done
owner: alex
created: 2026-10-18T02:19:04Z
updated: 2026-10-18T03:10:57Z
transitions:
  - to: ready
    at: 2026-10-18T02:53:21Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-18T02:53:29Z
    by: agent-S-0447
  - to: review
    at: 2026-10-18T03:08:37Z
    by: agent-S-0447
  - to: done
    at: 2026-10-18T03:10:57Z
    by: alex
tags: [flai, components, tests]
touches: [flai/internal/components/devserver.go, flai/internal/components/devserver_test.go, design/system/component-library.md, design/issues/I-0162-the-components-dev-server-test-picks-a-free-port-and-closes-it-before-the-fake-vite-binds-it-so-a-concurrent-test-run-on-the-host-can-take-the-port-and-fail-the-integration-tier.md, design/issues/summary.md, design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 922
  priced_by: harness
  turns:
    - day: 2026-10-18
      ceremony: 1
      test_runs: 1
      hand_edits: 2
      work: 39
  models:
    - model: claude-opus-5-5
      input: 88
      output: 19801
      cache_read: 5482418
      cache_write: 159231
      cost: 2.7667
  strategic:
    - kind: planner
      seconds: 152
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 172
          output: 14099
          cache_read: 33892515
          cache_write: 101666
          cost: 10.2167
    - kind: orchestrator
      seconds: 1639
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 86
          output: 1759
          cache_read: 40032556
          cache_write: 38205
          cost: 10.0454
cost_of_delay:
  inputs:
    time_lost_per_cycle: 7m
    by: flai
    at: 2026-10-18T02:19:04Z
  value: 17.5
  by: planner-S-0447
  at: 2026-10-18T02:42:06Z
forecast:
  duration: 30m
  delivery: 2026-10-18T22:09:00Z
  basis: "flai forecast's 13m (107 s per unit of size over 23 medium remediation stories, size 7), raised to 30m because the close-out runs the integration tier, 5m alone per I-0162, and the fix needs two new tests; delivery is flai's, 39th in the pull order, moved by the 17m added."
  by: planner-S-0447
  at: 2026-10-18T02:42:06Z
finalized:
  by: orchestrator
  at: 2026-10-18T02:42:39Z
---
# S-0447 The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier

## Goal

This story remediates [I-0162](../../../design/issues/I-0162-the-components-dev-server-test-picks-a-free-port-and-closes-it-before-the-fake-vite-binds-it-so-a-concurrent-test-run-on-the-host-can-take-the-port-and-fail-the-integration-tier.md), "The components dev server test picks a free port and closes it before the fake Vite binds it, so a concurrent test run on the host can take the port and fail the integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0162 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0162 is closed with `flai issue close I-0162 --reason` saying what fixed it

## Tasks
- T-1902 StartDevServer starts Vite again on another free port when the one it picked was taken before Vite bound it
- T-1903 Close I-0162 with the retry that fixed it

## Notes

Cost of delay inputs set by flai from I-0162. time_lost_per_cycle 7m: 7m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-17T18:10:30Z, 0.3 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0162's one instance: Vite takes a port number, not a listener, so the gap between `freePort` closing its listener and Vite binding the port cannot be closed. Instead, `StartDevServer` starts Vite again on another free port when Vite exits early saying the port is in use, at most three ports in all. A test makes the port picker return a held port first, which reproduces I-0162 without a race. The same race can hit `flai components serve` and the dashboard's dev server on any busy host, not only in tests, so the fix belongs in `StartDevServer` rather than in the test.

All touches are files. No folder touch is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/components/devserver.go` | design (I-0162 names it) | `freePort` and `StartDevServer`, where the retry goes |
| `flai/internal/components/devserver_test.go` | layout | the fake Vite and the dev server tests; the reproducing tests go here |
| `design/system/component-library.md` | design | § The dev server says how the port is chosen |
| `design/issues/I-0162-…md` | layout | closed by `flai issue close` |
| `design/issues/summary.md` | layout | `flai issue close` updates it |

`flai touches suggest`, from the first two files, listed co-changed files at 50% over 4 commits: `flai/cmd/serve_actions.go`, `flai/internal/hostapi/components.go`, `flai/internal/hostapi/writes.go`, `flai/internal/serve/components.go`, `flai/internal/config/config.go`, `docs/operators/settings.md`, and their tests. None is kept. They changed with the features that added the dev server, and the retry stays inside `StartDevServer`, whose signature does not change. `flai/internal/serve/components_test.go`, where I-0162 failed, gets the fix through the supervisor and needs no change.

S-0448, in progress, also claims `design/system/component-library.md`. This story is held until S-0448 is accepted, which its place in the pull order makes likely first.

Figures:

- Forecast 30m, not flai's 13m: the close-out's integration tier takes 5m alone (I-0162), and the fix needs two new tests. Delivery is flai's 2026-10-18T21:52:00Z plus the 17m added.
- Cost of delay 17.50 USD a week, as `flai cod` works it out from the input of 7m lost per 168h cycle at 150 USD an hour. It stands: the input is flai's from the issue's one occurrence.
