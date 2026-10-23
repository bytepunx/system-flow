---
id: T-1845
type: task
nature: feature
title: The dashboard's charts page draws threads per completed story per week with the repeats stacked
status: backlog
parent: S-0434
owner: alex
created: 2026-10-16T06:08:34Z
updated: 2026-10-23T17:23:57Z
transitions: []
stream: S-0434
tags: [dashboard]
touches: [flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-1843]
---
# T-1845 The dashboard's charts page draws threads per completed story per week with the repeats stacked

## Work

- `flaiover/src/lib/viz/charts.ts`: the `digest` types from `flai stats --json`, a chart kind `digest` with its label, and its builder: one bar per week of threads per completed story, split into the threads that asked again and the rest, beside the agent-waiting chart in the charts list.
- `flaiover/src/routes/charts/[kind]/+page.svelte`: render the kind, and the older-flai notice when the stats carry no `digest`.

Waits for T-1843, whose JSON shape it reads. Runs beside the CLI task.

## Done when

`charts.test.ts` builds the chart from a fixture, and `charts.svelte.test.ts` renders it and the older-flai notice; the flaiover tiers pass with `flai test flaiover/src/lib/viz flaiover/src/routes/charts`.

## Notes

Use the design system's tokens for the two series, as the other charts do (`palette.ts`).
