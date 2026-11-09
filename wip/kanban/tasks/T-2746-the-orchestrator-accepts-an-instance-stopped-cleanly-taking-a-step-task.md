---
id: T-2746
type: task
nature: improvement
title: The orchestrator accepts an instance stopped cleanly, taking a step task cancelled by the stop as its criterion's evidence
status: backlog
parent: S-0621
owner: alex
created: 2026-11-09T17:06:59Z
updated: 2026-11-09T17:06:59Z
transitions: []
stream: S-0621
tags: [flai]
touches: [flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/cmd/accept.go, flai/cmd/accept_formula_test.go]
after: [T-2745]
---
# T-2746 The orchestrator accepts an instance stopped cleanly, taking a step task cancelled by the stop as its criterion's evidence

## Work

Today `stepsUndone` in `flai/internal/preview/orchestrator.go` (about lines 353 to 384) wants every criterion's task done, so a step cancelled by a stop is the blocker `criterion_unevidenced` even with its criterion ticked (I-0283's second instance).

- Take the evidence form T-2743's ADR settles, proposed `- <n>: T-nnnn skipped`, for a step task cancelled with the stop's reason by `flai formula stop`. Refuse it for a task cancelled any other way, or not cancelled, naming it. Keep `- <n>: T-nnnn done` as it is.
- Update the evidence regex (about line 98) and the error text in `flai/cmd/accept.go` (about lines 1205 to 1210), which names the forms.
- Fix the third instance as T-2743 found its cause: S-0585's branch was gone at acceptance, so the dry run answered "the operator accepts it" (`BlockUnverified`, about line 296). Make an instance stopped cleanly, whose branch changes nothing, acceptable by the orchestrator. Widen the touches with any file the fix needs.
- Tests: in `flai/internal/preview/accept_test.go`, beside `TestStepsUndoneNamesEachCriterionWhoseTaskIsNotDone`, a skipped step's evidence passes and a task cancelled by hand does not. In `flai/cmd/accept_formula_test.go`, an end-to-end test reproduces I-0283: a publish-like instance stopped at its evaluate step with `flai formula stop`, closed out, is accepted by `flai accept --by orchestrator --verified <commit> --evidence -` with no blocker.

This task waits for T-2745, whose cancel reason the evidence check reads.

## Done when

- An instance stopped with `flai formula stop` is accepted by the orchestrator with no blocker and no operator, with a test that reproduces I-0283's second and third instances.
- A step task cancelled by hand is still `criterion_unevidenced`, with a test.
- `flai test` on the changed paths passes.

## Notes
