---
id: T-0213
type: task
nature: feature
title: "Board page: drag within backlog and ready with an insertion marker, nothing in unordered columns, and up and down controls with a keyboard shortcut"
status: done
parent: S-0057
owner: alex
created: 2026-09-19T06:41:25Z
updated: 2026-09-19T06:59:57Z
transitions:
  - to: ready
    at: 2026-09-19T06:50:27Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:50:27Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:59:57Z
    by: system-flow
stream: S-0057
tags: []
---

# T-0213 Board page: drag within backlog and ready with an insertion marker, nothing in unordered columns, and up and down controls with a keyboard shortcut

## Work
On the board page a story card dragged over another story card in the same backlog or ready column shows an insertion marker above or below it by the pointer's position, and the drop calls the order endpoint; the board refreshes to what flai wrote. In in-progress, review, done, and cancelled a drag within the column shows no marker and a drop does nothing. A drop on another column still moves the item, acceptance confirmation included. Epics and tasks are not reorderable. Story cards in backlog and ready get up and down buttons, outside the card's link, shown on hover and on keyboard focus, and Alt+ArrowUp and Alt+ArrowDown on a focused card do the same; the first card has no up and the last no down. Nothing is draggable and no controls show on a read-only board. Tests for the placement logic and the controls.

## Done when
- Placement logic (which call a drop or a button makes) is a tested function
- The card keeps being one link with no control inside it
- `make flaiover-test` and `make flaiover-build` pass

## Notes
