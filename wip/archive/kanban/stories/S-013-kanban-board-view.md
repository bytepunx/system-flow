---
id: S-013
type: story
nature: feature
title: Kanban board view with transitions
status: done
parent: E-003
owner: agent
created: 2026-09-15T16:09:00Z
updated: 2026-09-17T04:50:10Z
transitions:
  - to: ready
    at: 2026-09-17T04:32:27Z
    by: agent
  - to: in-progress
    at: 2026-09-17T04:32:27Z
    by: agent
  - to: review
    at: 2026-09-17T04:38:33Z
    by: agent
  - to: done
    at: 2026-09-17T04:50:10Z
    by: alex
tags: []
---

# S-013 Kanban board view with transitions

## Goal
Board with columns from board.md, cards with age and blocked flag, item detail page, and drag-to-transition that writes validated front matter.

## Acceptance criteria
- [x] Board matches flai board output
- [x] Transition via drag appends to transitions and updates status
- [x] Invalid transitions are refused with the rule named

## Tasks
- T-090 ADR-0016: dashboard writes and metrics delegate to the flai binary; server wrapper with JSON and error mapping
- T-091 API: GET /api/board; POST move, block, unblock on items; POST stream log
- T-092 Board page with columns, cards, WIP counts, drag to transition; item page with history, children, actions, narrative link
- T-093 Tests against a temp project with the built flai; docs and design

## Notes
- Writes delegate to the flai binary (ADR-0016) instead of a TypeScript port of the rules; the image must bundle flai (S-015).
