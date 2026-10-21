---
id: T-1965
type: task
nature: improvement
title: The review page's sidebar scrolls to its criteria, state, threads, each changed file, and release plan
status: done
parent: S-0461
owner: alex
created: 2026-10-21T07:04:13Z
updated: 2026-10-21T12:11:42Z
transitions:
  - to: ready
    at: 2026-10-21T11:58:56Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T11:58:56Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T12:11:42Z
    by: agent-S-0461
stream: S-0461
tags: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/DiffView.svelte]
usage:
  source: log
  seconds: 766
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 36627
      cache_read: 4679022
      cache_write: 162273
      cost: 2.6296
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
