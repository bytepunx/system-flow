---
id: S-0461
type: story
nature: improvement
title: The review, charts, and documents pages navigate within themselves from the left sidebar
status: backlog
parent: E-0027
owner: alex
created: 2026-10-21T07:00:44Z
updated: 2026-10-21T07:06:08Z
transitions: []
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", "flaiover/src/routes/docs/[...path]/+page.svelte", flaiover/src/routes/docs/docs.svelte.test.ts, flaiover/src/lib/components/DocTree.svelte, design/components/flaiover/src/lib/components/DocTree.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0460]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 45m
  delivery: 2026-10-22T02:20:00Z
  basis: "flai's 23m from 89s per unit of size, raised for three pages changed with their tests, the charts page at 807 lines; delivery is flai's play-out moved by its own 22m raise and those of S-0460 and S-0458."
  by: planner-E-0027
  at: 2026-10-21T07:06:08Z
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
