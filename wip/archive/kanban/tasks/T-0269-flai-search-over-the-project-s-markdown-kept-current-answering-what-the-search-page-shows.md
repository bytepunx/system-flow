---
id: T-0269
type: task
nature: feature
title: "flai: search over the project's Markdown, kept current, answering what the search page shows"
status: done
parent: S-0074
owner: alex
created: 2026-09-20T08:29:59Z
updated: 2026-09-20T08:37:39Z
transitions:
  - to: ready
    at: 2026-09-20T08:35:09Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T08:35:09Z
    by: system-flow
  - to: done
    at: 2026-09-20T08:37:39Z
    by: system-flow
stream: S-0074
tags: []
---
# T-0269 flai: search over the project's Markdown, kept current, answering what the search page shows

## Work
An index over the same files and fields as today (item ID, title, tags, headings, body; boosts 4, 3, 2), prefix and fuzzy matching, BM25 ranking, rebuilt when the files' signature changes, answering search.query with the hit fields and snippet the page shows. The approach and why no library, in design/tech/search.md.

## Done when
- Tests for exact, prefix, and fuzzy matches, field boosts, the docs scope, the snippet, and a rebuild after a change

## Notes
