---
id: T-2021
type: task
nature: improvement
title: chartexplainer.svelte.ts keeps whether the explainer is open per browser, open on the first visit, one choice for every chart
status: done
parent: S-0485
owner: alex
created: 2026-10-21T20:57:32Z
updated: 2026-10-22T20:51:52Z
transitions:
  - to: ready
    at: 2026-10-22T20:50:57Z
    by: agent-S-0485
  - to: in-progress
    at: 2026-10-22T20:51:00Z
    by: agent-S-0485
  - to: done
    at: 2026-10-22T20:51:52Z
    by: agent-S-0485
stream: S-0485
tags: [dashboard]
touches: [flaiover/src/lib/chartexplainer.svelte.ts, flaiover/src/lib/chartexplainer.svelte.test.ts]
usage:
  source: log
  seconds: 52
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 2214
      cache_read: 350481
      cache_write: 13365
      cost: 0.1977
---
# T-2021 chartexplainer.svelte.ts keeps whether the explainer is open per browser, open on the first visit, one choice for every chart

## Work

Write `flaiover/src/lib/chartexplainer.svelte.ts` after the pattern of `chartwindow.svelte.ts`: a state object with `open`, read from `localStorage` under `flaiover-chart-explainer` and `true` when nothing is stored or the stored value is not a boolean, and `set(open)`, which stores the choice and swallows a storage error so a private window keeps the choice for the page.

First layer: waits for nothing. Runs beside the explainers task, with which it shares no path.

## Done when

`chartexplainer.svelte.test.ts` covers the default, a stored `false`, a stored non-boolean, and `set` writing the key. `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. One choice for every chart: an operator who has read the explainers closes them once.
