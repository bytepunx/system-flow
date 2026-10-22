---
id: S-0485
type: story
nature: improvement
title: The chart page shows the title and the chart first, then a collapsible explainer of what the chart shows and how it is measured, then the controls, then the table view
status: ready
parent: E-0028
owner: alex
created: 2026-10-21T20:55:38Z
updated: 2026-10-22T20:17:34Z
transitions:
  - to: ready
    at: 2026-10-22T20:17:34Z
    by: alex
tags: [dashboard]
topics: [analysis]
touches: [flaiover/src/lib/viz/explain.ts, flaiover/src/lib/viz/explain.test.ts, flaiover/src/lib/chartexplainer.svelte.ts, flaiover/src/lib/chartexplainer.svelte.test.ts, flaiover/src/lib/components/ChartExplainer.svelte, flaiover/src/lib/components/ChartExplainer.svelte.test.ts, design/components/flaiover/src/lib/components/ChartExplainer.examples.json, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [S-0484]
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
      seconds: 41887
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 66
          output: 1031
          cache_read: 9622405
          cache_write: 12689
          cost: 2.1276
forecast:
  duration: 1h
  delivery: 2026-10-22T17:05:00Z
  basis: "flai forecast's 24m at size 16 (5 criteria, 11 touches), 89 s a unit over 86 done improvement stories on claude-opus-5-5 in the large band, raised to 1h because explain.ts counts one touch but holds three paragraphs written against metrics.md for each of the 28 chart kinds; delivery follows S-0484, which it waits for, at 2026-10-22T15:56Z, plus 1h at the pull order's cycle factor of 1.15."
  by: planner-E-0028
  at: 2026-10-21T21:03:27Z
finalized:
  by: alex
  at: 2026-10-22T20:16:59Z
---
# S-0485 The chart page shows the title and the chart first, then a collapsible explainer of what the chart shows and how it is measured, then the controls, then the table view

## Goal

Each chart's page reads top down: the chart's title and the chart itself first, with the figures the chart states beneath it; then a collapsible explainer that says what the chart shows, how each figure is measured, from which timestamps and fields of `flai stats`, and how to read it, linking the definition in `metrics.md` on the dashboard's Docs page; then the controls, the window and whatever filters the chart offers, with the summary strip; and last the expandable table view. Today the controls come first and the title and chart after them, and the chart's explanation is spread over a few conditional notes. The second half of E-0028's outcome.

## Acceptance criteria

- [ ] On `/charts/<kind>` the content is in this order: the `h1` title; the notices that replace or qualify the chart, as today; the chart, with the share chart of Time in State and the longest waits of Agent Waiting kept with it; the figures the chart states, `strategic-ratio`, `coordination-shares`, `cod-saving`, `cod-totals`, `cod-left-out`, and `cod-without-value`, as today; the explainer; the controls, the window and the filters the chart offers with the summary strip, unchanged in what they offer and do; and the table view last. A test in `charts.svelte.test.ts` asserts the order by the elements' positions in the DOM.
- [ ] `flaiover/src/lib/viz/explain.ts` gives every kind in `KINDS` an explainer with three parts, `what`, how the chart is to be read at a glance, `how`, how each figure it draws is measured, naming the transitions, timestamps, and `flai stats --json` fields `metrics.md` defines, such as cycle time from the first move to in progress to the move to done and the p50 and p85 as nearest-rank percentiles, and `read`, what a shape of the chart means, with a `metrics` anchor into `design/system/metrics.md`. The static notes the page shows under a chart today, `spend-note`, the two estimated-in-part asterisks, the explanation in `strategic-note`, `coordination-note`, and `queue-note`, move into the explainers and are drawn no more under the chart. `explain.test.ts` asserts every kind has non-empty `what` and `how`, and that each `metrics` anchor is a heading of `design/system/metrics.md`, read from the repository.
- [ ] `ChartExplainer.svelte` draws the explainer as a native `<details>` with `data-testid="explainer"`, its summary "About this chart", and the three parts as short paragraphs with the link to `/docs/design/system/metrics.md#<anchor>` through `resolve` and `withProject`. Its text may name the report's type and the bucket chosen, as the notes it replaces do. It is open on the first visit, and opening or closing it is kept per browser under `flaiover-chart-explainer`, one choice for every chart (`$lib/chartexplainer.svelte.ts`, as `chartwindow.svelte.ts` keeps the window). `design/components/flaiover/src/lib/components/ChartExplainer.examples.json` gives an example per group of charts, and `flai check` reports no `components.examples`.
- [ ] `charts.svelte.test.ts`, `explain.test.ts`, `chartexplainer.svelte.test.ts`, and `ChartExplainer.svelte.test.ts` cover the criteria above, the tests of the moved notes changed to look in the explainer; `flai test` on the changed paths passes and `flai check --strict` is clean.
- [ ] `docs/users/flaiover.md` § Charts describes the page's order and the explainer, and `design/system/flaiover-dashboard.md` § Views says where the explainers live and what they hold; `metrics.md` stays the source of the definitions and is not changed.

## Tasks
- T-2020 viz/explain.ts gives every chart kind an explainer, what it shows, how each figure is measured, how to read it, and its metrics.md anchor
- T-2021 chartexplainer.svelte.ts keeps whether the explainer is open per browser, open on the first visit, one choice for every chart
- T-2023 ChartExplainer.svelte draws a chart's explainer as a details element that remembers whether it is open, with its examples file
- T-2025 The chart page puts the title and the chart first, the explainer after the chart's figures, the controls beneath it, and the table view last
- T-2027 The user docs and the dashboard design describe the chart page's order and the explainer

## Notes

- Drafted by the planner for E-0028. Waits for S-0484, which removes the page's own chart list from `/charts/[kind]/+page.svelte` and changes its test and the same two documents; this story reorders what S-0484 leaves.
- The explainer's text is written from `docs/users/flaiover.md` § Charts and `design/system/metrics.md`, which define every figure; where the two disagree, `metrics.md` wins and the user docs are corrected in this story.
- The chart's own figures, the saving under CoD by Order or the ratio under Strategic Cost, stay under the chart rather than in the explainer: they are results of the window chosen, not explanation.

### Planning

- **Plan thread.** TH-0591 names the two stories, their order, these tasks and their layers, and the assumptions: the explainer open on the first visit and remembered per browser, the figures a chart states kept under the chart, and the summary strip kept with the controls.
- **Layers.** T-2020 (the explainers) and T-2021 (the open state) first, sharing no path; T-2023 (`ChartExplainer.svelte` and its examples) after both; T-2025 (the page's order) after T-2023, since it mounts the component; T-2027 (the docs) after T-2025, which settles what the page shows.
- **Touches.** 11, all files, every one on a task.
  - **From the code layout:** the new `viz/explain.ts` and its test, beside `viz/charts.ts` where the titles are; the new `chartexplainer.svelte.ts` and its test, after `chartwindow.svelte.ts`; the new `ChartExplainer.svelte` and its test under `lib/components`, with its `design/components/…/ChartExplainer.examples.json`, which `flai check` wants for each component (`components.examples`); `charts/[kind]/+page.svelte` and `charts.svelte.test.ts`, where the order and the notes are.
  - **From design:** `docs/users/flaiover.md` § Charts and `design/system/flaiover-dashboard.md` § Views, which describe the page and its notes. `design/system/metrics.md` is read for every explainer and deliberately not touched: it stays the definition the explainers link.
  - **Co-change, not added:** `flai touches suggest S-0485` lists `viz/charts.ts` and `viz/charts.test.ts` at 8% and otherwise paths of unrelated stories (`flai-cli.md` 42%, `docs/users/flai.md` 41%, `docs/operators/index.md` 30%); the explainers are a module of their own, so `charts.ts` is unchanged.
  - No folder touch.
- **Forecast.** flai's 24m at size 16 (5 criteria, 11 touches), 89 s a unit over 86 done improvement stories on claude-opus-5-5 in the large band, raised to 1h: `explain.ts` counts one touch but holds three paragraphs written against `metrics.md` for each of the 28 chart kinds, the bulk of the story's work. Delivery follows S-0484, which this story waits for, at 2026-10-22T15:56Z, plus 1h at the pull order's cycle factor of 1.15: 2026-10-22T17:05Z. flai's own play-out with 24m gave 2026-10-22T16:12Z.
- **Cost of delay.** None yet: E-0028 carries no inputs, and `flai cod S-0485` refuses without them. TH-0590 asks the operator for them; the value is written here once the input is on the epic.
