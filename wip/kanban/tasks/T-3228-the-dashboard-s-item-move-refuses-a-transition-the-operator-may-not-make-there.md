---
id: T-3228
type: task
nature: feature
title: The dashboard's item.move refuses a transition the operator may not make there, and an acceptance from the board makes the acceptance transition
status: backlog
parent: S-0747
owner: alex
created: 2026-11-12T16:08:01Z
updated: 2026-11-12T16:08:01Z
transitions: []
stream: S-0747
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_run.go, flai/internal/hostapi/writes_test.go]
---
# T-3228 The dashboard's item.move refuses a transition the operator may not make there, and an acceptance from the board makes the acceptance transition

## Work

Replace the references to a state's name in the write methods: `item.move` refuses a transition the operator may not make from the dashboard, naming the transition and who may, and a move along the acceptance transition goes through acceptance as today. The outcome and run files report moves in the workflow's names. Test a move and an acceptance on a workflow with a state added between `in-progress` and `review`.

## Done when

These files name no state, and `writes_test.go` passes on both workflows.

## Notes
