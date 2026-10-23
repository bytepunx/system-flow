---
id: T-2082
type: task
nature: feature
title: The nav rail's Configuration group links the Host page
status: done
parent: S-0489
owner: alex
created: 2026-10-23T18:06:57Z
updated: 2026-10-23T18:42:50Z
transitions:
  - to: ready
    at: 2026-10-23T18:21:59Z
    by: agent-S-0489
  - to: in-progress
    at: 2026-10-23T18:22:00Z
    by: agent-S-0489
  - to: done
    at: 2026-10-23T18:42:50Z
    by: agent-S-0489
stream: S-0489
tags: [flaiover, dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/NavRail.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 1250
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 4594
      cache_read: 957135
      cache_write: 26731
      cost: 0.4429
---
# T-2082 The nav rail's Configuration group links the Host page

## Work

The Host page exists at `/host` (`flaiover/src/routes/host/+page.svelte`) but no group of the rail lists it since S-0458: only the host-flai pill in the top bar links it. Put it in the Configuration group of `SITE_MENU` in `flaiover/src/lib/sitemenu.ts`: add `host` to `PageKey` and a sixth page after Harnesses, label `Host`, a lucide icon imported as the others are (such as `@lucide/svelte/icons/server`), `path: '/host'`, `href: resolve('/host')`. The route stays `/host`, so the pill's link, `layout.svelte.test.ts`, and bookmarks keep working; the group's comment says the Host page is the one page of the group that is not a settings page.

`locate('/host')` then answers `{ group: 'configuration', page: 'host', child: null }`: change the test in `sitemenu.test.ts` that expects null for it, add `/host` to the Configuration group's hrefs in the groups-and-pages test, and cover the rendered link in `NavRail.svelte.test.ts`.

Update the rail table in `design/system/flaiover-dashboard.md` § Views (the Configuration row, and the line that says Updates is in no group) and `docs/users/flaiover.md` (the header table's Configuration row, and the first line of § Host, which says no group of the rail has it).

This task waits for no other.

## Done when

- The rail shows Host under Configuration, linking to `/host`, and the page is marked current there when it is open.
- `locate('/host')` places it in the Configuration group, and the sitemenu and NavRail tests say so and pass.
- The rail tables in `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` list Host under Configuration, and `flai test` on the changed paths passes.

## Notes

Planned by the planner on 2026-10-23: the route stays `/host` so nothing that links it changes; only the rail gains the entry.
