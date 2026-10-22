---
id: S-0484
type: story
nature: improvement
title: The rail's Charts item expands to list the charts by group, each a link that opens it, and the charts page keeps no chart list of its own
status: backlog
parent: E-0028
owner: alex
created: 2026-10-21T20:55:01Z
updated: 2026-10-22T20:17:07Z
transitions: []
tags: [dashboard]
topics: [navigation]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/navrail.svelte.test.ts, flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/viz/catalog.ts, flaiover/src/lib/viz/catalog.test.ts, flaiover/src/lib/viz/charts.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", flaiover/e2e/page-sidebar.e2e.ts, flaiover/e2e/layout.e2e.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 45m
  delivery: 2026-10-22T15:56:00Z
  basis: "flai forecast's 33m at size 22 (7 criteria, 15 touches), 89 s a unit over 86 done improvement stories on claude-opus-5-5 in the large band, raised to 45m for the two Playwright files, which the size counts as one unit each but whose every change runs the browser suite in the docked rail and the drawer; delivery is flai's play-out, 39th in the pull order with an in-progress limit of 5, 2026-10-22T15:44Z, plus the 12m raised."
  by: planner-E-0028
  at: 2026-10-21T21:02:56Z
finalized:
  by: alex
  at: 2026-10-22T20:17:07Z
---
# S-0484 The rail's Charts item expands to list the charts by group, each a link that opens it, and the charts page keeps no chart list of its own

## Goal

The dashboard's charts are chosen from the nav rail. Charts, in the rail's Flow group, becomes an item that expands when clicked, like a group, and lists every chart in a vertical list under four muted labels, flow, usage, planning, and strategic, in the order the charts page lists them today. Clicking a chart opens `/charts/<kind>`; the chart shown is marked current, and Charts stays current while any chart is shown. The charts page then draws no list of charts of its own: the sidebar S-0461 gave it goes, and its content takes the width beside the rail. The first half of E-0028's outcome.

## Acceptance criteria

- [ ] In `NavRail.svelte`, Charts is a button with `aria-expanded` and `aria-controls`, its icon in the accent colour while expanded and muted while collapsed, as a group's is. Expanded, it lists under it the four chart groups, `flow`, `usage`, `planning`, and `strategic`, each a muted label that is not a link, and under each label its charts in the order `KINDS` gives them, each an `<a>` to `/charts/<kind>` through `resolve` and `withProject`, labelled by `TITLES`. Collapsed, nothing below Charts is in the DOM, and Charts still links nowhere but through its charts.
- [ ] The chart shown is bold in the accent colour with `aria-current="page"`, and Charts is marked current while any `/charts/*` path is shown; `locate` in `sitemenu.ts` gives the chart's kind beside the group and the page, and an unknown kind under `/charts` marks Charts current with no chart.
- [ ] Arriving on `/charts/<kind>` expands Flow and Charts for the page's life without storing it. Clicking Charts expands or collapses it and keeps that choice per browser under `flaiover-nav-rail` beside the groups' choices; `parseChosen` keeps a boolean for Charts and drops anything else unknown, so a stored value from an older dashboard still parses. Choosing a chart in the drawer below `md` closes the drawer, as choosing a page does.
- [ ] The chart catalogue the rail and the page share, the kinds, their groups, and their titles, lives in `flaiover/src/lib/viz/catalog.ts`, imported by `sitemenu.ts` without pulling `viz/charts.ts` and its builders into the layout's bundle; `viz/charts.ts` re-exports it so its callers and tests are unchanged.
- [ ] `/charts/[kind]/+page.svelte` no longer uses `PageSidebar`: the page has no `nav[aria-label="Charts"]`, and its controls, chart, notes, and table view are drawn as they are today, full width beside the rail. `flaiover/src/routes/charts/+page.ts` still redirects `/charts` to cycle time.
- [ ] `sitemenu.test.ts`, `navrail.svelte.test.ts`, `NavRail.svelte.test.ts`, and `charts.svelte.test.ts` cover the criteria above, the S-0461 sidebar tests in `charts.svelte.test.ts` and the charts cases of `e2e/page-sidebar.e2e.ts` replaced by tests of the rail's chart list, `e2e/layout.e2e.ts` opening a chart from the rail; `flai test` on the changed paths passes, `design/components/flaiover/src/lib/components/NavRail.examples.json` stays in step with `NavRail`'s props, and `flai check --strict` is clean.
- [ ] `docs/users/flaiover.md` § The header and § Charts, and `design/system/flaiover-dashboard.md` § Views, the rail's groups table, the rail behaviour table, the page sidebar paragraph, and the `/charts/*` rows, say that the charts are chosen from the rail and that the page has no list of its own.

