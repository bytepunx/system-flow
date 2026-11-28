---
id: T-3652
type: task
nature: improvement
title: The curate scheduler matches an acceptance by the project's workflow, not the baseline's
status: backlog
parent: S-0851
owner: alex
created: 2026-11-28T00:03:37Z
updated: 2026-11-28T00:03:37Z
transitions: []
stream: S-0851
tags: [flai, workflow]
touches: [flai/internal/serve/curate_schedule.go, flai/internal/serve/curate_schedule_test.go, flai/internal/serve/orchestrate_reasons_test.go]
---
# T-3652 The curate scheduler matches an acceptance by the project's workflow, not the baseline's

## Work

I-0406's third read is `acceptance()` in `flai/internal/serve/curate_schedule.go`. It builds the occurrence the curator's triggers are matched against from `baselineFlow.WithEffect(workflow.TypeStory, workflow.EffectAccept)`. Its comment says roles' conditions name the baseline's states until S-0748. S-0748 is done, and since ADR-0304 a role file names the project's own states, so the reason the comment gives no longer holds. Its one caller, the scheduler's look, has the repository at hand: it calls `acceptedIn(repo.Flow(), ...)` on the next line.

- Give `acceptance` the workflow, as `acceptanceIn(wf workflow.Workflow)` or a parameter, and pass `repo.Flow()` from the look. Take `To` from the transition with the `accept` effect in that workflow. When the workflow has none, give no `To`, as the baseline never does.
- Rewrite its comment to say it reads the project's workflow, and drop the S-0748 sentence.
- Remove `accepted()`, the baseline wrapper of `acceptedIn` whose comment names kluster.go: since S-0843 no code calls it, only `curate_schedule_test.go`. Have that test call `acceptedIn` with `workflow.Baseline()`.
- Update the test of `acceptance()` in `flai/internal/serve/orchestrate_reasons_test.go` to the new form. Add a test that builds a workflow renaming done, as S-0843's tests build one, and checks that the occurrence's `To` is the renamed state.

It waits for nothing. It shares no path with the components task, so the two run together.

## Done when

- `grep -n baselineFlow flai/internal/serve/curate_schedule.go` finds nothing.
- A test gives the scheduler's acceptance a workflow whose accepted state is not `done` and gets that state as `To`. It passes, and the baseline still gives `done`.
- `flai test` on the three paths passes.

## Notes
