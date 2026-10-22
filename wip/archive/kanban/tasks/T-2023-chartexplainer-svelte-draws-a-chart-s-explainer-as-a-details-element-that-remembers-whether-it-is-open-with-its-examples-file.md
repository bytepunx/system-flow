---
id: T-2023
type: task
nature: improvement
title: ChartExplainer.svelte draws a chart's explainer as a details element that remembers whether it is open, with its examples file
status: done
parent: S-0485
owner: alex
created: 2026-10-21T20:58:02Z
updated: 2026-10-22T21:02:12Z
transitions:
  - to: ready
    at: 2026-10-22T21:00:28Z
    by: agent-S-0485
  - to: in-progress
    at: 2026-10-22T21:00:30Z
    by: agent-S-0485
  - to: done
    at: 2026-10-22T21:02:12Z
    by: agent-S-0485
stream: S-0485
tags: [dashboard]
touches: [flaiover/src/lib/components/ChartExplainer.svelte, flaiover/src/lib/components/ChartExplainer.svelte.test.ts, design/components/flaiover/src/lib/components/ChartExplainer.examples.json]
after: [T-2020, T-2021]
usage:
  source: log
  seconds: 102
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 15638
      cache_read: 2475534
      cache_write: 94398
      cost: 1.3964
---
# T-2023 ChartExplainer.svelte draws a chart's explainer as a details element that remembers whether it is open, with its examples file

## Work

Write `flaiover/src/lib/components/ChartExplainer.svelte` with props `kind: Kind`, `type: string`, `bucket: BucketSize`, and an optional `state` defaulting to the shared explainer state, for tests. It draws a native `<details data-testid="explainer">`, `open` from the state, whose `<summary>` reads "About this chart", and inside it the explainer's `what`, `how`, and `read` from `explain(kind, { type, bucket })` as three short paragraphs in the muted small text the notes use today, with a last line linking "Defined in metrics.md" to `/docs/design/system/metrics.md#<metrics>` through `resolve` and `withProject(…, projectState.current)`. On the element's `toggle` event it calls `state.set(open)`. Write `design/components/flaiover/src/lib/components/ChartExplainer.examples.json` with one example per chart group, cycle time, $ / Day with a bucket, Forecast Accuracy, and Strategic Cost.

Second layer: waits for T-2020 for the explainers and T-2021 for the state.

## Done when

`ChartExplainer.svelte.test.ts` renders the component open and closed from a given state, asserts the three paragraphs and the link's href for a kind, and asserts a toggle writes the state. `flai check` reports no `components.examples` on the new file, and `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. The page keeps its notices and figures; the component holds only the explanation.
