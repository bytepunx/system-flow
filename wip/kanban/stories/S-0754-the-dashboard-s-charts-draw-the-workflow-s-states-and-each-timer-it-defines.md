---
id: S-0754
type: story
nature: feature
title: The dashboard's charts draw the workflow's states and each timer it defines from what flai gives, and an operator's own timer gets a chart beside lead and cycle time
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:01:59Z
updated: 2026-11-13T07:23:13Z
transitions: []
tags: [dashboard]
topics: [dashboard, analysis]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, flaiover/src/lib/server/stats.ts, flaiover/src/lib/server/stats.test.ts, flaiover/src/routes/charts/+page.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", design/system/metrics.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0752, S-0753]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 15
          cache_read: 121381
          cache_write: 20507
          cost: 0.0405
draft: true
forecast:
  duration: 30m
  delivery: 2026-11-13T19:43:00Z
  basis: "Its own forecast of 30m; 28th in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663, S-0730, S-0733, S-0734, S-0735, S-0736, S-0737, S-0738, S-0739, S-0740, S-0741, S-0742, S-0743, S-0744, S-0745, S-0746, S-0747, S-0748, S-0749, S-0750, S-0751, S-0752 and S-0753."
  by: flai
  at: 2026-11-13T07:23:13Z
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

- T-3245 The charts stack and colour the workflow's states from the columns board.get answers, in place of STATES and COD_COLUMNS
- T-3251 Each timer stats.get reports gets a chart in the catalogue with its distribution over the window
- T-3254 metrics.md § Charts gains a row for a workflow's own timer, and the dashboard's design and guide say the charts follow the workflow

## Notes

- `charts.ts` line 497 (`COD_COLUMNS`) and line 616 (`STATES`, five states with no `cancelled`) are the lists today. ADR-0219 asks `metrics.md` § Charts for a row under its catalogue title for every chart.

### Planning

- Touches, all files:
  - `flaiover/src/lib/viz/charts.ts`, `charts.test.ts`, `flaiover/src/lib/server/stats.ts`, `stats.test.ts`, `routes/charts/+page.ts`: layout.
  - Co-change: `routes/charts/[kind]/+page.svelte` and `charts.svelte.test.ts` (6% to 7% in `flai touches suggest`), where a chart kind is drawn.
  - `design/system/metrics.md`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`: design (criterion 5).
  - `flaiover/src/lib/viz/palette.ts` is S-0758's, which colours states by kind for every view.
  - Left out from `flai touches suggest`: `docs/users/flai.md`, `flai-cli.md`, and `docs/operators/index.md` (25% to 38%), flai's side.
- `after`: S-0752 and S-0753. The timer charts read the section S-0752 adds to `stats.get`, and the stacked charts read the columns and kinds S-0753 adds to `board.get`.
- Tasks: T-3245 (stacked charts), then T-3251 (timer charts), since both change `charts.ts`, then T-3254 (the documents), which describe T-3251's chart.
- Forecast 30m, up from the 24m `flai forecast` gave for the first draft's size 13; size 15 now gives 27m, and a new chart kind on the catalogue has taken more than its size in earlier dashboard stories.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
