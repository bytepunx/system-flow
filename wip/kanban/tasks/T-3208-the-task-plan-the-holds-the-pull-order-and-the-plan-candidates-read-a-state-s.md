---
id: T-3208
type: task
nature: feature
title: The task plan, the holds, the pull order, and the plan candidates read a state's kind
status: done
parent: S-0739
owner: alex
created: 2026-11-12T16:06:41Z
updated: 2026-11-15T00:55:05Z
transitions:
  - to: ready
    at: 2026-11-14T23:56:43Z
    by: agent-S-0739
  - to: in-progress
    at: 2026-11-15T00:26:04Z
    by: agent-S-0739
  - to: done
    at: 2026-11-15T00:55:05Z
    by: agent-S-0739
stream: S-0739
tags: [cli]
touches: [flai/internal/workitem/taskplan.go, flai/internal/workitem/taskplan_test.go, flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/order.go, flai/internal/workitem/plancandidates.go, flai/internal/workitem/order_test.go, flai/internal/workitem/plancandidates_test.go]
after: [T-3203]
usage:
  source: log
  seconds: 1741
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 34169
      cache_read: 3661727
      cache_write: 198916
      cost: 2.5081
---
# T-3208 The task plan, the holds, the pull order, and the plan candidates read a state's kind

## Work

Waits for T-3203: its tests use the shared test workflow that task adds. Replace the references to a state's name in `taskplan.go`, `hold.go`, `order.go`, and `plancandidates.go` with reads of kinds and of which states are ordered: a task is done or waiting by kind, a ready story is held against stories in active-work states, the pull order lives in the ordered states, and the planner's candidates are drawn from the queue states. Test each on the baseline and on the fixture.

## Done when

The four files name no state and their tests pass on both workflows.

## Notes