## Tasks
- T-2018 The chart catalogue, the kinds, their four groups, and their titles, moves to viz/catalog.ts, which viz/charts.ts re-exports
- T-2019 The nav rail's state keeps whether Charts is expanded beside the groups, under the same localStorage key, and expands it on arrival at a chart
- T-2022 sitemenu.ts gives the Charts page its chart list by group from the catalogue, and locate names the chart a /charts path shows
- T-2024 NavRail.svelte draws Charts as an expandable item listing the chart groups and their charts, the chart shown current, and the e2e layout test opens a chart from the rail
- T-2026 The charts page drops its PageSidebar and chart list, and the sidebar tests of charts move to tests of the rail's list
- T-2028 The user docs and the dashboard design say the charts are chosen from the rail's Charts item and the page has no list of its own

## Notes

- Drafted by the planner for E-0028. The second story of the epic, the chart page's order and explainer, waits for this one, since both change `/charts/[kind]/+page.svelte`, its test, and the same two documents.
- This reverses the chart list S-0461 put in a page sidebar beside the chart, accepted on 2026-10-21: the epic's outcome asks for the list in the rail instead, and one picker is enough. The operator confirms this on the plan thread when finalizing.
- The rail labels the spend charts by their bucket-free titles, Tokens / Day and $ / Day, since the rail does not know the bucket the page has chosen; the page's own title still follows the bucket (`titleOf`).

### Planning

- **Plan thread.** TH-0591 names the two stories, their order, these tasks and their layers, and the assumptions, the reversal of S-0461's sidebar first among them.
- **Layers.** T-2018 (the catalogue) and T-2019 (the rail's state) first, sharing no path; T-2022 (`sitemenu.ts`) after T-2018, since it imports the catalogue; T-2024 (`NavRail.svelte`, `layout.e2e.ts`) after T-2019 and T-2022; T-2026 (the page drops its sidebar, `page-sidebar.e2e.ts`) after T-2024, so the rail lists the charts before the page stops listing them and every commit leaves a way to pick a chart; T-2028 (the docs) last, after T-2026.
- **Touches.** 15, all files, every one on a task.
  - **From the code layout** (the explorer's map of the charts page, the rail, and their tests): `sitemenu.ts` and its test, where the rail's pages and `locate` are; `navrail.svelte.ts` and its test, the expansion state; `NavRail.svelte` and its test, the rendering; `viz/charts.ts`, which holds the kind lists and titles today, and the new `viz/catalog.ts` with its test, so the rail imports the catalogue without the builders and ECharts; `charts/[kind]/+page.svelte` and `charts.svelte.test.ts`, where the sidebar and its S-0461 tests are; `e2e/page-sidebar.e2e.ts`, whose charts cases go, and `e2e/layout.e2e.ts`, which gains the rail's chart case.
  - **From design:** `docs/users/flaiover.md` § The header and § Charts, and `design/system/flaiover-dashboard.md` § Views, which describe the rail, the page sidebar, and the `/charts/*` routes.
  - **Co-change, not added:** `flai touches suggest S-0484` lists only paths these files co-changed with in unrelated stories, `design/system/flai-cli.md` 38%, `docs/users/flai.md` 38%, `docs/operators/index.md` 27%, and `viz/charts.test.ts` 10%; `charts.ts` re-exports what moves, so its test is unchanged.
  - **Checked, not added:** `NavRail.examples.json`, whose props do not change (`flai task done` widens the touches if they do); `PageSidebar.svelte` and the Docs page, untouched; `routes/charts/+page.ts`, whose redirect stays; `+layout.svelte`, which mounts the rail unchanged.
  - No folder touch.
- **Forecast.** flai's 33m at size 22 (7 criteria, 15 touches), 89 s a unit over 86 done improvement stories on claude-opus-5-5 in the large band, raised to 45m: the two Playwright files count one unit each, but each change to them runs the browser suite, in the docked rail and in the drawer, which the unit does not weigh. Delivery is flai's play-out, 39th in the pull order with an in-progress limit of 5, 2026-10-22T15:44Z, plus the 12m raised: 2026-10-22T15:56Z.
- **Cost of delay.** None yet: E-0028 carries no inputs, and `flai cod S-0484` refuses without them. TH-0590 asks the operator for them, recommending 1h of time lost per cycle, 150 USD a week at `planning.hour_rate` 150; the value is written here once the input is on the epic.
