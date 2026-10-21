---
id: S-0461
type: story
nature: improvement
title: The review, charts, and documents pages navigate within themselves from the left sidebar
status: ready
parent: E-0027
owner: alex
created: 2026-10-21T07:00:44Z
updated: 2026-10-21T08:22:18Z
transitions:
  - to: ready
    at: 2026-10-21T07:54:29Z
    by: alex
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", "flaiover/src/routes/docs/[...path]/+page.svelte", flaiover/src/routes/docs/docs.svelte.test.ts, flaiover/src/lib/components/DocTree.svelte, design/components/flaiover/src/lib/components/DocTree.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0460]
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
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 59
          cache_read: 271604
          cache_write: 1064
          cost: 0.0684
forecast:
  duration: 45m
  delivery: 2026-10-21T09:14:00Z
  basis: "Its own forecast of 45m; 2nd in the pull order with an in-progress limit of 5, behind S-0456, S-0458, S-0463, S-0465 and S-0459."
  by: flai
  at: 2026-10-21T08:22:18Z
finalized:
  by: alex
  at: 2026-10-21T07:53:56Z
---
# S-0461 The review, charts, and documents pages navigate within themselves from the left sidebar

## Goal

Three more of the pages the operator scrolls through most carry the sidebar S-0460 builds, so that each is navigated from its left edge.

- **Review**: the review of a story in `Review.svelte` scrolls through its criteria, the narrative's state, its threads, a diff of every file, and the release plan.
- **Charts**: its group and kind links sit as a row above the chart.
- **Documents**: the document tree and the document's own headings have no shared place.

## Acceptance criteria

- [ ] `/review/:id` has a `PageSidebar` whose scroll entries are Criteria, State, Threads, Changes, and Release plan, in the page's order. Under Changes there is one entry per changed file, and choosing one opens that file's panel and scrolls to it. Blockers and the undone-main cards stay above the sections, where Accept is decided.
- [ ] `/charts/:kind` lists its groups and kinds in a `PageSidebar` in place of the row above the chart. Each kind is a pane entry that routes to `/charts/<kind>` as the links do now, and the chart window control stays with the chart.
- [ ] `/docs/<path>` puts `DocTree` in the sidebar's place. Below it, the open document's headings, from `headingsOf`, are scroll entries, with the heading in view current. Editing a document and its threads work as before.
- [ ] Each page's existing tests pass, and new tests cover its sidebar's entries and what choosing one does. Below Tailwind's `md` width no page scrolls sideways.
- [ ] `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md` (§ Reviewing a story, § Charts, § Docs) describe each sidebar. `DocTree`'s examples follow any change to its props.

## Tasks
- T-1965 The review page's sidebar scrolls to its criteria, state, threads, each changed file, and release plan
- T-1966 The charts page lists its groups and kinds in the sidebar in place of the row above the chart
- T-1967 The documents page holds the document tree and the open document's headings in the sidebar
- T-1968 The dashboard's design and user docs describe the review, charts, and documents sidebars

## Notes

- Drafted by the planner for E-0027. Fourth, after S-0460, whose `PageSidebar.svelte` it uses.
- Other pages could take the sidebar later: Inbox by kind of entry, Harnesses by harness, and the Design page by token group. They are left out to keep the story to one or two sessions. They are proposed on the plan thread.

### Planning

- **Touches.** None were declared, so every touch is the planner's. All are files, and no folder touch was kept.
  - **Layout:** `flaiover/src/lib/components/Review.svelte` and its test. The route `review/[id]/+page.svelte` is 13 lines that only render it, so it is not touched. `flaiover/src/routes/charts/[kind]/+page.svelte` (807 lines, its group-and-kind `<nav>` at line 267) and `charts.svelte.test.ts`. `flaiover/src/routes/docs/[...path]/+page.svelte`, which imports `DocTree` and `headingsOf`, with `docs.svelte.test.ts`. `DocTree.svelte`, in case fitting the column changes its props.
  - **Design:** `design/system/flaiover-dashboard.md` § Views, and `docs/users/flaiover.md` § Reviewing a story, § Charts, and § Docs. `DocTree.examples.json` follows `DocTree`'s props, as ADR-0164 requires.
  - **Co-change:** `flai touches suggest` lists `flaiover/src/lib/viz/charts.ts` and its test at 7%. They build the charts' data, not the navigation, so none was added.
- **Forecast.** 45m against flai's 23m. The story changes three pages, each with its own tests, and the charts page alone is 807 lines. The delivery is flai's play-out, moved by its own 22m raise and the raises of S-0460 and S-0458.
- **Cost of delay.** Not yet set. E-0027 has no inputs, so `flai cod` refuses. The input is asked on TH-0553 question 1.
