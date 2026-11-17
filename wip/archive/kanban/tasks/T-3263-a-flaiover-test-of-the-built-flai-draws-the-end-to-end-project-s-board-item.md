---
id: T-3263
type: task
nature: feature
title: A flaiover test of the built flai draws the end-to-end project's board, item.move, and charts with the added state and its timer
status: done
parent: S-0756
owner: alex
created: 2026-11-12T16:10:32Z
updated: 2026-11-17T11:34:57Z
transitions:
  - to: ready
    at: 2026-11-17T11:29:38Z
    by: agent-S-0756
  - to: in-progress
    at: 2026-11-17T11:29:39Z
    by: agent-S-0756
  - to: done
    at: 2026-11-17T11:34:57Z
    by: agent-S-0756
stream: S-0756
tags: [dashboard]
touches: [flaiover/src/lib/server/workflow.test.ts]
after: [T-3261]
usage:
  source: log
  seconds: 318
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 10339
      cache_read: 2199299
      cache_write: 79349
      cost: 1.1145
---
# T-3263 A flaiover test of the built flai draws the end-to-end project's board, item.move, and charts with the added state and its timer

## Work

Waits for T-3261: it uses the project and workflow file that task's fixture defines. Write `workflow.test.ts`, importing `flaiover/src/lib/server/testing.ts` so the `flaiover-flai` tier runs it against the built flai: open the end-to-end project, read `board.get` and check the new column, make and refuse a move with `item.move`, and read `stats.get` and the chart data for the new timer and for time in state with the new state stacked.

## Done when

The test passes in the `flaiover-flai` tier against the built flai.

## Notes
