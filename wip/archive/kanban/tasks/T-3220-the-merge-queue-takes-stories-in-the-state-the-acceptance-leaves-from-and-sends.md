---
id: T-3220
type: task
nature: feature
title: The merge queue takes stories in the state the acceptance leaves from and sends a bisected one back by the send-back effect
status: done
parent: S-0743
owner: alex
created: 2026-11-12T16:07:20Z
updated: 2026-11-15T01:44:19Z
transitions:
  - to: ready
    at: 2026-11-15T01:07:35Z
    by: agent-S-0743
  - to: in-progress
    at: 2026-11-15T01:07:35Z
    by: agent-S-0743
  - to: done
    at: 2026-11-15T01:44:19Z
    by: agent-S-0743
stream: S-0743
tags: [cli]
touches: [flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/internal/mergequeue/batch.go, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/eligible.go, flai/internal/mergequeue/eligible_test.go, flai/internal/mergequeue/bisect.go, flai/internal/mergequeue/bisect_test.go, flai/internal/mergequeue/state.go]
after: [T-3369]
usage:
  source: log
  seconds: 2204
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 19727
      cache_read: 10280441
      cache_write: 250917
      cost: 3.9924
---
# T-3220 The merge queue takes stories in the state the acceptance leaves from and sends a bisected one back by the send-back effect

## Work

Replace the references to a state's name in `serve/mergequeue.go` and `flai/internal/mergequeue`: a story is eligible in the state the acceptance leaves from, a landed one is in an accepted state, and a story a red batch is bisected to goes back along the send-back transition (ADR-0176). Test eligibility and the send back on a workflow with a state added.

## Done when

These files name no state, and their tests pass on both workflows.

## Notes
