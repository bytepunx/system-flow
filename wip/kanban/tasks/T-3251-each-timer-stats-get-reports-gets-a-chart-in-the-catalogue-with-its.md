---
id: T-3251
type: task
nature: feature
title: Each timer stats.get reports gets a chart in the catalogue with its distribution over the window
status: backlog
parent: S-0754
owner: alex
created: 2026-11-12T16:09:40Z
updated: 2026-11-12T16:09:40Z
transitions: []
stream: S-0754
tags: [dashboard]
touches: [flaiover/src/lib/server/stats.ts, flaiover/src/lib/server/stats.test.ts, flaiover/src/routes/charts/+page.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", flaiover/src/lib/viz/charts.ts]
after: [T-3245]
---
# T-3251 Each timer stats.get reports gets a chart in the catalogue with its distribution over the window

## Work

Waits for T-3245: both change `charts.ts`, and the timer charts use its columns. Read the timers section of `stats.get` in `stats.ts`, add a chart kind per timer to the catalogue in `routes/charts/+page.ts`, and draw each as cycle time is drawn, its distribution over the window, in `charts/[kind]`. Draw the baseline's charts from their timers. Test a chart for a timer of the operator's own.

## Done when

A timer `stats.get` reports appears on the Charts page, and the flaiover tests pass.

## Notes
