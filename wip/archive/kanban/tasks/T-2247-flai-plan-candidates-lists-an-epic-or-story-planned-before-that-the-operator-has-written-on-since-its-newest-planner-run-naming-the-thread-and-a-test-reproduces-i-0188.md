---
id: T-2247
type: task
nature: improvement
title: flai plan --candidates lists an epic or story planned before that the operator has written on since its newest planner run, naming the thread, and a test reproduces I-0188
status: done
parent: S-0524
owner: alex
created: 2026-10-28T14:41:23Z
updated: 2026-10-30T13:32:44Z
transitions:
  - to: ready
    at: 2026-10-30T13:08:52Z
    by: agent-S-0524
  - to: in-progress
    at: 2026-10-30T13:08:54Z
    by: agent-S-0524
  - to: done
    at: 2026-10-30T13:32:44Z
    by: agent-S-0524
stream: S-0524
tags: [flai, orchestrator, planner]
touches: [flai/internal/workitem/plancandidates.go, flai/internal/workitem/plancandidates_test.go, flai/cmd/plan.go, flai/cmd/plan_candidates_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 1430
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 19311
      cache_read: 4037909
      cache_write: 174311
      cost: 2.2107
---
# T-2247 flai plan --candidates lists an epic or story planned before that the operator has written on since its newest planner run, naming the thread, and a test reproduces I-0188

## Work

First layer: waits for nothing; it runs beside the ADR task, which shares no path with it.

Add the third kind of candidate, *planned before, asked again*, to `flai plan --candidates`:

- In `flai/internal/workitem/plancandidates.go`, `PlanCandidatesOf` lists, after the epics and the stories of today, each epic that is not done, cancelled, or archived and that neither kind lists, and each story in the backlog or ready that `storyPlanLacks` finds complete, when the caller says the operator has written on it since its newest planner run. Keep the package free of serve and threads: take the judgement as an input, such as a map from item ID to the reason text and thread, filled by `flai/cmd/plan.go`, the way `planHeld` is passed today. `PlanCandidate` gains `thread` (the thread ID) beside `id`, `type`, `title`, `status`, and `reason`, empty for the other kinds, so the orchestrator and the planner's start can name it.
- In `flai/cmd/plan.go`, beside `planHeld`, `plannerAsked`, and `updatedAfter`, work out for each epic and story with an ended planner run in `serve.AgentState.Plans` whether a thread anchored on the item, or on an item under it (an epic's stories and their tasks, a story's tasks), has an entry after the run's `Ended` by an author `threads.OperatorName` recognises for that story (the manifest's owner, the story's owner, or `designer`). Take the newest such entry; the reason reads `planned before; <author> wrote on <TH-nnnn> (<anchor item>) at <time>, after its planner ran until <ended>`. The holds stay as they are: a running planner for the item or its epic, or a planner's question awaiting the operator (`plannerAsked`), keep it in `left_out`. A story planned before whose own words changed but that nobody wrote on stays left out as today, since the replanner plans it again behind the `plan` host action.
- The run records that count are the planner's for the item (`st.Plans[id]`), as `planHeld` reads them. A run that has not ended counts as running, as today.

Tests:

- `flai/internal/workitem/plancandidates_test.go`: an epic with open stories and a complete backlog story are listed with the reason given when the input names them, and left out when it does not; the two kinds of today are unchanged.
- `flai/cmd/plan_candidates_test.go`: reproduce I-0188 on a fixture: an epic with stories whose planner run ended, then an entry by the manifest's owner on a thread anchored on one of its stories after the run's end; before the change the epic is not listed, after it the epic is listed with the thread in `thread` and the reason above. Also: an entry by the planner or the orchestrator on the same thread lists nothing; an entry older than the run's end lists nothing; a running planner for the epic keeps it in `left_out`.

## Done when

- `flai plan --candidates` and `flai plan --candidates --json` list the third kind with `reason` and `thread`, and the tests above pass with `go test ./internal/workitem/... ./cmd/...` in `flai/`.
- The fixture test reproducing I-0188 fails on the code before this task and passes after it.
- `flai test` on the changed paths passes.

## Notes

The reason text and the `thread` field are what the next task passes to the planner's start and the orchestrator reads; keep them stable.
