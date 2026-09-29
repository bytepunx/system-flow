---
title: Charts
updated: 2026-09-29
status: active
topics: [dashboard]
---

# Charts

| | |
|-|-|
| Library | Apache ECharts 6.1, `echarts/core` with only scatter, line, bar, grid, tooltip, legend, mark line, and the canvas renderer, loaded lazily by `Chart.svelte` |
| Used in | `flaiover` |

## Why

- Every chart in [metrics.md](../system/metrics.md) is covered natively: scatter with reference lines, stacked area for the cumulative flow diagram, stacked bars, line with forecast segment.
- Mature, framework-agnostic, wraps in a small Svelte action with resize handling.

## Palette and rules

Colours come from the validated reference palette in `flaiover/src/lib/viz/palette.ts` (eight categorical slots for light and for dark, both modes checked with the dataviz validator on 2026-09-17). Workflow states, natures, model families (S-0143), and item types (S-0163) have fixed slots so a state, a nature, a model, or a type keeps its colour across charts. On the line charts of spend over time a model and an item type also have a mark each (circle, square, triangle, diamond), shown in the legend: the opus and fable slots are close for a reader who sees little red or green (the validator's all-pairs check gives 0.6 in light and 2.2 in dark), so colour is not the only thing that tells two models apart. Those charts draw time in UTC, with ticks no finer than a bucket and the axis a bucket wider than the data on each side. Rules applied in every builder (`src/lib/viz/charts.ts`): one y-axis per chart, thin marks (2px lines, bars capped at 24px, a surface-coloured seam between stacked segments), a legend whenever two or more series are shown, a tooltip on every mark, p50 and p85 as dashed reference lines, scatter colours at most three groups (the all-pairs rule) and folds the rest into "other", and every chart page offers a table view.

## Considered

- LayerChart: Svelte-native and pleasant, but younger and its Svelte 5 support was still moving when this was decided. Revisit if ECharts styling fights Tailwind too much.
- Chart.js: fine for bars and lines, weaker for the CFD and annotated scatter.
- D3 directly: maximum control, too much code for the initial delivery.
