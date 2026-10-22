---
id: T-2020
type: task
nature: improvement
title: viz/explain.ts gives every chart kind an explainer, what it shows, how each figure is measured, how to read it, and its metrics.md anchor
status: backlog
parent: S-0485
owner: alex
created: 2026-10-21T20:57:25Z
updated: 2026-10-21T20:57:25Z
transitions: []
stream: S-0485
tags: [dashboard]
touches: [flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts]
---
# T-2020 viz/explain.ts gives every chart kind an explainer, what it shows, how each figure is measured, how to read it, and its metrics.md anchor

## Work

Write `flaiover/src/lib/viz/explain.ts`: `export type Explainer = { what: string; how: string; read: string; metrics: string }` and `EXPLAINERS: Record<Kind, Explainer>`, with a function `explain(kind, { type, bucket })` that fills the type and the bucket into the text where a chart names them, as the notes it replaces do (`spend-note` names the bucket, `strategic-note` the type). Write each explainer from `docs/users/flaiover.md` § Charts and `design/system/metrics.md`: `what` in one or two sentences; `how` naming the transitions, timestamps, and `flai stats --json` fields the figures come from, such as cycle time as the first move to in progress to the move to done, the p50 and p85 as nearest-rank percentiles, a spend bucket holding the items that entered done in it, a hold's one reason; `read` saying what a shape means, as the planning and strategic subsections of the user docs do; `metrics` the anchor of the `metrics.md` heading that defines the chart, such as `per-item-durations`, `spend-over-time-s-0163`, `cost-of-delay`, `waiting`, `claims-and-touches`, `coordination`, `merge-queue-s-0404`, `strategic-use-per-day`. Fold in the text of `spend-note`, the two estimated-in-part asterisks, the explanatory sentences of `strategic-note`, `coordination-note`, and `queue-note`. Where the user docs and `metrics.md` disagree, follow `metrics.md` and note the difference for the docs task.

First layer: waits for nothing. Runs beside the explainer state task, with which it shares no path.

## Done when

`explain.test.ts` asserts that every kind of `KINDS` has an explainer whose `what`, `how`, and `metrics` are non-empty, that `explain` fills the type and the bucket into the spend and strategic texts, and that each `metrics` anchor matches a heading of `design/system/metrics.md` read from the repository, slugged as markdown-it-anchor slugs it. `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. The explainers are the one place the dashboard explains a chart; `metrics.md` stays the definition and is linked, not copied at length.
