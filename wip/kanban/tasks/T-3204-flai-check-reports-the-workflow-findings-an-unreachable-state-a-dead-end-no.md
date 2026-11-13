---
id: T-3204
type: task
nature: feature
title: "flai check reports the workflow findings: an unreachable state, a dead end, no pull, no acceptance, a bad timer, an unknown rule, and an item in a missing state"
status: backlog
parent: S-0741
owner: alex
created: 2026-11-12T16:06:27Z
updated: 2026-11-12T16:06:27Z
transitions: []
stream: S-0741
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go, flai/internal/check/sequence_test.go, flai/internal/check/vocabulary_test.go]
---
# T-3204 flai check reports the workflow findings: an unreachable state, a dead end, no pull, no acceptance, a bad timer, an unknown rule, and an item in a missing state

## Work

Add `check_workflow.go`, modelled on `check_roles.go`, and register it in `check.go`'s `sequence`. Report each finding under its own `workflow.` code, naming the state, transition, timer, or item and what to do: a state no transition reaches, a state with no path to an accepted or a cancelled state, no transition with the pull effect, none with the acceptance effect, a timer naming a state the workflow lacks, a rule the vocabulary lacks, and an item whose status the workflow lacks. Keep `sequence_test.go` and `vocabulary_test.go` in step with the new codes.

## Done when

`workflow_test.go` gives each finding on a workflow written to fail it, and the baseline gives none.

## Notes
