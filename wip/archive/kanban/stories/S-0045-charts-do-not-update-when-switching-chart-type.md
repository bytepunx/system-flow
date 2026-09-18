---
id: S-0045
type: story
nature: remediation
title: Charts do not update when switching chart type
status: done
parent: E-0003
owner: alex
created: 2026-09-18T16:15:39Z
updated: 2026-09-18T16:30:53Z
transitions:
  - to: ready
    at: 2026-09-18T16:19:25Z
    by: alex
  - to: in-progress
    at: 2026-09-18T16:19:25Z
    by: alex
  - to: review
    at: 2026-09-18T16:25:57Z
    by: alex
  - to: done
    at: 2026-09-18T16:30:53Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/lib/components/Chart.svelte, flaiover/src/routes/charts]
---

# S-0045 Charts do not update when switching chart type

## Goal
Clicking between chart types on the charts page redraws the chart immediately. Today the tabs, title, and table view change but the plot keeps showing the previous chart until the page is refreshed by hand.

## Acceptance criteria
- [x] Switching between every pair of chart types redraws the plot without a reload, including the height change for aging work in progress
- [x] Changing the window, type, or epic filter, a theme change, and a live update over the event stream also redraw the plot
- [x] A component test mounts the chart, changes its option, and asserts the chart instance received the new option; it fails on the current code
- [x] Checked in a browser on the served build across all seven chart types in both themes
- [x] I-0014 is closed with a reference to this story
- [x] An empty selection draws an empty chart instead of throwing: the dashboard tolerates `null` lists from older flai builds, and `flai stats --json` emits `[]` for empty lists

## Tasks
- T-0134 Confirm the cause with a failing component test that changes the chart option after mount
- T-0135 Fix the redraw in Chart.svelte so option, height, and theme changes reach the chart instance
- T-0136 Browser check across the seven chart types in both themes; close I-0014

## Notes
- Reported by the operator on 2026-09-18 (I-0014).
- Suspected cause, to confirm first: in `src/lib/components/Chart.svelte` the effect runs `chart?.setOption(option, true)` while `chart` is still null on its first run, so the optional chain short-circuits before `option` is read, the effect registers no dependency on it, and never runs again. The first draw comes from `onMount`, which is why a refresh fixes it. Filters and live updates would be affected the same way.
- Found while verifying the fix: with type set to task the aging list came back as `null` from `flai stats`, and the aging chart threw on `.filter`, freezing the page. The redraw defect had hidden it, because the chart never redrew. Fixed on both sides in this story.
