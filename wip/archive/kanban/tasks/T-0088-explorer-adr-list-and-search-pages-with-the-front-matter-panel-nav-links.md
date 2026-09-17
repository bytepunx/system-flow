---
id: T-0088
type: task
nature: feature
title: "Explorer, ADR list, and search pages with the front matter panel; nav links"
status: done
parent: S-0012
owner: alex
created: 2026-09-17T04:24:14Z
updated: 2026-09-17T04:29:51Z
transitions:
  - to: ready
    at: 2026-09-17T04:29:51Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:29:51Z
    by: agent
  - to: done
    at: 2026-09-17T04:29:51Z
    by: agent
stream: S-0012
tags: [dashboard]
---

# T-0088 Explorer, ADR list, and search pages with the front matter panel; nav links

## Work
/docs/[...path]: collapsible tree from /api/docs/tree, rendered content, front matter panel; /adrs table with status and supersedes/superseded_by links; /search page with query, scope toggle for docs, results with snippets; nav gains Docs, ADRs, Search via resolve().

## Done when
Every markdown file in this repository opens from the tree and renders; ADR list matches design/adrs/README.md; search finds S-0011 and 'front matter'.

## Notes
