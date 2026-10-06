---
id: T-1132
type: task
nature: improvement
title: Document the nature filter in the dashboard design and the user guide
status: backlog
parent: S-0302
owner: alex
created: 2026-10-06T23:20:19Z
updated: 2026-10-06T23:20:19Z
transitions: []
stream: S-0302
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1130, T-1131]
---
# T-1132 Document the nature filter in the dashboard design and the user guide

## Work

In `design/system/flaiover-dashboard.md`, say in the `/board` row of the views table that the legend's nature tags toggle which natures' cards are shown, every nature shown by default, the choice kept per browser in localStorage (`flaiover-board-natures`), and that counts and limits still count every card; and add the shown tints `nature-<nature>-shown` to the palette table under Theme, with their light and dark values and their use (a legend tag whose nature is shown, S-0302).

In `docs/users/flaiover.md`, section Board, tell the user that clicking a nature tag in the legend hides or shows that nature's cards, that a lit tag means shown and a dim one hidden, that every nature starts shown, and that the choice lasts across reloads and navigation in that browser. Waits for the legend and board tasks, so it describes what they built.

## Done when

- [ ] The design's `/board` row and palette table describe the nature filter and the shown tints with their values
- [ ] The user guide's Board section explains toggling a nature and that the choice is remembered
- [ ] The markdown lint passes on both files

## Notes
