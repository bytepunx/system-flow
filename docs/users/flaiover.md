---
title: flaiover dashboard
updated: 2026-09-18
status: draft
---

# flaiover

Run it with `flai dashboard` and open <http://localhost:4242>. The design is in [design/system/flaiover-dashboard.md](../../design/system/flaiover-dashboard.md); the full guide arrives with story S-0019.

## Theme

flaiover uses the brand palette in a light and a dark theme. It follows your system preference by default; the button at the right of the navigation cycles system, light, and dark, and the choice is remembered per browser. Both themes are checked for readable contrast, and the charts use palettes validated for colour-vision deficiency in each theme.

## Documentation explorer

Docs shows every markdown file under `design/`, `docs/`, and `wip/` in a collapsible tree. A document renders with its Mermaid diagrams, highlighted code, task-list checkboxes, and heading anchors; links between documents open in the explorer. The front matter is shown in a panel above the text. ADRs lists the architecture decisions with status, date, and which decisions supersede which.

## Board

Board shows one column per state with the WIP count against the limit from `wip/kanban/board.md`. Cards are stories by default; a toggle adds epics and tasks. Each card shows the ID, title, nature, how long it has sat in its column, and a flag when it is blocked. In its bottom right corner a story card shows the ID of its epic, and a task card the ID of its story; hover over it, or use a screen reader, for the parent's title. The whole card is one link to the item, so the parent ID is not a link of its own. Drag a card to another column to move it: the same rules `flai move` enforces apply, and a refused move shows the rule. Dropping a story on done accepts it, so the dashboard first shows what acceptance will do, the branch to be merged and the release it will cut, and only proceeds when you confirm; cancelling leaves the card in review. When acceptance cannot run from the dashboard, the confirmation lists why and what to do, and its accept button stays disabled: the usual case is a dashboard that sees the repository at a different path than your machine does, so git cannot open the story's worktree; accept from a shell with `flai accept S-nnnn`, or stop the dashboard and start it with `flai dashboard`, which mounts the repository at its own path. If the dashboard cannot push, it says the story was accepted locally and asks you to push from a shell. Every write is made by the bundled `flai`, so the files change exactly as they would from the terminal and appear in `git status` for you to commit.

Click a card for the item page: front matter, the rendered body, children, the transition history, blocked intervals, a link to the narrative, and buttons for the allowed moves, block, unblock, and adding a narrative log entry.

## Charts

Charts plots the flow metrics `flai stats` computes, so the numbers are the same in both places. Pick a window, an item type, and where it applies an epic. Cycle time shows one point per completed item with the p50 and p85 lines; burn-up shows scope against done; cumulative flow shows how many items sit in each state each day; time in state shows where each completed item spent its time and the share across all of them; throughput counts completions per week by nature; aging lists in-progress work against the p85 line; estimates compares estimated with actual hours. Every chart has a table view and follows the light or dark theme.

## Search

Search covers `design/` and `wip/` by default and `docs/` when you tick the box. Type an item ID, a title, or words from the body; each result shows where it is, its status, and a snippet around the match. Results open in the explorer.
