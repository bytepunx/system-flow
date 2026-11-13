---
id: T-3220
type: task
nature: feature
title: The merge queue takes stories in the state the acceptance leaves from and sends a bisected one back by the send-back effect
status: backlog
parent: S-0743
owner: alex
created: 2026-11-12T16:07:20Z
updated: 2026-11-12T16:07:20Z
transitions: []
stream: S-0743
tags: [cli]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go]
---
# T-3220 The merge queue takes stories in the state the acceptance leaves from and sends a bisected one back by the send-back effect

## Work

Replace the references to a state's name in `serve/mergequeue.go` and `flai/internal/mergequeue`: a story is eligible in the state the acceptance leaves from, a landed one is in an accepted state, and a story a red batch is bisected to goes back along the send-back transition (ADR-0176). Test eligibility and the send back on a workflow with a state added.

## Done when

These files name no state, and their tests pass on both workflows.

## Notes
