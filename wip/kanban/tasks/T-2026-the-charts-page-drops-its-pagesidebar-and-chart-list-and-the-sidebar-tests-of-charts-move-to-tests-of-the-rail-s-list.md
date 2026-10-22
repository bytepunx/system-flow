---
id: T-2026
type: task
nature: improvement
title: The charts page drops its PageSidebar and chart list, and the sidebar tests of charts move to tests of the rail's list
status: done
parent: S-0484
owner: alex
created: 2026-10-21T20:59:05Z
updated: 2026-10-22T20:35:17Z
transitions:
  - to: ready
    at: 2026-10-22T20:19:32Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:31:22Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:35:17Z
    by: agent-S-0484
stream: S-0484
tags: [dashboard]
touches: ["flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", flaiover/e2e/page-sidebar.e2e.ts]
after: [T-2024]
usage:
  source: log
  seconds: 235
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 10072
      cache_read: 1576331
      cache_write: 56515
      cost: 0.8609
---
# T-2026 The charts page drops its PageSidebar and chart list, and the sidebar tests of charts move to tests of the rail's list

## Work

In `flaiover/src/routes/charts/[kind]/+page.svelte`, remove the `PageSidebar` wrapper, the local `groups` const, the `entries` derived, and the `SidebarEntry` and `PageSidebar` imports, so the controls, the title, the chart, its notes, and the table view are the page's content at full width beside the rail; nothing else in the page changes. In `charts.svelte.test.ts`, remove the S-0461 sidebar tests (the `nav[aria-label="Charts"]`, `[data-sidebar-entry]`, and controls-inside-`[data-sidebar-content]` assertions, and the group assertions in the planning, strategic, and cost of delay blocks) and add one test that the page has no `nav` of its own and draws its controls and chart directly. In `flaiover/e2e/page-sidebar.e2e.ts`, remove the charts helper and its two cases; the rail's chart list is covered by `layout.e2e.ts` since T-2024.

Fourth layer: waits for T-2024, so the rail lists the charts before the page stops listing them, and every commit leaves a way to pick a chart.

## Done when

`charts.svelte.test.ts` passes with the sidebar tests gone and the new one in, `page-sidebar.e2e.ts` compiles without the charts cases, and `flai test` on the changed paths passes.

## Notes

Drafted by the planner for E-0028. `PageSidebar.svelte` itself and the Docs page's use of it are not touched.
