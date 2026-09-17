---
id: S-0012
type: story
nature: feature
title: Documentation explorer and search
status: done
parent: E-0003
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T04:32:23Z
transitions:
  - to: ready
    at: 2026-09-17T04:24:14Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:24:14Z
    by: agent
  - to: review
    at: 2026-09-17T04:29:52Z
    by: agent
  - to: done
    at: 2026-09-17T04:32:23Z
    by: alex
tags: []
---

# S-0012 Documentation explorer and search

## Goal
Tree navigation, rendered markdown with Mermaid and code highlighting, front matter panel, ADR list, and MiniSearch across design and wip.

## Acceptance criteria
- [x] Every markdown file in the repo renders with working relative links
- [x] Mermaid blocks render
- [x] Search returns items and documents with snippets

## Tasks
- T-0086 Search index (MiniSearch) over design and wip with docs opt-in, rebuilt on change; /api/search with snippets; /api/docs/adrs
- T-0087 Markdown rendering: markdown-it with anchors and task lists, relative links rewritten to explorer routes, mermaid and shiki lazy-loaded
- T-0088 Explorer, ADR list, and search pages with the front matter panel; nav links
- T-0089 Tests for the index and the renderer, docs, tech pins

## Notes
- Verified against this repository: every markdown file is in the tree, Mermaid and code fences render client side, search finds items and phrases with snippets.
