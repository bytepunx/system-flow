---
id: T-0095
type: task
nature: feature
title: Chart component on ECharts core with light and dark themes and a shared palette
status: done
parent: S-0014
owner: alex
created: 2026-09-17T04:50:13Z
updated: 2026-09-17T05:02:57Z
transitions:
  - to: ready
    at: 2026-09-17T05:02:57Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:02:57Z
    by: agent
  - to: done
    at: 2026-09-17T05:02:57Z
    by: agent
stream: S-0014
tags: [dashboard, charts]
---

# T-0095 Chart component on ECharts core with light and dark themes and a shared palette

## Work
src/lib/viz/palette.ts carries the validated categorical palette for light and dark with fixed slots for states and natures; Chart.svelte initialises ECharts core with only scatter, line, bar, grid, tooltip, legend, mark line, and the canvas renderer, follows option changes, resizes with its container; charts.ts builds one option per chart from the flai report as pure functions (one axis, thin marks, legend for two or more series, tooltips, reference lines for p50 and p85).

## Done when
Option builders unit-tested; palette validated with the dataviz script in both modes.

## Notes
