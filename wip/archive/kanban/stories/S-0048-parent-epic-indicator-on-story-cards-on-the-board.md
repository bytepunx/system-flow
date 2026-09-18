---
id: S-0048
type: story
nature: improvement
title: Parent epic indicator on story cards on the board
status: done
parent: E-0006
owner: alex
created: 2026-09-18T17:55:45Z
updated: 2026-09-18T20:55:46Z
transitions:
  - to: ready
    at: 2026-09-18T20:39:52Z
    by: alex
  - to: in-progress
    at: 2026-09-18T20:40:27Z
    by: alex
  - to: review
    at: 2026-09-18T20:49:05Z
    by: alex
  - to: done
    at: 2026-09-18T20:55:46Z
    by: alex
tags: [dashboard]
touches: [flaiover]
---

# S-0048 Parent epic indicator on story cards on the board

## Goal
A story card on the kanban board shows which epic the story belongs to, in the bottom right corner of the card, so the designer can read the board by epic without opening each story.

## Acceptance criteria
- [x] Every story card on `/board` shows its parent epic ID (for example `E-0006`) in the bottom right corner, on the same row as the nature and blocked flag, which stay on the left
- [x] The indicator carries the epic's title as its tooltip and accessible name, so the ID is not the only way to tell epics apart
- [x] A story with no parent shows no indicator and leaves no gap; epic cards show none; task cards, shown with "epics and tasks too", show their parent story the same way
- [x] The indicator does not break the card: the card stays one link to the item, stays draggable, and long titles still wrap above the bottom row at every column width
- [x] Component test on fixtures covering a story with a parent, a story without one, and a task; `design/system/flaiover-dashboard.md` and `docs/users` updated

## Tasks
- T-0159 Board API returns parent_title on every card
- T-0160 BoardCard component with the parent indicator bottom right, used by the board page
- T-0161 Docs, all checks, and a look at the real board at several widths

## Notes
`GET /api/board` already returns `parent` on every card and the board page's `Card` type already declares it, so the ID needs no API change. The epic title is not in the card payload today: either add `parent_title` to the card in the API (and to the contract in `flaiover-dashboard.md`), or resolve it on the page from the epic cards already in the response. Decide when refining.

The card is a single `<a>` to the item, so the indicator cannot be a nested link to the epic; making it one means restructuring the card. Left out of scope unless the operator asks for it.

Not in the board `order` yet: E-0006 is pulled in story order and the operator has not said where this sits against S-0040 to S-0043.

Decided when pulled, 2026-09-18: the parent's title comes from the API as `parent_title`, and the card became `BoardCard.svelte` so it could have a component test.

Verification, 2026-09-18. The dev server from the story branch was run against this repository with a throwaway token and the board opened in a browser at 1440, 768, and 390 pixels, with "epics and tasks too" off and on. Measured on every card: 11 stories each show their epic, 3 tasks show `S-0048`, 4 epics show nothing; the indicator's right and bottom edges are 9 pixels inside the card's (its padding and border) on all 14; the title ends above the bottom row on all 18; the bottom row stays on one line; no card overflows, at card widths from the six-column layout down to a single 325 pixel column. Tooltips read, for example, `E-0004 User-facing documentation`. Dragging was not performed in the browser, because the board was this repository's real one; the component test covers the draggable attribute and the drag handler. At 390 pixels the page scrolls sideways because of the navigation bar in the root layout, not the board; recorded as I-0018 and left alone. Screenshots are in the session scratchpad, not the repository.
