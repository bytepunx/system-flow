---
id: T-3221
type: task
nature: feature
title: flai move refuses a transition the mover may not make, naming it and who may, and cancels along the workflow's transitions
status: done
parent: S-0744
owner: alex
created: 2026-11-12T16:07:31Z
updated: 2026-11-15T01:09:57Z
transitions:
  - to: ready
    at: 2026-11-15T00:59:28Z
    by: agent-S-0744
  - to: in-progress
    at: 2026-11-15T00:59:28Z
    by: agent-S-0744
  - to: done
    at: 2026-11-15T01:09:57Z
    by: agent-S-0744
stream: S-0744
tags: [cli]
touches: [flai/cmd/move.go, flai/cmd/move_cancel.go, flai/cmd/move_workflow_test.go, flai/cmd/project.go, flai/internal/workitem/transition.go, flai/internal/preview/cancel.go, docs/users/flai.md, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 629
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 18101
      cache_read: 5634766
      cache_write: 165488
      cost: 2.6293
---
# T-3221 flai move refuses a transition the mover may not make, naming it and who may, and cancels along the workflow's transitions

## Work

Make `flai move` and `move_cancel.go` use the workflow package's mover check, which S-0737 adds, in place of their state names, through a helper in `workitem/transition.go` that `flai accept` uses too: a mover the transition does not name is refused with the transition, who may make it, and the command that does. The move to an accepted state stays `flai accept`'s. Test a refused and an allowed move on a workflow with a state added between `in-progress` and `review`.

## Done when

`move.go` and `move_cancel.go` name no state, and the move tests pass on both workflows.

## Notes
