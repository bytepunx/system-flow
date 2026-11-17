---
id: T-3244
type: task
nature: feature
title: board.get answers the workflow's columns with their kinds, WIP limits, and order, and its transitions with who may make each and its reason
status: done
parent: S-0753
owner: alex
created: 2026-11-12T16:09:08Z
updated: 2026-11-17T03:49:01Z
transitions:
  - to: ready
    at: 2026-11-17T03:31:44Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T03:31:44Z
    by: system-flow
  - to: done
    at: 2026-11-17T03:49:01Z
    by: agent-S-0753
stream: S-0753
tags: [cli]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go]
usage:
  source: log
  seconds: 1037
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 72
      output: 344
      cache_read: 2534293
      cache_write: 279787
      cost: 1.0712
---
# T-3244 board.get answers the workflow's columns with their kinds, WIP limits, and order, and its transitions with who may make each and its reason

## Work

Add to the board view, and so to `board.get`, the workflow's transitions, each with from, to, the item types it takes, who may make it, and whether it needs a reason, beside the columns S-0739 gives with their kinds, WIP limits, and order. Keep every field the answer has today. Test the answer on the baseline and on a workflow with a state added.

## Done when

`board.get` answers the columns and transitions, and `hostapi_test.go` and flaiover's tests of the built flai pass.

## Notes
