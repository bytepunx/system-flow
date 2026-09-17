---
title: flaiover dashboard
updated: 2026-09-15
status: draft
---

# flaiover

Run it with `flai dashboard` and open <http://localhost:4242>. The design is in [design/system/flaiover-dashboard.md](../../design/system/flaiover-dashboard.md); the full guide arrives with story S-019.

## Documentation explorer

Docs shows every markdown file under `design/`, `docs/`, and `wip/` in a collapsible tree. A document renders with its Mermaid diagrams, highlighted code, task-list checkboxes, and heading anchors; links between documents open in the explorer. The front matter is shown in a panel above the text. ADRs lists the architecture decisions with status, date, and which decisions supersede which.

## Board

Board shows one column per state with the WIP count against the limit from `wip/kanban/board.md`. Cards are stories by default; a toggle adds epics and tasks. Each card shows the ID, title, nature, how long it has sat in its column, and a flag when it is blocked. Drag a card to another column to move it: the same rules `flai move` enforces apply, and a refused move shows the rule. Every write is made by the bundled `flai`, so the files change exactly as they would from the terminal and appear in `git status` for you to commit.

Click a card for the item page: front matter, the rendered body, children, the transition history, blocked intervals, a link to the narrative, and buttons for the allowed moves, block, unblock, and adding a narrative log entry.

## Search

Search covers `design/` and `wip/` by default and `docs/` when you tick the box. Type an item ID, a title, or words from the body; each result shows where it is, its status, and a snippet around the match. Results open in the explorer.
