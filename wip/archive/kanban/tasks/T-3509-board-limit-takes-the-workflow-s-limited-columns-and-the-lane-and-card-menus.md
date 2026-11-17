---
id: T-3509
type: task
nature: feature
title: board.limit takes the workflow's limited columns, and the lane and card menus' examples carry the flow
status: done
parent: S-0753
owner: alex
created: 2026-11-17T04:15:55Z
updated: 2026-11-17T04:21:57Z
transitions:
  - to: ready
    at: 2026-11-17T04:16:35Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T04:16:36Z
    by: system-flow
  - to: done
    at: 2026-11-17T04:21:57Z
    by: agent-S-0753
stream: S-0753
tags: []
touches: [flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_test.go, design/components/flaiover/src/lib/components/CardMenu.examples.json, design/components/flaiover/src/lib/components/LaneMenu.examples.json]
usage:
  source: log
  seconds: 321
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1249
      cache_read: 312522
      cache_write: 7595
      cost: 0.1336
---

# T-3509 board.limit takes the workflow's limited columns, and the lane and card menus' examples carry the flow

## Work

Found in the review of T-3249 and T-3250. `board.limit` in `flai/internal/hostapi/writes_items.go` checks the column against `workitem.LimitedColumns`, the baseline's, so the lane menu's **Change WIP limit** on a project's own limited state is refused before `flai board limit`, which checks the project's workflow, runs. Check it against the project's workflow instead, and test it with `testingWorkflow`. T-3249 gave `CardMenu` and `LaneMenu` a required `flow` prop, so their examples under `design/components/` gain the baseline's lanes and transitions. Waits for T-3249, whose components it follows.

## Done when

`board.limit` takes `testing` under `testingWorkflow` and still refuses `backlog`, the examples check clean, and the hostapi tests pass.

## Notes
