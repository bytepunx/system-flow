---
id: T-2763
type: task
nature: remediation
title: Test that the operator's confirmation of a cost of delay recommendation on a planner's thread gives the orchestrator a thread reason
status: done
parent: S-0622
owner: alex
created: 2026-11-09T18:05:16Z
updated: 2026-11-09T22:09:58Z
transitions:
  - to: ready
    at: 2026-11-09T21:55:51Z
    by: agent-S-0622
  - to: in-progress
    at: 2026-11-09T21:55:52Z
    by: agent-S-0622
  - to: done
    at: 2026-11-09T22:09:58Z
    by: agent-S-0622
stream: S-0622
tags: [cli]
touches: [flai/internal/serve/orchestrate_confirm_test.go]
usage:
  source: log
  seconds: 846
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 29
      output: 2892
      cache_read: 4030948
      cache_write: 82229
      cost: 1.3439
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 26247
          cache_write: 880
          cost: 0.0072
---
# T-2763 Test that the operator's confirmation of a cost of delay recommendation on a planner's thread gives the orchestrator a thread reason

## Work

Write the regression test for I-0284's first half: a confirmed input must reach the orchestrator. Put it in a new `flai/internal/serve/orchestrate_confirm_test.go`. Build the fixture like the tests in `orchestrate_reasons_test.go`, and like `TestAPlannersQuestionAwaitsTheOperatorInTheOrchestratorsRole` in `flai/internal/inbox/threadevents_test.go`.

1. Make a story in progress, as S-0549 was, and a thread on it opened by `planner-S-nnnn` that asks for `time_lost_per_cycle`.
2. Add the orchestrator's recommendation with `threads.ReplyWith`, then the operator's confirmation with `threads.Confirm`.
3. Run the orchestrator's look and assert that it queues the reason `thread:TH-nnnn`, as `reasonOf` and `occurrenceOf` in `orchestrate_reasons.go` give it.
4. Assert the same for flai serve's own confirmation after `confirm_after`, `threads.ConfirmByTimeout`.

The code path is read to hold: `confirm` sets the thread `answered`, the orchestrator is a participant, and `Summarize` makes the thread await it. If the test fails, fix the cause in the reason path. Say in the narrative's `## Decisions` what it was. The touches then widen to the fixed file when the task closes.

The task waits for nothing. It shares no path with the other test task, so the two run together.

## Done when

- `flai/internal/serve/orchestrate_confirm_test.go` holds the test. It fails when a confirmation queues no orchestrator reason, and passes on main as it is, or after the fix this task makes.
- `flai test flai/internal/serve/orchestrate_confirm_test.go` passes.

## Notes
