---
title: Charts
updated: 2026-09-15
status: active
---

# Charts

| | |
|-|-|
| Library | Apache ECharts 6, imported per chart type to keep bundle size down |
| Used in | `flaiover` |

## Why

- Every chart in [metrics.md](../system/metrics.md) is covered natively: scatter with reference lines, stacked area for the cumulative flow diagram, stacked bars, line with forecast segment.
- Mature, framework-agnostic, wraps in a small Svelte action with resize handling.

## Considered

- LayerChart: Svelte-native and pleasant, but younger and its Svelte 5 support was still moving when this was decided. Revisit if ECharts styling fights Tailwind too much.
- Chart.js: fine for bars and lines, weaker for the CFD and annotated scatter.
- D3 directly: maximum control, too much code for the initial delivery.
