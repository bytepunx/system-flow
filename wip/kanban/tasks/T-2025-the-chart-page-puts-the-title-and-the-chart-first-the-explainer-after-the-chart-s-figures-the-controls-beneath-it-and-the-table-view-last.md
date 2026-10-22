---
id: T-2025
type: task
nature: improvement
title: The chart page puts the title and the chart first, the explainer after the chart's figures, the controls beneath it, and the table view last
status: backlog
parent: S-0485
owner: alex
created: 2026-10-21T20:58:40Z
updated: 2026-10-21T20:58:40Z
transitions: []
stream: S-0485
tags: [dashboard]
touches: ["flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-2023]
---
# T-2025 The chart page puts the title and the chart first, the explainer after the chart's figures, the controls beneath it, and the table view last

## Work

Reorder the template of `flaiover/src/routes/charts/[kind]/+page.svelte`: the `h1` title first, then the notices that replace or qualify the chart as they are, then `<Chart>`, with Time in State's share chart and Agent Waiting's longest waits kept after it, then the figures the chart states (`strategic-ratio`, `coordination-shares`, `cod-saving`, `cod-totals`, `cod-left-out`, `cod-without-value`), then `<ChartExplainer {kind} {type} {bucket} />`, then the controls `div` moved down with everything it holds, the window, per, type, epic, nature, model, and the summary strip, and the `<details>` table view last. Drop the static notes the explainers now hold: `spend-note`, the two estimated-in-part asterisk paragraphs, the explanatory sentences of `strategic-note` (keep the `strategic-ratio` span and its no-ratio sentences as the figure under the chart), `coordination-note`, and `queue-note`. The title must show while the report loads too, so the controls are reachable before the chart draws.

Third layer: waits for T-2023, the component it mounts.

## Done when

`charts.svelte.test.ts` asserts the order by `compareDocumentPosition` for a flow chart, a spend chart, Strategic Cost, and CoD by Order; the tests that looked for the dropped notes look in the explainer instead; the controls still change the report as before. `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. S-0484 has removed the page's sidebar before this task runs; if it has not been accepted, the story waits in ready on `after`.
