---
title: Search
updated: 2026-09-15
status: active
---

# Search in flaiover

| | |
|-|-|
| Library | MiniSearch 7.2 |
| Where | Server side, in the SvelteKit process, index rebuilt on file change |
| Scope | `design/` and `wip/` including `archive/`, `docs/` opt-in via a toggle |

Fields: `itemId`, `title`, `tags`, `headings`, `body`; stored `path`, `kind`, `scope`, `status`, `type`. Prefix and fuzzy matching on, boosts on `itemId`, `title`, `headings`. Results carry a snippet around the first match and the route to open (`/items/<id>` for work items, `/docs/<path>` for documents). Rebuilt, debounced, on watcher events. `src/lib/server/search.ts` (S-012).

Chosen over FlexSearch for a simpler API and over a full-text engine because a repo of a few thousand markdown files indexes in well under a second.
