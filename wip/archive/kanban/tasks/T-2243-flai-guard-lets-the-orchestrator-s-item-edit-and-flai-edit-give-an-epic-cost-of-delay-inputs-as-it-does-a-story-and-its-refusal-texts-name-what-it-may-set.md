---
id: T-2243
type: task
nature: improvement
title: flai guard lets the orchestrator's item_edit and flai edit give an epic cost of delay inputs as it does a story, and its refusal texts name what it may set
status: done
parent: S-0523
owner: alex
created: 2026-10-28T14:32:11Z
updated: 2026-10-30T14:03:35Z
transitions:
  - to: ready
    at: 2026-10-30T13:07:23Z
    by: agent-S-0523
  - to: in-progress
    at: 2026-10-30T13:07:24Z
    by: agent-S-0523
  - to: done
    at: 2026-10-30T14:03:35Z
    by: agent-S-0523
stream: S-0523
tags: [flai, guard, orchestrator, cost-of-delay]
touches: [flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/guard_test.go]
usage:
  source: log
  seconds: 3371
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 15873
      cache_read: 4564768
      cache_write: 199173
      cost: 2.4036
---
# T-2243 flai guard lets the orchestrator's item_edit and flai edit give an epic cost of delay inputs as it does a story, and its refusal texts name what it may set

## Work

Second layer: waits for T-2241, since the refusal texts cite the ADR's rule. Runs beside T-2242, with which it shares no path.

In `flai/internal/guard/guard.go`, `setsCostInputs` (MCP `item_edit`) and `setsCostInputsOnly` (`flai edit`) admit a story ID alone through `isStory`; admit an epic ID too (`isEpic` exists), and nothing else: a task stays refused. Keep every other condition: one input or more, none null or empty, no value, no `clear_cost_of_delay`, and no field besides `cost_of_delay`, `hash`, `id`, and `project` (`costFields`), or besides the flags finalizing allows but `--no-draft`. Reword `editsMCP` and `editsCLI`, and the comments, from "give a story cost of delay inputs" to "give an epic or a story cost of delay inputs". The guard still does not see an item's status: whether the epic or story may take the inputs stays flai's to judge (T-2242).

Tests: in `flai/internal/guard/orchestrate_test.go`, `TestTheOrchestratorNeverPlansATaskEditsAnItemOrPlacesOneByHand` refuses `E-0001` with `cost_of_delay` today, which reproduces I-0187's first instance; move that case to `TestTheOrchestratorPlansFinalizesPromotesAndOrdersByItsPermissions` beside `S-0328`, allowed with the permission on and refused with it off, and keep a task refused as never. Update the pinned refusal texts in `TestTheOrchestratorsStoryPlanningRefusalsNamePlanBacklogStories` and `TestAnInputKnowsTheFieldsItGives` where they quote `editsMCP`. In `flai/internal/guard/guard_test.go`, the CLI case that allows `flai edit S-nnnn --time-lost-per-cycle` gains an epic twin.

## Done when

- `go test ./internal/guard/...` passes in `flai/`, and `scripts/flai.sh test` on the changed paths is clean.
- Under `FLAI_ROLE=orchestrate` with `plan_backlog_stories` on, `flai guard` lets through `item_edit` with `id` an epic's and `cost_of_delay` inputs alone, and `flai edit E-nnnn --time-lost-per-cycle 4h`; with the permission off both are refused naming `plan_backlog_stories`.
- A task's ID, a value, `clear_cost_of_delay`, an empty input, and any other field are still refused as never, and the refusal text says the orchestrator gives an epic or a story cost of delay inputs.

## Notes

Drafted by planner-S-0523. `setsCostInputs` is at `guard.go` near line 1035 and `setsCostInputsOnly` near line 1088; `editsMCP` and `editsCLI` near line 492.
