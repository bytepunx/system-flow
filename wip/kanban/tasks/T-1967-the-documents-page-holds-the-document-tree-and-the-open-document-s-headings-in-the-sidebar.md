---
id: T-1967
type: task
nature: improvement
title: The documents page holds the document tree and the open document's headings in the sidebar
status: done
parent: S-0461
owner: alex
created: 2026-10-21T07:04:24Z
updated: 2026-10-21T12:17:35Z
transitions:
  - to: ready
    at: 2026-10-21T12:12:07Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T12:12:08Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T12:17:35Z
    by: agent-S-0461
stream: S-0461
tags: [dashboard]
touches: ["flaiover/src/routes/docs/[...path]/+page.svelte", flaiover/src/routes/docs/docs.svelte.test.ts, flaiover/src/lib/components/DocTree.svelte, design/components/flaiover/src/lib/components/DocTree.examples.json]
after: [T-1984]
usage:
  source: log
  seconds: 327
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 30
      output: 167
      cache_read: 1033899
      cache_write: 72265
      cost: 0.3747
---
# T-1967 The documents page holds the document tree and the open document's headings in the sidebar

## Work

- In `flaiover/src/routes/docs/[...path]/+page.svelte`, put `DocTree` in the sidebar column. Below it, add a `PageSidebar` of scroll entries from `headingsOf` on the open document, with the heading in view current. Give each rendered heading the id its entry scrolls to, unless the markdown renderer already gives one.
- Change `DocTree.svelte` only as far as fitting the column needs. Keep `DocTree.examples.json` in step with any prop it gains.
- Extend `docs.svelte.test.ts`: the tree is in the sidebar, the headings list the document's, and editing and threads still work.
- Waits for T-1984, which gives the sidebar a header for `DocTree`. It runs beside the review and charts tasks.

## Done when

`docs.svelte.test.ts` passes, and a long document's headings are one choice away in the sidebar.

## Notes
