---
id: S-0048
type: story
nature: improvement
title: Parent epic indicator on story cards on the board
status: backlog
parent: E-0006
owner: alex
created: 2026-09-18T17:55:45Z
updated: 2026-09-18T17:55:45Z
transitions: []
tags: [dashboard]
touches: [flaiover]
---

# S-0048 Parent epic indicator on story cards on the board

## Goal
A story card on the kanban board shows which epic the story belongs to, in the bottom right corner of the card, so the designer can read the board by epic without opening each story.

## Acceptance criteria
- [ ] Every story card on `/board` shows its parent epic ID (for example `E-0006`) in the bottom right corner, on the same row as the nature and blocked flag, which stay on the left
- [ ] The indicator carries the epic's title as its tooltip and accessible name, so the ID is not the only way to tell epics apart
- [ ] A story with no parent shows no indicator and leaves no gap; epic cards show none; task cards, shown with "epics and tasks too", show their parent story the same way
- [ ] The indicator does not break the card: the card stays one link to the item, stays draggable, and long titles still wrap above the bottom row at every column width
- [ ] Component test on fixtures covering a story with a parent, a story without one, and a task; `design/system/flaiover-dashboard.md` and `docs/users` updated

## Tasks

## Notes
`GET /api/board` already returns `parent` on every card and the board page's `Card` type already declares it, so the ID needs no API change. The epic title is not in the card payload today: either add `parent_title` to the card in the API (and to the contract in `flaiover-dashboard.md`), or resolve it on the page from the epic cards already in the response. Decide when refining.

The card is a single `<a>` to the item, so the indicator cannot be a nested link to the epic; making it one means restructuring the card. Left out of scope unless the operator asks for it.

Not in the board `order` yet: E-0006 is pulled in story order and the operator has not said where this sits against S-0040 to S-0043.
