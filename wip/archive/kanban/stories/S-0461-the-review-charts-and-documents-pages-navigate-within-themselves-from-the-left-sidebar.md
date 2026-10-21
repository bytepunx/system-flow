---
id: S-0461
type: story
nature: improvement
title: The review, charts, and documents pages navigate within themselves from the left sidebar
status: done
parent: E-0027
owner: alex
created: 2026-10-21T07:00:44Z
updated: 2026-10-21T13:07:23Z
transitions:
  - to: ready
    at: 2026-10-21T07:54:29Z
    by: alex
  - to: in-progress
    at: 2026-10-21T11:57:41Z
    by: agent-S-0461
  - to: review
    at: 2026-10-21T13:02:10Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T13:07:23Z
    by: orchestrator
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, "flaiover/src/routes/charts/[kind]/+page.svelte", "flaiover/src/routes/charts/[kind]/charts.svelte.test.ts", "flaiover/src/routes/docs/[...path]/+page.svelte", flaiover/src/routes/docs/docs.svelte.test.ts, flaiover/src/lib/components/DocTree.svelte, design/components/flaiover/src/lib/components/DocTree.examples.json, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/PageSidebar.svelte, flaiover/src/lib/components/PageSidebar.svelte.test.ts, flaiover/src/lib/pagesidebar.svelte.ts, flaiover/src/lib/pagesidebar.svelte.test.ts, design/components/flaiover/src/lib/components/PageSidebar.examples.json, flaiover/src/lib/components/DiffView.svelte, flaiover/e2e/page-sidebar.e2e.ts, flaiover/src/lib/components/Threads.svelte, design/issues/I-0185-search-test-ts-finds-a-body-phrase-across-wip-and-design-times-out-at-5-s-on-the-monorepo-while-other-stories-close-outs-run.md, design/issues/summary.md]
after: [S-0460]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3880
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 3
      hand_edits: 7
      work: 56
  models:
    - model: claude-opus-5-5
      input: 424
      output: 196642
      cache_read: 25120399
      cache_write: 871197
      cost: 14.1175
  strategic:
    - kind: orchestrator
      seconds: 619
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 162
          output: 2789
          cache_read: 28248489
          cache_write: 43059
          cost: 6.8489
        - model: claude-sonnet-5-5
          input: 10
          output: 52
          cache_read: 171466
          cache_write: 40828
          cost: 0.1461
forecast:
  duration: 45m
  delivery: 2026-10-21T12:05:00Z
  basis: "Its own forecast of 45m; 1st in the pull order with an in-progress limit of 5, behind S-0460."
  by: flai
  at: 2026-10-21T11:12:19Z
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

- [x] `/review/:id` has a `PageSidebar` whose scroll entries are Criteria, State, Threads, Changes, and Release plan, in the page's order. Under Changes there is one entry per changed file, and choosing one opens that file's panel and scrolls to it. Blockers and the undone-main cards stay above the sections, where Accept is decided.
- [x] `/charts/:kind` lists its groups and kinds in a `PageSidebar` in place of the row above the chart. Each kind is a pane entry that routes to `/charts/<kind>` as the links do now, and the chart window control stays with the chart.
- [x] `/docs/<path>` puts `DocTree` in the sidebar's place. Below it, the open document's headings, from `headingsOf`, are scroll entries, with the heading in view current. Editing a document and its threads work as before.
- [x] Each page's existing tests pass, and new tests cover its sidebar's entries and what choosing one does. Below Tailwind's `md` width no page scrolls sideways.
- [x] `design/system/flaiover-dashboard.md` § Views and `docs/users/flaiover.md` (§ Reviewing a story, § Charts, § Docs) describe each sidebar. `DocTree`'s examples follow any change to its props.

## Tasks
- T-1965 The review page's sidebar scrolls to its criteria, state, threads, each changed file, and release plan
- T-1966 The charts page lists its groups and kinds in the sidebar in place of the row above the chart
- T-1967 The documents page holds the document tree and the open document's headings in the sidebar
- T-1968 The dashboard's design and user docs describe the review, charts, and documents sidebars
- T-1984 PageSidebar takes entries that link to another page and a header shown above its entries
- T-2010 An end-to-end test finds that the review, charts, and documents pages never scroll sideways below md and keep the sidebar beside the content from md
- T-2011 A thread's message, the review narrative's long paths, and the charts page's epic list wrap or narrow within a phone's width

