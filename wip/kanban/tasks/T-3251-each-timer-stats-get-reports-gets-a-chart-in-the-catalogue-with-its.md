---
id: T-3251
type: task
nature: feature
title: Each timer stats.get reports gets a chart in the catalogue with its distribution over the window
status: done
parent: S-0754
owner: alex
created: 2026-11-12T16:09:40Z
updated: 2026-11-17T10:14:13Z
transitions:
  - to: ready
    at: 2026-11-17T09:59:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T09:59:52Z
    by: system-flow
  - to: done
    at: 2026-11-17T10:14:13Z
    by: agent-S-0754
stream: S-0754
tags: [dashboard]
touches: [flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", flaiover/src/lib/components/ChartExplainer.svelte]
after: [T-3245]
---
# T-3251 Each timer stats.get reports gets a chart in the catalogue with its distribution over the window

## Work

Waits for T-3245: both change `charts.ts`, and the timer charts use its columns. Read the timers section of `stats.get` in `stats.ts`, add a chart kind per timer to the catalogue in `routes/charts/+page.ts`, and draw each as cycle time is drawn, its distribution over the window, in `charts/[kind]`. Draw the baseline's charts from their timers. Test a chart for a timer of the operator's own.

## Done when

A timer `stats.get` reports appears on the Charts page, and the flaiover tests pass.

## Notes
