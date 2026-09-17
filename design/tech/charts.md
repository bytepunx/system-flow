---
title: Charts
updated: 2026-09-15
status: active
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

Colours come from the validated reference palette in `flaiover/src/lib/viz/palette.ts` (eight categorical slots for light and for dark, both modes checked with the dataviz validator on 2026-09-17). Workflow states and natures have fixed slots so a state keeps its colour across charts. Rules applied in every builder (`src/lib/viz/charts.ts`): one y-axis per chart, thin marks (2px lines, bars capped at 24px, a surface-coloured seam between stacked segments), a legend whenever two or more series are shown, a tooltip on every mark, p50 and p85 as dashed reference lines, scatter colours at most three groups (the all-pairs rule) and folds the rest into "other", and every chart page offers a table view.

## Considered

- LayerChart: Svelte-native and pleasant, but younger and its Svelte 5 support was still moving when this was decided. Revisit if ECharts styling fights Tailwind too much.
- Chart.js: fine for bars and lines, weaker for the CFD and annotated scatter.
- D3 directly: maximum control, too much code for the initial delivery.
