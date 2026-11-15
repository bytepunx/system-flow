---
id: T-3222
type: task
nature: feature
title: flai formula stop closes the stopped step and cancels the steps after it along the workflow's transitions
status: done
parent: S-0744
owner: alex
created: 2026-11-12T16:07:34Z
updated: 2026-11-15T01:35:52Z
transitions:
  - to: ready
    at: 2026-11-15T00:59:28Z
    by: agent-S-0744
  - to: in-progress
    at: 2026-11-15T01:10:11Z
    by: agent-S-0744
  - to: done
    at: 2026-11-15T01:35:52Z
    by: agent-S-0744
stream: S-0744
tags: [cli]
touches: [flai/cmd/formula_stop.go, flai/cmd/formula.go, flai/cmd/formula_stop_workflow_test.go]
after: [T-3221]
usage:
  source: log
  seconds: 1541
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 8289
      cache_read: 2133741
      cache_write: 52610
      cost: 0.9272
---
# T-3222 flai formula stop closes the stopped step and cancels the steps after it along the workflow's transitions

## Work

Replace the 10 references to a state's name in `formula_stop.go` and the 4 in `formula.go`: the stopped step's task closes along the task's transition to an accepted state, and the steps after it go along the transition with the cancellation effect (ADR-0271). Test a stop on a workflow with a state renamed.

## Done when

Both files name no state, and `formula_stop_test.go` passes on both workflows.

## Notes
