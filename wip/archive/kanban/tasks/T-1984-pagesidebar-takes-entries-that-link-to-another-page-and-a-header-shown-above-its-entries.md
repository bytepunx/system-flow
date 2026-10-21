---
id: T-1984
type: task
nature: feature
title: PageSidebar takes entries that link to another page and a header shown above its entries
status: done
parent: S-0461
owner: alex
created: 2026-10-21T11:58:27Z
updated: 2026-10-21T12:11:32Z
transitions:
  - to: ready
    at: 2026-10-21T11:58:57Z
    by: agent-S-0461
  - to: in-progress
    at: 2026-10-21T11:58:58Z
    by: agent-S-0461
  - to: done
    at: 2026-10-21T12:11:32Z
    by: agent-S-0461
stream: S-0461
tags: []
touches: [flaiover/src/lib/components/PageSidebar.svelte, flaiover/src/lib/components/PageSidebar.svelte.test.ts, flaiover/src/lib/pagesidebar.svelte.ts, flaiover/src/lib/pagesidebar.svelte.test.ts, design/components/flaiover/src/lib/components/PageSidebar.examples.json]
usage:
  source: log
  seconds: 754
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 13774
      cache_read: 1759541
      cache_write: 61022
      cost: 0.9889
---

# T-1984 PageSidebar takes entries that link to another page and a header shown above its entries

## Work

- `SidebarEntry` gains an optional `href`: an entry with one is a link to another page, followed by the router as any link is, with no hash written and no section scrolled to. It is current while the bound `current` names it, as the charts page sets it from its route.
- `PageSidebar` gains an optional `header` snippet, drawn in the sidebar's column above its entries, from `md` sticky with them and below `md` above the row of entries, so that the documents page can put `DocTree` there.
- Tests in `PageSidebar.svelte.test.ts` and `pagesidebar.svelte.test.ts` cover a linked entry and the header. `PageSidebar.examples.json` gains a set with linked entries.
- Waits for nothing. T-1966 and T-1967 wait for it, as both use what it adds.

## Done when

`flai test` passes on the paths it changed, and a page can list links to other pages and a header of its own in the sidebar.

## Notes
