---
id: T-0608
type: task
nature: feature
title: The header's navigation is a two-tier site menu
status: done
parent: S-0172
owner: alex
created: 2026-09-30T01:18:30Z
updated: 2026-09-30T01:24:07Z
transitions:
  - to: ready
    at: 2026-09-30T01:18:44Z
    by: agent-S-0172
  - to: in-progress
    at: 2026-09-30T01:18:44Z
    by: agent-S-0172
  - to: done
    at: 2026-09-30T01:24:07Z
    by: agent-S-0172
stream: S-0172
tags: []
usage:
  source: log
  seconds: 323
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 62
      output: 18298
      cache_read: 3152280
      cache_write: 69472
      cost: 1.5524
---

# T-0608 The header's navigation is a two-tier site menu

## Work

- `$lib/sitemenu.ts`: the three groups (Workflow, Status, Host) and their pages, and which group and page a path belongs to.
- `SiteMenu.svelte`: the parents in the header's top row; the shown group's pages in a row below. Hover, focus, click, or touch on a parent shows its pages; leaving the menu shows the current page's group again. The current page and its parent are bold in the accent colour.
- A page's indicator (the Inbox count) shows on its parent too.
- `+layout.svelte` uses `SiteMenu` in place of the flat list.

## Done when

- Behaviour tests for `sitemenu.ts` and `SiteMenu.svelte` cover showing a group, the active parent and child, and the parent's indicator.
- `npm run check`, `npm run lint`, and `npm run test:unit` pass in `flaiover/`.

## Notes
