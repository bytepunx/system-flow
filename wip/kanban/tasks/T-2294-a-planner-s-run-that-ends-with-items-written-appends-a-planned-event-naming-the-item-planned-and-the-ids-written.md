---
id: T-2294
type: task
nature: feature
title: A planner's run that ends with items written appends a planned event naming the item planned and the IDs written
status: in-progress
parent: S-0535
owner: alex
created: 2026-10-28T21:38:09Z
updated: 2026-10-28T22:45:26Z
transitions:
  - to: ready
    at: 2026-10-28T22:45:25Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:45:26Z
    by: system-flow
stream: S-0535
tags: [flai, events, serve, planner]
touches: [flai/internal/serve/plan.go, flai/internal/serve/plan_test.go]
after: [T-2292]
---
# T-2294 A planner's run that ends with items written appends a planned event naming the item planned and the IDs written

## Work

`planEnded` in `flai/internal/serve/plan.go` judges a planner's run, lists the items it wrote with `plannedItems`, and logs the activity. `settleOrphans` in `agents.go` calls it too, so a run started from the CLI and handed over is covered.

- In `planEnded`, when `plannedItems` lists more than the item planned, append a `planned` event through T-2292's log: `id` the item planned, the IDs written, `by` the planner's agent name, and the run's outcome.
- Append nothing for a run that wrote nothing beyond the item, or that failed before writing.

It waits for T-2292, whose log it appends to. It shares no path with T-2295 or T-2296, so the three run together.

## Done when

- A test ends a planner run that wrote a story's tasks and reads one `planned` event with the story and the task IDs.
- A test ends a run that wrote nothing and reads no event.
- `flai test` on the changed paths passes.

## Notes

ADR-0209 §1 starts the orchestrator on `planned` to settle the planner's thread and finalize what it drafted.
