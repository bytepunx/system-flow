---
id: T-1660
type: task
nature: feature
title: "The dashboard's queue panel shows each bisection step and each story's outcome: landed, sent back, or still enqueued"
status: backlog
parent: S-0400
owner: alex
created: 2026-10-16T05:02:11Z
updated: 2026-10-16T05:02:11Z
transitions: []
stream: S-0400
tags: [dashboard]
touches: [flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts]
after: [T-1644]
---
# T-1660 The dashboard's queue panel shows each bisection step and each story's outcome: landed, sent back, or still enqueued

## Work

- In `MergeQueue.svelte`, list the batch's bisection steps, each prefix with its result, and each story's last outcome with its time.
- Link a sent-back story to its thread.
- Test the steps and the outcomes from a fixture state.
- Waits for T-1644, which defines the steps' shape in the state.

## Done when

The panel shows the steps and outcomes, and the vitest tests pass.

## Notes
