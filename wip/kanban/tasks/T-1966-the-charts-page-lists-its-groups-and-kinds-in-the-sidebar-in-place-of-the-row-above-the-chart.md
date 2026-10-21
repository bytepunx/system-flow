---
id: T-1966
type: task
nature: improvement
title: The charts page lists its groups and kinds in the sidebar in place of the row above the chart
status: backlog
parent: S-0461
owner: alex
created: 2026-10-21T07:04:18Z
updated: 2026-10-21T07:04:18Z
transitions: []
stream: S-0461
tags: [dashboard]
touches: ["flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
---
# T-1966 The charts page lists its groups and kinds in the sidebar in place of the row above the chart

## Work

- In `flaiover/src/routes/charts/[kind]/+page.svelte`, replace the group-and-kind `<nav>` with a `PageSidebar`. Each group is an entry, and its kinds are child pane entries linking `/charts/<kind>`, with the current kind marked.
- Keep the chart window control with the chart.
- Update `charts.svelte.test.ts` to find the kind links in the sidebar, with `aria-current` on the kind shown.
- Waits for nothing in this story. It runs beside the review and documents tasks.

## Done when

`charts.svelte.test.ts` passes, and every chart kind is reached from the sidebar.

## Notes
