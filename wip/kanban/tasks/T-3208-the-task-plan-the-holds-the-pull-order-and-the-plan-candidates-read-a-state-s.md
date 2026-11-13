---
id: T-3208
type: task
nature: feature
title: The task plan, the holds, the pull order, and the plan candidates read a state's kind
status: backlog
parent: S-0739
owner: alex
created: 2026-11-12T16:06:41Z
updated: 2026-11-12T16:06:41Z
transitions: []
stream: S-0739
tags: [cli]
touches: [flai/internal/workitem/taskplan.go, flai/internal/workitem/taskplan_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/order.go, flai/internal/workitem/plancandidates.go]
after: [T-3203]
---
# T-3208 The task plan, the holds, the pull order, and the plan candidates read a state's kind

## Work

Waits for T-3203: its tests use the shared test workflow that task adds. Replace the references to a state's name in `taskplan.go`, `hold.go`, `order.go`, and `plancandidates.go` with reads of kinds and of which states are ordered: a task is done or waiting by kind, a ready story is held against stories in active-work states, the pull order lives in the ordered states, and the planner's candidates are drawn from the queue states. Test each on the baseline and on the fixture.

## Done when

The four files name no state and their tests pass on both workflows.

## Notes
