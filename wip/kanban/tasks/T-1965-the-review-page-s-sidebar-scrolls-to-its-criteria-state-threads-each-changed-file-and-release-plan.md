---
id: T-1965
type: task
nature: improvement
title: The review page's sidebar scrolls to its criteria, state, threads, each changed file, and release plan
status: backlog
parent: S-0461
owner: alex
created: 2026-10-21T07:04:13Z
updated: 2026-10-21T07:04:13Z
transitions: []
stream: S-0461
tags: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts]
---
# T-1965 The review page's sidebar scrolls to its criteria, state, threads, each changed file, and release plan

## Work

- In `flaiover/src/lib/components/Review.svelte`, keep the blockers, the undone-main cards, and Accept above the sections. Give each section an id, and add a `PageSidebar` of scroll entries: Criteria, State, Threads, Changes, and Release plan.
- Under Changes, add one child entry per changed file from the diff summary. Choosing one opens that file's panel in `DiffView` and scrolls to it.
- Extend `Review.svelte.test.ts`: the entries match the sections shown, and a file entry opens its panel.
- Waits for nothing in this story. It runs beside the charts and documents tasks, which touch other files.

## Done when

`Review.svelte.test.ts` passes, and every section and changed file of a review is one choice away in the sidebar.

## Notes
