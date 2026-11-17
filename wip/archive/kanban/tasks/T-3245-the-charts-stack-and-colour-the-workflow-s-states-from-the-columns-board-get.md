---
id: T-3245
type: task
nature: feature
title: The charts stack and colour the workflow's states from the columns board.get answers, in place of STATES and COD_COLUMNS
status: done
parent: S-0754
owner: alex
created: 2026-11-12T16:09:11Z
updated: 2026-11-17T09:59:37Z
transitions:
  - to: ready
    at: 2026-11-17T09:48:27Z
    by: agent-S-0754
  - to: in-progress
    at: 2026-11-17T09:48:30Z
    by: agent-S-0754
  - to: done
    at: 2026-11-17T09:59:37Z
    by: agent-S-0754
stream: S-0754
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
usage:
  source: log
  seconds: 667
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 27630
      cache_read: 5745066
      cache_write: 148227
      cost: 2.6834
---
# T-3245 The charts stack and colour the workflow's states from the columns board.get answers, in place of STATES and COD_COLUMNS

## Work

Replace `charts.ts`'s `STATES` and `COD_COLUMNS` with the columns `board.get` answers, in the workflow's order, coloured by kind: time in state, cumulative flow, and cost of delay per column. Test time in state and cumulative flow for a workflow with a state added between `in-progress` and `review`, and that the baseline draws as today.

## Done when

`charts.ts` holds no list of states, and `charts.test.ts` passes on both workflows.

## Notes
