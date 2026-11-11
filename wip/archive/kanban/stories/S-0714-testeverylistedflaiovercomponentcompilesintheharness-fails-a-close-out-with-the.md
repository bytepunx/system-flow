---
id: S-0714
type: story
nature: remediation
title: "TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with \"the setup src/routes/layout.css does not compile: 404, 404 page not found\" while the main checkout's flaiover node_modules is intact"
status: done
owner: alex
created: 2026-11-10T10:54:41Z
updated: 2026-11-11T21:56:03Z
transitions:
  - to: ready
    at: 2026-11-11T16:15:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-11T18:16:21Z
    by: agent-S-0714
  - to: review
    at: 2026-11-11T19:10:21Z
    by: agent-S-0714
  - to: done
    at: 2026-11-11T21:56:03Z
    by: alex
tags: []
touches: [flai/internal/components/devserver.go, flai/internal/components/devserver_test.go, flai/internal/components/harness.go, flai/internal/components/harness_test.go, flai/internal/components/harness/plugin.mjs, flai/internal/components/serve.go, flai/internal/components/flaiover_test.go, design/system/component-library.md, design/issues/I-0342-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md, design/issues/summary.md, flai/cmd/components_test.go, flai/internal/serve/components_test.go, design/issues/I-0347-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3261
  priced_by: harness
  turns:
    - day: 2026-11-11
      ceremony: 1
      test_runs: 3
      work: 41
  models:
    - model: claude-opus-5-5
      input: 158
      output: 58544
      cache_read: 14246542
      cache_write: 736817
      cost: 8.6992
  strategic:
    - kind: planner
      seconds: 4
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 7
          output: 999
          cache_read: 400440
          cache_write: 10446
          cost: 0.1649
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 16
          cache_read: 77254
          cache_write: 9334
          cost: 0.0244
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-10T10:54:41Z
  value: 25
  by: planner-S-0714
  at: 2026-11-11T09:06:26Z
forecast:
  duration: 45m
  delivery: 2026-11-11T19:18:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 3, behind S-0715 and S-0707."
  by: flai
  at: 2026-11-11T18:16:00Z
finalized:
  by: alex
  at: 2026-11-11T09:02:34Z
---
# S-0714 TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with "the setup src/routes/layout.css does not compile: 404, 404 page not found" while the main checkout's flaiover node_modules is intact

## Goal

This story remediates [I-0342](../../../design/issues/I-0342-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md), "TestEveryListedFlaioverComponentCompilesInTheHarness fails a close-out with "the setup src/routes/layout.css does not compile: 404, 404 page not found" while the main checkout's flaiover node_modules is intact". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0342 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0342 is closed with `flai issue close I-0342 --reason` saying what fixed it

## Tasks
- T-3114 The components dev server takes its port from outside the host's ephemeral range, so no other process's port 0 listener can be given it
- T-3115 The harness marks each answer with a token flai gave it, and the dev server's wait takes only an answer that carries it
- T-3116 TestEveryListedFlaioverComponentCompilesInTheHarness checks each answer is its own vite's and says so when another server answered
- T-3117 component-library.md says where the dev server's port comes from and how flai knows the answer is its vite's, and I-0342 is closed

## Notes

Cost of delay inputs set by flai from I-0342. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-10T10:40:56Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed cause, from I-0342's one instance. The body `404 page not found` is Go's `net/http` default, so a Go server answered on the port where vite had answered the harness page with 200 a moment before. `freePort` in `flai/internal/components/devserver.go` takes the port from a listen on `127.0.0.1:0`, which is the kernel's ephemeral range, 32768 to 60999 on this host. Every `httptest` server and every outgoing connection of the close-outs running at the same time draws from that range. Once vite stops listening, by exiting or by restarting its server, one of them can be given the port. This is not checked against vite's last lines, which the instance does not record.

Proposed fix, in tasks:

- T-3114: pick the port below the ephemeral range, which the kernel never hands out.
- T-3115: the harness marks every answer with a per-start token, and the start's wait checks it. Its test reproduces the cause with a Go server on the picked port.
- T-3116: the integration test checks the token on every answer, and on a mismatch says another server answered.
- T-3117: the design is updated, and the issue is closed.

Touches:

- Layout: `flai/internal/components/devserver.go` and `devserver_test.go` hold `freePort`, `await`, and their tests.
- Layout: `harness.go`, `harness_test.go`, and `harness/plugin.mjs` write the harness and its plugin. `serve.go` passes options from `Serve` to both.
- Declared by the failing test's path: `flai/internal/components/flaiover_test.go`.
- Design: `design/system/component-library.md` § the dev server's free port, which says the gap cannot be closed.
- Criteria: the I-0342 file and `design/issues/summary.md`, through `flai issue close`.
- Left out: the co-change suggestions from `flai touches suggest`. They are hub and issue files that this change does not reach, such as `docs/users/flai.md` and `design/system/flai-cli.md`, because no command or user-facing behaviour changes.

No folder touch. No `after`: the story reads nothing that another open story writes.

Forecast: 45m. `flai forecast` gave 27m, from 135 s per unit of size over 59 medium remediation stories, at size 12. I raised it because the fix spans Go and the harness's JavaScript plugin. Its proof also needs real vite runs in the integration tier, which took 8m39s to reach in I-0342's close-out. Delivery moves by the same 18m, from 22:25Z to 22:43Z.

Cost of delay: 25 USD a week, as `flai cod` gives it from the 10m `time_lost_per_cycle` flai set from I-0342, at 150 USD an hour, once a week. It stands: the issue has one occurrence. A second close-out that fails this way would raise it.
