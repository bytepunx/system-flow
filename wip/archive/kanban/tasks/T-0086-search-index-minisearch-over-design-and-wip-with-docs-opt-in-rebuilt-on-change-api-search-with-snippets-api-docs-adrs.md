---
id: T-0086
type: task
nature: feature
title: "Search index (MiniSearch) over design and wip with docs opt-in, rebuilt on change; /api/search with snippets; /api/docs/adrs"
status: done
parent: S-0012
owner: alex
created: 2026-09-17T04:24:13Z
updated: 2026-09-17T04:29:50Z
transitions:
  - to: ready
    at: 2026-09-17T04:29:50Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:29:50Z
    by: agent
  - to: done
    at: 2026-09-17T04:29:50Z
    by: agent
stream: S-0012
tags: [dashboard, search]
---

# T-0086 Search index (MiniSearch) over design and wip with docs opt-in, rebuilt on change; /api/search with snippets; /api/docs/adrs

## Work
src/lib/server/search.ts builds a MiniSearch index over every markdown file under the design and wip folders (docs included when requested) with fields id, title, tags, headings, body, type, nature, status, path, kind; rebuilds debounced on Repo change events; /api/search?q=&docs=true returns hits with a snippet around the first match and the route to open; /api/docs/adrs lists ADR front matter with the supersession chain.

## Done when
Unit tests on the fixture and this repo: hits for an item ID, a title word, and a body phrase; ADR list shape.

## Notes