## Notes

- Drafted by the planner for E-0027. Fourth, after S-0460, whose `PageSidebar.svelte` it uses.
- Other pages could take the sidebar later: Inbox by kind of entry, Harnesses by harness, and the Design page by token group. They are left out to keep the story to one or two sessions. They are proposed on the plan thread.

### As built

- **Criterion 2.** A chart kind is an entry that links to `/charts/<kind>`, which `PageSidebar` gained in T-1984, not a hash pane: a hash pane cannot route.
- **Criterion 3.** The heading entries are read from the rendered document, not from `headingsOf`. The renderer slugs a heading's plain text and numbers repeats, so `headingsOf`'s raw markdown would name no element for a heading with inline code, a link, or a repeated title. `Threads` still takes `headingsOf`.
- **Criterion 4.** `flaiover/e2e/page-sidebar.e2e.ts` opens the three pages in Chromium at 375px and 1280px. It found three overflows at 375px, which T-2011 fixed: a thread message's long path, an inline path in the review narrative, and the charts page's epic list. Run it with `scripts/flaiover-e2e.sh e2e/page-sidebar.e2e.ts`; no tier of the manifest runs the e2e tests.

### Planning

- **Touches.** None were declared, so every touch is the planner's. All are files, and no folder touch was kept.
  - **Layout:** `flaiover/src/lib/components/Review.svelte` and its test. The route `review/[id]/+page.svelte` is 13 lines that only render it, so it is not touched. `flaiover/src/routes/charts/[kind]/+page.svelte` (807 lines, its group-and-kind `<nav>` at line 267) and `charts.svelte.test.ts`. `flaiover/src/routes/docs/[...path]/+page.svelte`, which imports `DocTree` and `headingsOf`, with `docs.svelte.test.ts`. `DocTree.svelte`, in case fitting the column changes its props.
  - **Design:** `design/system/flaiover-dashboard.md` § Views, and `docs/users/flaiover.md` § Reviewing a story, § Charts, and § Docs. `DocTree.examples.json` follows `DocTree`'s props, as ADR-0164 requires.
  - **Co-change:** `flai touches suggest` lists `flaiover/src/lib/viz/charts.ts` and its test at 7%. They build the charts' data, not the navigation, so none was added.
- **Forecast.** 45m against flai's 23m. The story changes three pages, each with its own tests, and the charts page alone is 807 lines. The delivery is flai's play-out, moved by its own 22m raise and the raises of S-0460 and S-0458.
- **Cost of delay.** Not yet set. E-0027 has no inputs, so `flai cod` refuses. The input is asked on TH-0553 question 1.

### Accepted by the orchestrator

- Verified: 71dd293e8e2808b60aff65de3c647b9d109c3784
- Batch: 20261021-130439
- Enqueued verified: c4d04d3bdee1ca4662815db8c6721ceb56f0d596
- At: 2026-10-21T13:07:23Z

Verdict: PASS. All five criteria are met. Criterion 3 reads headings from the rendered HTML in place of headingsOf, for the reason recorded in the story's notes.
- 1: flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/DiffView.svelte, flaiover/src/lib/components/Review.svelte.test.ts
- 2: flaiover/src/routes/charts/[kind]/+page.svelte, flaiover/src/routes/charts/[kind]/charts.svelte.test.ts, flaiover/src/lib/components/PageSidebar.svelte, flaiover/src/lib/pagesidebar.svelte.ts, design/components/flaiover/src/lib/components/PageSidebar.examples.json
- 3: flaiover/src/routes/docs/[...path]/+page.svelte, flaiover/src/routes/docs/docs.svelte.test.ts
- 4: flaiover/src/lib/components/PageSidebar.svelte.test.ts, flaiover/src/lib/pagesidebar.svelte.test.ts, flaiover/e2e/page-sidebar.e2e.ts, flaiover/src/lib/components/Threads.svelte
- 5: design/system/flaiover-dashboard.md, docs/users/flaiover.md
