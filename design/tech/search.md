---
title: Search
updated: 2026-09-15
status: active
---

# Search in flaiover

| | |
|-|-|
| Library | MiniSearch 7 |
| Where | Server side, in the SvelteKit process, index rebuilt on file change |
| Scope | `design/` and `wip/` including `archive/`, `docs/` opt-in via a toggle |

Fields: `id`, `title`, `tags`, `headings`, `body`, `type`, `nature`, `status`, `path`. Prefix and fuzzy matching on, field boosts on `id` and `title`. Results carry a snippet and the route to open.

Chosen over FlexSearch for a simpler API and over a full-text engine because a repo of a few thousand markdown files indexes in well under a second.
