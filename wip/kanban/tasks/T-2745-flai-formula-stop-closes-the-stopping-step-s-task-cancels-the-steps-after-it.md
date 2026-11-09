---
id: T-2745
type: task
nature: improvement
title: flai formula stop closes the stopping step's task, cancels the steps after it, and ticks their criteria as skipped, so the instance can go to review
status: backlog
parent: S-0621
owner: alex
created: 2026-11-09T17:06:48Z
updated: 2026-11-09T17:06:48Z
transitions: []
stream: S-0621
tags: [flai]
touches: [flai/cmd/formula_stop.go, flai/cmd/formula_stop_test.go, flai/cmd/formula.go, flai/internal/guard/guard.go, flai/internal/guard/formula_test.go]
after: [T-2744]
---
# T-2745 flai formula stop closes the stopping step's task, cancels the steps after it, and ticks their criteria as skipped, so the instance can go to review

## Work

Build the command T-2743's ADR settles, proposed `flai formula stop <S-nnnn> --step <id> --exit <n>`, in a new `flai/cmd/formula_stop.go`, registered in `newFormulaCmd` in `flai/cmd/formula.go`.

- Refuse, writing nothing, a story with no `formula:`, a step id its formula does not have, an exit its step's `stop` does not list, and a step whose task is not open or whose `after` tasks are not all done. Name the reason in each refusal.
- Otherwise, in one write through the items lock:
  - write the exit and the last 40 lines of the output into the step task's `## Notes`, and move it to done as `task_done` does, committing nothing on a branch that changes nothing;
  - cancel each step task after it, its layers' tasks transitively, with the reason `step <id> stopped the run with exit <n>`;
  - tick the criterion of the stopping step and of each skipped step, each skipped one marked with that reason, through the criteria code in `flai/internal/itemedit/criteria.go`;
  - log the stop in the story's narrative.
- Answer with `--json` the tasks closed and cancelled and the criteria ticked.
- In `flai/internal/guard/guard.go`, add `stop` to `cliAuthors["formula"]`, so a read-only sub-agent and the analyzer are refused it and the story's agent is not.
- Tests in `flai/cmd/formula_stop_test.go`: reproduce I-0283 with a three-step instance whose second step lists `stop: [2]`. Stopping at step 2 leaves step 2's task done, step 3's cancelled, and all three criteria ticked. Test each refusal. In `flai/internal/guard/formula_test.go`, a sub-agent is refused `flai formula stop` and the story's agent runs it.

This task waits for T-2744, whose `Step` field the command reads.

## Done when

- `flai formula stop` on a stopped instance leaves the stopping step's task done with its exit in its Notes, every later step task cancelled with the reason, and every criterion ticked, with a test that reproduces I-0283.
- Each refusal writes nothing and has a test; the guard refuses a sub-agent the command.
- `flai test` on the changed paths passes.

## Notes
