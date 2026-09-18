---
id: S-0045
type: story
nature: remediation
title: Charts do not update when switching chart type
status: backlog
parent: E-0003
owner: alex
created: 2026-09-18T16:15:39Z
updated: 2026-09-18T16:15:39Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/lib/components/Chart.svelte, flaiover/src/routes/charts]
---

# S-0045 Charts do not update when switching chart type

## Goal
Clicking between chart types on the charts page redraws the chart immediately. Today the tabs, title, and table view change but the plot keeps showing the previous chart until the page is refreshed by hand.

## Acceptance criteria
- [ ] Switching between every pair of chart types redraws the plot without a reload, including the height change for aging work in progress
- [ ] Changing the window, type, or epic filter, a theme change, and a live update over the event stream also redraw the plot
- [ ] A component test mounts the chart, changes its option, and asserts the chart instance received the new option; it fails on the current code
- [ ] Checked in a browser on the served build across all seven chart types in both themes
- [ ] I-0014 is closed with a reference to this story

## Tasks

## Notes
- Reported by the operator on 2026-09-18 (I-0014).
- Suspected cause, to confirm first: in `src/lib/components/Chart.svelte` the effect runs `chart?.setOption(option, true)` while `chart` is still null on its first run, so the optional chain short-circuits before `option` is read, the effect registers no dependency on it, and never runs again. The first draw comes from `onMount`, which is why a refresh fixes it. Filters and live updates would be affected the same way.
