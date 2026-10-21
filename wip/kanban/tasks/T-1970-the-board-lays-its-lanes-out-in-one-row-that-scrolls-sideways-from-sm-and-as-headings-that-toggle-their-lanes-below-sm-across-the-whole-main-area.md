---
id: T-1970
type: task
nature: improvement
title: The board lays its lanes out in one row that scrolls sideways from sm and as headings that toggle their lanes below sm, across the whole main area
status: backlog
parent: S-0462
owner: alex
created: 2026-10-21T07:04:56Z
updated: 2026-10-21T07:04:56Z
transitions: []
stream: S-0462
tags: [dashboard]
touches: [flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/+layout.svelte]
after: [T-1969]
---
# T-1970 The board lays its lanes out in one row that scrolls sideways from sm and as headings that toggle their lanes below sm, across the whole main area

## Work

- In `flaiover/src/routes/board/+page.svelte`, replace `grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-6` as follows:
  - From `sm`, a flex row of lanes with a minimum width and no wrap, inside a container that scrolls sideways (`overflow-x-auto`, scroll snap to lanes). Each lane's heading is sticky while its cards scroll.
  - Below `sm`, each lane's `<h2>` holds a button with `aria-expanded` and `aria-controls` on its card list. The button carries the name, the counts line, and the WIP count, and toggles the lane through `lanecollapse.svelte.ts`. A closed lane hides its cards and keeps its menu.
- In `flaiover/src/routes/+layout.svelte`, let `/board` take the whole main area instead of `max-w-6xl`.
- Check that the lane menu, the card menu, reordering, the lane-move dialog, the merge queue, and the publish banner work in both layouts, and that a card moved into a closed lane updates its counts.
- Extend `board.svelte.test.ts`: the toggle buttons and their `aria-expanded`, the default open lanes, a move into a closed lane, and the menus in a closed lane.
- Waits for T-1969, the open-lanes state it reads.

## Done when

`board.svelte.test.ts` and the board's other tests pass, and the board shows one row of lanes from `sm` and toggling headings below it.

## Notes
