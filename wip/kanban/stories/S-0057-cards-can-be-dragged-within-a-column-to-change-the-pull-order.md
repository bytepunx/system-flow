---
id: S-0057
type: story
nature: feature
title: Cards can be dragged within a column to change the pull order
status: ready
parent: E-0006
owner: alex
created: 2026-09-19T02:05:29Z
updated: 2026-09-19T05:35:38Z
transitions:
  - to: ready
    at: 2026-09-19T05:35:38Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover, flai/cmd]
---

# S-0057 Cards can be dragged within a column to change the pull order

## Goal
The designer reorders work by dragging a card up or down within its column, and that order is the pull order agents follow.

## Acceptance criteria
- [ ] Cards in the backlog and ready columns appear in the board's pull order (`order` in `wip/kanban/board.md`), with items not in the list after those that are; today the dashboard ignores the order when laying out a column
- [ ] Dragging a story card to a new position within backlog or within ready changes the pull order, the card stays where it was dropped after the board refreshes, and an agent reading `flai board` sees the same order
- [ ] The write goes through flai (ADR-0016): a command sets an item's position in the order (before or after another item, or to the top or bottom), keeps the rule that ready stories come before backlog stories, and the dashboard calls it; `board.md` is not edited by the dashboard directly
- [ ] Dragging within a column that has no meaningful order (in-progress, review, done, cancelled) does nothing and does not look as if it worked; dragging between columns still moves the item as it does today, including the acceptance confirmation on done
- [ ] It works with a mouse, and there is a keyboard or button alternative for moving a card up and down; on a read-only board nothing is draggable
- [ ] Tests for the command's ordering rules, the endpoint, and the board's sorting; the board tried in a browser; `design/system/workflow.md`, `flai-cli.md`, `flaiover-dashboard.md`, and the user documentation updated

## Tasks

## Notes
Raised by the operator on 2026-09-18: "I would like for cards to be draggable within a column if possible to allow for re-prioritization."

It is possible, with two pieces of groundwork found when queuing this. The pull order exists (`order` in `board.md`; `flai move` appends a story when it becomes ready and removes it when it starts or is cancelled, see `design/system/workflow.md`) but nothing can change it except a hand edit: `flai` has no ordering command. And the dashboard does not use it: `flaiover/src/lib/server/board.ts` returns cards in the reader's order and the page only prints the list as text above the columns. So the story is a flai command, the columns sorted by the order, and then the drag.

Several stories created on 2026-09-18 were never added to `order` because the operator had not placed them (S-0047, S-0052 to S-0057 among them); once columns sort by the order, where unlisted items go becomes visible, hence the first criterion.

Tasks and epics are not in the pull order. Whether tasks should be reorderable within their story is a separate question and out of scope here.

Shares the card's drag handlers with `BoardCard.svelte` (S-0048); S-0054 and S-0055 touch the same component.
