---
id: T-0506
type: task
nature: feature
title: Rank the sections not selected with flai/internal/search, and catalog the documents not loaded
status: done
parent: S-0137
owner: alex
created: 2026-09-29T00:36:07Z
updated: 2026-09-29T00:40:24Z
transitions:
  - to: ready
    at: 2026-09-29T00:36:13Z
    by: agent-S-0137
  - to: in-progress
    at: 2026-09-29T00:39:13Z
    by: agent-S-0137
  - to: done
    at: 2026-09-29T00:40:24Z
    by: agent-S-0137
stream: S-0137
tags: [cli]
touches: [flai/internal/context, flai/internal/search]
---
# T-0506 Rank the sections not selected with flai/internal/search, and catalog the documents not loaded

## Work

- Index the sections left unselected as `search.Doc`s and pick the five best ADRs and five best design sections against the story's title, goal, and criteria.
- Items: the selection grouped into printed items with path, heading path, reason, other reasons, size.
- Catalog: one line per document not loaded, and the heading outline of each loaded in part.

## Done when

- Behavior tests for ranking and the catalog pass; no new dependency.

## Notes
