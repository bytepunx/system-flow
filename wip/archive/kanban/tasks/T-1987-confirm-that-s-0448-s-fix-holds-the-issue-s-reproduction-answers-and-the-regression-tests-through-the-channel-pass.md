---
id: T-1987
type: task
nature: remediation
title: "Confirm that S-0448's fix holds: the issue's reproduction answers and the regression tests through the channel pass"
status: done
parent: S-0475
owner: alex
created: 2026-10-21T12:06:49Z
updated: 2026-10-21T12:11:22Z
transitions:
  - to: ready
    at: 2026-10-21T12:09:51Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T12:09:52Z
    by: system-flow
  - to: done
    at: 2026-10-21T12:11:22Z
    by: agent-S-0475
stream: S-0475
tags: [cli, dashboard]
usage:
  source: log
  seconds: 90
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 912
      cache_read: 336657
      cache_write: 18024
      cost: 0.2298
---
# T-1987 Confirm that S-0448's fix holds: the issue's reproduction answers and the regression tests through the channel pass

## Work

S-0448, accepted 2026-10-18, fixed the cause I-0163 records: the components methods take their sub-project as `sub_project` (382fab56f), and the dashboard's agent puts the channel's project key after a call's params (d94c114d4). This task confirms the fix still holds on main. It waits for nothing.

1. Run the issue's reproduction with the tree's flai: `scripts/flai.sh hostapi components.status '{"project":"sf"}'` and `scripts/flai.sh hostapi components.list '{"project":"sf"}'`. Each must answer, not refuse with `system-flow.yaml has no sub-project "sf"`.
2. Run the regression tests S-0448 added with `flai test` on these paths:
   - `flai/internal/hostapi/components_test.go`, which has `TestComponentsAnswerThroughTheChannelWithTheProjectKeyInTheirParams`
   - `flai/internal/hostapi/params_test.go`, which has `TestNoMethodTakesAParameterNamedProject`
   - `flaiover/src/lib/server/agent.test.ts`
   - `flaiover/src/routes/api/components/components.test.ts`
3. Record in the narrative's `## Decisions` what each step answered. These tests are the test that reproduces the cause, as the first criterion asks.
4. If a step fails, the cause recurs. Fix it, add a test that sends the call through the channel, and widen the task's touches to the files changed.

The task expects no file change, so it declares no touches.

## Done when

- Both `flai hostapi` calls with `{"project":"sf"}` answer without the sub-project error.
- The four test files pass under `flai test`.
- The narrative's `## Decisions` names the tests that cover the cause.

## Notes
