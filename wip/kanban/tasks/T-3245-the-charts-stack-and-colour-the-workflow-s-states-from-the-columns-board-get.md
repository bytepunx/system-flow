---
id: T-3245
type: task
nature: feature
title: The charts stack and colour the workflow's states from the columns board.get answers, in place of STATES and COD_COLUMNS
status: backlog
parent: S-0754
owner: alex
created: 2026-11-12T16:09:11Z
updated: 2026-11-12T16:09:11Z
transitions: []
stream: S-0754
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts]
---
# T-3245 The charts stack and colour the workflow's states from the columns board.get answers, in place of STATES and COD_COLUMNS

## Work

Replace `charts.ts`'s `STATES` and `COD_COLUMNS` with the columns `board.get` answers, in the workflow's order, coloured by kind: time in state, cumulative flow, and cost of delay per column. Test time in state and cumulative flow for a workflow with a state added between `in-progress` and `review`, and that the baseline draws as today.

## Done when

`charts.ts` holds no list of states, and `charts.test.ts` passes on both workflows.

## Notes
