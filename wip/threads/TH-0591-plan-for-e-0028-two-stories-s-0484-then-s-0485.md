---
id: TH-0591
title: "Plan for E-0028: two stories, S-0484 then S-0485"
anchor:
  path: wip/kanban/epics/E-0028-charts-improvements.md
  item: E-0028
status: open
participants: [planner-E-0028, orchestrator]
created: 2026-10-21T21:02:07Z
updated: 2026-10-21T21:02:16Z
---

# TH-0591 Plan for E-0028: two stories, S-0484 then S-0485

On wip/kanban/epics/E-0028-charts-improvements.md.

## Entries

### 2026-10-21T21:02:07Z planner-E-0028
Two draft stories deliver the outcome, one per bullet of it, and both are enriched with touches, a forecast, and tasks in layers. Their cost of delay waits on TH-0590.

**S-0484** The rail's Charts item expands to list the charts by group, each a link that opens it, and the charts page keeps no chart list of its own. Forecast 45m.

- Layer 1: T-2018 the chart catalogue (kinds, groups, titles) moves to `viz/catalog.ts`, re-exported by `viz/charts.ts`, so the rail imports it without the builders; T-2019 `navrail.svelte.ts` keeps whether Charts is expanded under `flaiover-nav-rail`.
- Layer 2: T-2022 `sitemenu.ts` gives Charts `children` from the catalogue and `locate` names the chart shown.
- Layer 3: T-2024 `NavRail.svelte` draws Charts as an expandable item with the four labels and the chart links, the chart shown current, and `layout.e2e.ts` opens a chart from the rail.
- Layer 4: T-2026 the page drops `PageSidebar` and its S-0461 tests, `page-sidebar.e2e.ts` loses its charts cases.
- Layer 5: the docs task, `docs/users/flaiover.md` and `flaiover-dashboard.md`.

**S-0485** The chart page shows the title and the chart first, then a collapsible explainer of what the chart shows and how it is measured, then the controls, then the table view. `after: S-0484`, since both change the page, its test, and the same two documents. Forecast 1h.

- Layer 1: T-2020 `viz/explain.ts`, an explainer per kind (`what`, `how` it is measured from `metrics.md`'s transitions and fields, `read`, and the `metrics.md` anchor), with a test that every kind has one and each anchor is a heading; T-2021 `chartexplainer.svelte.ts` remembers open or closed per browser.
- Layer 2: T-2023 `ChartExplainer.svelte`, a `<details>` "About this chart" linking the definition on the Docs page, with its examples file.
- Layer 3: T-2025 the page reordered: title, notices, chart, the figures it states, explainer, controls with the summary strip, table view; the static notes under charts fold into the explainers.
- Layer 4: T-2027 the docs.

**Assumptions, each yours to overturn when finalizing:**

1. The page's own chart list goes. S-0461 (accepted today) put the list in a page sidebar beside the chart; the outcome's first sentence describes the row before it. I read "instead" as one picker, in the rail. Say so if you want both.
2. Charts in the rail links nowhere itself: clicking it expands the list, and `/charts` typed still redirects to cycle time.
3. The rail labels spend charts by their bucket-free titles (Tokens / Day, $ / Day), since it does not know the bucket; the page's title still follows the bucket.
4. The explainer is open on the first visit and stays as you last left it, one choice for every chart, kept in the browser like the window.
5. The figures a chart states (the saving under CoD by Order, the ratio under Strategic Cost, the stories without a value) stay directly under the chart; only the explanatory notes move into the explainer. The summary strip (completed, p50, p85, spend) stays with the controls.
6. The explainers link `metrics.md` on the Docs page rather than copying it at length; `metrics.md` is not changed.

**Proposed, not drafted:** nothing to split, merge, or drop. If you would rather keep the page sidebar as well, S-0484 shrinks to the rail alone (T-2026 and the sidebar part of its docs go) and S-0485 reorders inside the sidebar's content.

### 2026-10-21T21:02:16Z orchestrator (recommendation)
Recommended for the operator: approve the plan as drafted, with its six assumptions, and resolve this thread.

- **The two stories map one-to-one onto the outcome's two bullets.** S-0485 waits on S-0484 by `after`, since both change the chart page, its test, and the same two documents. Forecasts are 45m and 1h. Each story's tasks form layers with no shared path within a layer.
- **Assumption 1, the only one that changes scope:** take it. Your outcome says clicking a chart in the expanded Charts item "is what displays it on the chart page". That makes the rail the one picker, so S-0461's page sidebar goes. Keeping both would give two lists of the same charts. If you want both, the planner's fallback shrinks S-0484 to the rail alone.
- **Assumptions 2 to 6** follow the outcome's wording or settle details. Each is yours to overturn when you finalize.

The stories get their values once you set E-0028's input on TH-0590. I then finalize them under finalize_drafts.

Source: design/conventions/strategic-agents.md § As the planner
