---
id: T-3215
type: task
nature: feature
title: The ready-empty event, the starved-board thread, the replanner, and the planner's starts count by kind and WIP limit
status: done
parent: S-0742
owner: alex
created: 2026-11-12T16:07:01Z
updated: 2026-11-15T01:55:04Z
transitions:
  - to: ready
    at: 2026-11-15T01:27:00Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-15T01:27:00Z
    by: agent-S-0742
  - to: done
    at: 2026-11-15T01:55:04Z
    by: agent-S-0742
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go]
after: [T-3370]
usage:
  source: log
  seconds: 1684
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 78
      output: 22322
      cache_read: 7570494
      cache_write: 173971
      cost: 2.9646
---
# T-3215 The ready-empty event, the starved-board thread, the replanner, and the planner's starts count by kind and WIP limit

## Work

Replace the references to a state's name in `readyempty.go`, `replan.go`, and `plan.go`: the `ready-empty` event and the starved-board thread (ADR-0285, ADR-0279) count the stories in the state the pull takes from against the room the active-work states' WIP limit leaves, and the replanner and the planner's starts pick their stories from the queue states. Leave `replan.go`'s use of the role conditions to S-0748. Test each on a workflow with a state renamed and one added.

## Done when

The three files name no state, and their tests pass on both workflows.

## Notes
