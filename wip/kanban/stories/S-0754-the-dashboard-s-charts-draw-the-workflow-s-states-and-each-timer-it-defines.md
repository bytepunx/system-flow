---
id: S-0754
type: story
nature: feature
title: The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:59Z
updated: 2026-11-12T16:01:59Z
transitions: []
tags: [dashboard]
topics: [dashboard, analysis]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/server/stats.ts, flaiover/src/lib/server/stats.test.ts, flaiover/src/routes/charts/+page.ts, design/system/metrics.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0752, S-0753]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0754 The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time

## Goal

Every chart follows the workflow. Time in state, cumulative flow, and the cost of delay per column stack the workflow's states in its order and colour them by kind, read from flai, and each timer the workflow defines, the baseline's and the operator's own, is charted as lead and cycle time are, so a timer added to the workflow appears on the Charts page with no change to flaiover.

## Acceptance criteria

- [ ] `flaiover/src/lib/viz/charts.ts` holds no list of states: `STATES` and `COD_COLUMNS` give way to the columns `board.get` answers, in the workflow's order and by kind
- [ ] Each timer `stats.get` reports gets a chart in the catalogue, with its distribution over the window as cycle time has, and the baseline's charts are drawn from their timers
- [ ] Tests draw time in state and cumulative flow for a workflow with a state added between `in-progress` and `review`, and a chart for a timer of the operator's own
- [ ] With the baseline, every chart draws as today and every existing flaiover test passes
- [ ] `design/system/metrics.md` § Charts has a row for a workflow's own timer, as ADR-0219 asks, and `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say the charts follow the workflow

## Tasks

## Notes

- `charts.ts` line 497 (`COD_COLUMNS`) and line 616 (`STATES`, five states with no `cancelled`) are the lists today. ADR-0219 asks `metrics.md` § Charts for a row under its catalogue title for every chart.
