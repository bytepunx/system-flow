---
title: flaiover dashboard
updated: 2026-09-15
status: draft
---

# flaiover

Run it with `flai dashboard` and open <http://localhost:4242>. The design is in [design/system/flaiover-dashboard.md](../../design/system/flaiover-dashboard.md); the full guide arrives with story S-019.

## Documentation explorer

Docs shows every markdown file under `design/`, `docs/`, and `wip/` in a collapsible tree. A document renders with its Mermaid diagrams, highlighted code, task-list checkboxes, and heading anchors; links between documents open in the explorer. The front matter is shown in a panel above the text. ADRs lists the architecture decisions with status, date, and which decisions supersede which.

## Search

Search covers `design/` and `wip/` by default and `docs/` when you tick the box. Type an item ID, a title, or words from the body; each result shows where it is, its status, and a snippet around the match. Results open in the explorer.
