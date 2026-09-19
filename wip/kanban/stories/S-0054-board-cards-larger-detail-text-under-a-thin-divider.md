---
id: S-0054
type: story
nature: improvement
title: "Board cards: larger detail text under a thin divider"
status: backlog
parent: E-0006
owner: alex
created: 2026-09-19T02:05:28Z
updated: 2026-09-19T02:05:28Z
transitions: []
tags: [dashboard]
touches: [flaiover/src/lib/components]
---

# S-0054 Board cards: larger detail text under a thin divider

## Goal
The details at the bottom of a board card (nature, type, blocked flag, parent) are easier to read: set in a larger font, and separated from the title above by a thin horizontal line.

## Acceptance criteria
- [ ] The bottom row of `BoardCard.svelte` uses a larger font size than today's 10 pixels; the size is chosen with the operator from a rendered comparison and is the same for every item in the row, including the parent ID
- [ ] A thin horizontal line in the theme's line colour separates the title from the bottom row, with even spacing above and below, in the light and the dark theme
- [ ] The card still reads well in every column layout: the bottom row stays on one line where it did before or wraps cleanly, the parent stays bottom right, long titles wrap above the line, no card overflows, at desktop, tablet, and phone widths
- [ ] The ID and age row at the top is unchanged unless the operator asks for it when choosing the size
- [ ] Component test updated for the divider and the row; the board looked at in a browser at three widths in both themes; `docs/users/flaiover.md` updated only if what a card shows changes

## Tasks

## Notes
Raised by the operator on 2026-09-18: "use a larger font for the card details at the bottom and put a thin horizontal line to divide those from the card description/title visible on the card face."

The card is `flaiover/src/lib/components/BoardCard.svelte` since S-0048. The bottom row is `text-[10px] text-muted`; the title is the card's `text-xs` (12 pixels). Making the details as large as the title, or larger, changes which one the eye lands on, so show the operator two or three sizes rendered before settling.

S-0055 (colour coding) changes the same component. Do this one first or together with it; if they are pulled separately, the second rebases on the first.
