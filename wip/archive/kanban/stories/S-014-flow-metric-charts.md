---
id: S-014
type: story
nature: feature
title: Flow metric charts
status: done
parent: E-003
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T05:29:29Z
transitions:
  - to: ready
    at: 2026-09-17T04:50:14Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:50:14Z
    by: agent
  - to: review
    at: 2026-09-17T05:02:59Z
    by: agent
  - to: done
    at: 2026-09-17T05:29:29Z
    by: alex
tags: []
---

# S-014 Flow metric charts

## Goal
Cycle time scatter, burn-up, cumulative flow, time in state, throughput, aging WIP, and estimate versus actual, with a shared filter bar.

## Acceptance criteria
- [x] Each chart in metrics.md has a route
- [x] Numbers come from `flai stats --json` invoked by the server (ADR-0016), so they match by construction; the client charts them
- [x] Charts work in light and dark theme

## Tasks
- T-094 GET /api/stats delegating to flai stats --json with a change-invalidated cache
- T-095 Chart component on ECharts core with light and dark themes and a shared palette
- T-096 Chart routes: cycle time, burn-up, cumulative flow, time in state, throughput, aging, estimates, with a shared filter bar
- T-097 Tests, docs, tech pins

## Notes
- Numbers come from flai stats --json (ADR-0016); the visual check found and fixed an SSE teardown crash in the server.
