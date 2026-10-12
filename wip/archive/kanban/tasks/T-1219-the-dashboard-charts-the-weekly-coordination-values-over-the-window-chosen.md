---
id: T-1219
type: task
nature: feature
title: The dashboard charts the weekly coordination values over the window chosen
status: done
parent: S-0337
owner: alex
created: 2026-10-07T20:17:32Z
updated: 2026-10-12T16:14:32Z
transitions:
  - to: ready
    at: 2026-10-12T15:35:43Z
    by: agent-S-0337
  - to: in-progress
    at: 2026-10-12T16:08:23Z
    by: agent-S-0337
  - to: done
    at: 2026-10-12T16:14:32Z
    by: agent-S-0337
stream: S-0337
tags: [flaiover]
touches: [flaiover/src/lib/viz/charts.ts, flaiover/src/lib/viz/charts.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-1218]
usage:
  source: log
  seconds: 369
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 19341
      cache_read: 3599285
      cache_write: 117273
      cost: 1.8881
---
# T-1219 The dashboard charts the weekly coordination values over the window chosen

## Work

Chart what T-1218 reports. It waits for T-1218, whose `coordination` it reads.

- A coordination chart kind in `charts.ts`: conflicts found early against conflicts met late per week, and conversations by how they ended, spanning the window chosen (ADR-0054).
- Render it on `src/routes/charts/[kind]/+page.svelte`, with the stories started on a share and the hold time saved beside it.

## Done when

- Tests cover the series built from a fixture, an empty window, and the page rendering the kind.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
