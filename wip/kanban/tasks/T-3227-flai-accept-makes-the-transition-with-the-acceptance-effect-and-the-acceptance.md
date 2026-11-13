---
id: T-3227
type: task
nature: feature
title: flai accept makes the transition with the acceptance effect, and the acceptance preview lists the blockers its rules give
status: backlog
parent: S-0744
owner: alex
created: 2026-11-12T16:07:58Z
updated: 2026-11-12T16:11:00Z
transitions: []
stream: S-0744
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/cancel.go]
after: [T-3221]
---
# T-3227 flai accept makes the transition with the acceptance effect, and the acceptance preview lists the blockers its rules give

## Work

Waits for T-3221: the acceptance checks who may make its transition through the helper that task adds to `workitem/transition.go`. Replace the references to a state's name in `cmd/accept.go` and `flai/internal/preview`: `flai accept` makes the transition with the acceptance effect from the state it leaves, refusing a story elsewhere with today's words and that state's name, and the preview's blockers come from that transition's rules, word for word as today with the baseline. `preview/cancel.go` lists what goes with a cancellation by kind. Test accept and the preview on a workflow with a state added between `in-progress` and `review`.

## Done when

These files name no state, and `accept_test.go` in `cmd` and `preview` pass on both workflows.

## Notes
