---
id: T-0096
type: task
nature: feature
title: "Chart routes: cycle time, burn-up, cumulative flow, time in state, throughput, aging, estimates, with a shared filter bar"
status: done
parent: S-0014
owner: alex
created: 2026-09-17T04:50:13Z
updated: 2026-09-17T05:02:58Z
transitions:
  - to: ready
    at: 2026-09-17T05:02:58Z
    by: agent
  - to: in-progress
    at: 2026-09-17T05:02:58Z
    by: agent
  - to: done
    at: 2026-09-17T05:02:58Z
    by: agent
stream: S-0014
tags: [dashboard, charts]
---

# T-0096 Chart routes: cycle time, burn-up, cumulative flow, time in state, throughput, aging, estimates, with a shared filter bar

## Work
/charts/[kind] with a kind switcher, a filter bar (window, type, epic where it applies), summary strip, the chart, the state-share bar under time in state, and a table view; /charts redirects to cycle time; nav gains Charts.

## Done when
Every chart in metrics.md has a route and renders against this repository in light and dark.

## Notes
