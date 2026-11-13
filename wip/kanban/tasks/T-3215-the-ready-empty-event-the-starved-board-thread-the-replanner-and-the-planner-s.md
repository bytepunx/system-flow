---
id: T-3215
type: task
nature: feature
title: The ready-empty event, the starved-board thread, the replanner, and the planner's starts count by kind and WIP limit
status: backlog
parent: S-0742
owner: alex
created: 2026-11-12T16:07:01Z
updated: 2026-11-12T16:07:01Z
transitions: []
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/readyempty.go, flai/internal/serve/readyempty_test.go, flai/internal/serve/replan.go, flai/internal/serve/replan_test.go, flai/internal/serve/plan.go, flai/internal/serve/plan_test.go]
---
# T-3215 The ready-empty event, the starved-board thread, the replanner, and the planner's starts count by kind and WIP limit

## Work

Replace the references to a state's name in `readyempty.go`, `replan.go`, and `plan.go`: the `ready-empty` event and the starved-board thread (ADR-0285, ADR-0279) count the stories in the state the pull takes from against the room the active-work states' WIP limit leaves, and the replanner and the planner's starts pick their stories from the queue states. Leave `replan.go`'s use of the role conditions to S-0748. Test each on a workflow with a state renamed and one added.

## Done when

The three files name no state, and their tests pass on both workflows.

## Notes
