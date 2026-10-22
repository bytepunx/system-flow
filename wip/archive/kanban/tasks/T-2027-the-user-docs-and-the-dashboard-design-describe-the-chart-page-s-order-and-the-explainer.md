---
id: T-2027
type: task
nature: improvement
title: The user docs and the dashboard design describe the chart page's order and the explainer
status: done
parent: S-0485
owner: alex
created: 2026-10-21T20:59:15Z
updated: 2026-10-22T21:09:33Z
transitions:
  - to: ready
    at: 2026-10-22T21:07:18Z
    by: agent-S-0485
  - to: in-progress
    at: 2026-10-22T21:07:19Z
    by: agent-S-0485
  - to: done
    at: 2026-10-22T21:09:33Z
    by: agent-S-0485
stream: S-0485
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-2025]
usage:
  source: log
  seconds: 134
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 14991
      cache_read: 2373058
      cache_write: 90490
      cost: 1.3386
---
# T-2027 The user docs and the dashboard design describe the chart page's order and the explainer

## Work

In `docs/users/flaiover.md` § Charts, say the page's order: the title and the chart, the figures it states, **About this chart**, which says what the chart shows, how its figures are measured, and how to read it, links the definition in `metrics.md`, and stays open or closed as you last left it in this browser, then the window and filters, then the table view. Correct any sentence the explainers task found at odds with `metrics.md`, and drop the notes the page no longer shows under a chart from the chart table's rows where they are described. In `design/system/flaiover-dashboard.md` § Views, add a paragraph on the chart page's order, `ChartExplainer.svelte`, `viz/explain.ts`, and `chartexplainer.svelte.ts` with its key, and amend the `/charts/*` rows that name the notes now in the explainer (`spend-note`, `strategic-note`, `coordination-note`, `queue-note`). Bump each document's `updated`.

Fourth layer: waits for T-2025, which settles what the page shows.

## Done when

Both documents describe the page as T-2025 built it, the markdown lint and `flai check --strict` pass on them, and `metrics.md` is unchanged.

## Notes

Drafted by the planner for E-0028. S-0484 changes the same two documents for the rail; this task writes after its acceptance, as the story's `after` ensures.
