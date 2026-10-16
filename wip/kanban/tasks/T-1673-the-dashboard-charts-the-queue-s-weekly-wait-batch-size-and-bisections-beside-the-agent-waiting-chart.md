---
id: T-1673
type: task
nature: feature
title: The dashboard charts the queue's weekly wait, batch size, and bisections beside the agent-waiting chart
status: backlog
parent: S-0404
owner: alex
created: 2026-10-16T05:03:11Z
updated: 2026-10-16T05:03:11Z
transitions: []
stream: S-0404
tags: [dashboard]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-1664]
---
# T-1673 The dashboard charts the queue's weekly wait, batch size, and bisections beside the agent-waiting chart

## Work

- Add `queue` to the report types in `charts.ts` and a chart kind for it. It shows per week the mean and longest queue wait, the mean batch size, and the bisections.
- On the charts page, place it beside the agent-waiting chart. Say so when the flai on the host sends no `queue`, as the coordination chart does.
- Extend the port fixture in `charts.test.ts` with the new fields.
- Waits for T-1664, which defines the report's fields.

## Done when

The chart draws from a fixture report, the older-flai note shows without one, and the vitest tests pass.

## Notes
