---
id: T-2019
type: task
nature: improvement
title: The nav rail's state keeps whether Charts is expanded beside the groups, under the same localStorage key, and expands it on arrival at a chart
status: done
parent: S-0484
owner: alex
created: 2026-10-21T20:57:15Z
updated: 2026-10-22T20:23:19Z
transitions:
  - to: ready
    at: 2026-10-22T20:19:30Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:19:34Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:23:19Z
    by: agent-S-0484
stream: S-0484
tags: [dashboard]
touches: [flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/navrail.svelte.test.ts]
usage:
  source: log
  seconds: 225
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 10
      output: 3100
      cache_read: 485181
      cache_write: 17395
      cost: 0.265
---
# T-2019 The nav rail's state keeps whether Charts is expanded beside the groups, under the same localStorage key, and expands it on arrival at a chart

## Work

Widen `NavRailState` in `flaiover/src/lib/navrail.svelte.ts` so that an expandable key is a `GroupKey` or `'charts'`: `Chosen`, `parseChosen`, `isExpanded`, `toggle`, and `expand` take the wider key, and `parseChosen` keeps a boolean under `charts` and still drops anything else, so a value stored by an older dashboard parses. `toggle` stores the choice under `flaiover-nav-rail` as it does for a group; `expand` sets it for the page's life without storing it.

First layer: waits for nothing. Runs beside the catalogue task, with which it shares no path.

## Done when

`navrail.svelte.test.ts` covers `parseChosen` with `charts` true, false, and a non-boolean, `toggle('charts')` stored beside a group's choice, and `expand('charts')` not stored. `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. The key stays `flaiover-nav-rail` so the groups the operator already chose are kept.
