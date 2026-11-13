---
id: T-3232
type: task
nature: feature
title: flai guard refuses an agent a move by who may make its transition, with today's refusals on the baseline
status: backlog
parent: S-0749
owner: alex
created: 2026-11-12T16:08:20Z
updated: 2026-11-12T16:08:20Z
transitions: []
stream: S-0749
tags: [cli]
touches: [flai/internal/guard/guard_rules.go, flai/internal/guard/guard.go, flai/internal/guard/guard_refusals.go, flai/internal/guard/call.go, flai/internal/guard/guard_test.go]
---
# T-3232 flai guard refuses an agent a move by who may make its transition, with today's refusals on the baseline

## Work

Replace the references to a state's name in the guard's files: a move an agent's role may not make is refused by the transition's who, read from the workflow, with today's refusal text on the baseline. Leave `guard/roles.go`'s role conditions to S-0748. Test a refused move into a state added between `in-progress` and `review`.

## Done when

These files name no state, and `guard_test.go` passes on both workflows.

## Notes
