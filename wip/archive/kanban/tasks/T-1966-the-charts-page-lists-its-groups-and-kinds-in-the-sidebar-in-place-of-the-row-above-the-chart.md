---
id: T-1966
type: task
nature: improvement
title: The charts page lists its groups and kinds in the sidebar in place of the row above the chart
status: done
parent: S-0461
owner: alex
created: 2026-10-21T07:04:18Z
updated: 2026-10-21T12:17:28Z
transitions:
  - to: ready
    at: 2026-10-21T12:12:05Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T12:12:06Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T12:17:28Z
    by: agent-S-0461
stream: S-0461
tags: [dashboard]
touches: ["flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts"]
after: [T-1984]
usage:
  source: log
  seconds: 322
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 25
      output: 11648
      cache_read: 1488017
      cache_write: 51606
      cost: 0.8363
---
# T-1966 The charts page lists its groups and kinds in the sidebar in place of the row above the chart

## Work

- In `flaiover/src/routes/charts/[kind]/+page.svelte`, replace the group-and-kind `<nav>` with a `PageSidebar`. Each group is an entry, and its kinds are child pane entries linking `/charts/<kind>`, with the current kind marked.
- Keep the chart window control with the chart.
- Update `charts.svelte.test.ts` to find the kind links in the sidebar, with `aria-current` on the kind shown.
- Waits for T-1984, which lets a sidebar entry link to another page. It runs beside the review and documents tasks.

## Done when

`charts.svelte.test.ts` passes, and every chart kind is reached from the sidebar.

## Notes
