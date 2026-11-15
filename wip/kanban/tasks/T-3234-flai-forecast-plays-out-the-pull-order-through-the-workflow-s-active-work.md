---
id: T-3234
type: task
nature: feature
title: flai forecast plays out the pull order through the workflow's active-work states and their WIP limit
status: done
parent: S-0749
owner: alex
created: 2026-11-12T16:08:26Z
updated: 2026-11-15T19:57:35Z
transitions:
  - to: ready
    at: 2026-11-15T19:45:27Z
    by: agent-S-0749
  - to: in-progress
    at: 2026-11-15T19:45:27Z
    by: agent-S-0749
  - to: done
    at: 2026-11-15T19:57:35Z
    by: agent-S-0749
stream: S-0749
tags: [cli]
touches: [flai/internal/planning/forecast.go, flai/internal/planning/forecast_test.go, flai/cmd/forecast.go]
usage:
  source: log
  seconds: 706
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 50
      output: 228
      cache_read: 1619059
      cache_write: 171684
      cost: 0.6766
---
# T-3234 flai forecast plays out the pull order through the workflow's active-work states and their WIP limit

## Work

Replace the 9 references to a state's name in `forecast.go`: the stories in active-work states hold their places, then the stories in the ordered queue states start in order as the active-work states' WIP limit leaves room, each after the stories in its `after`. Test a forecast on a workflow with a state added between `in-progress` and `review`, whose stories hold a place as `in-progress` ones do when it is of the active-work kind.

## Done when

`forecast.go` names no state, and `forecast_test.go` passes on both workflows.

## Notes
