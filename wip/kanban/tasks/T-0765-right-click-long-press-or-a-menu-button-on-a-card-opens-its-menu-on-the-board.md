---
id: T-0765
type: task
nature: feature
title: Right-click, long press, or a menu button on a card opens its menu on the board
status: ready
parent: S-0202
owner: alex
created: 2026-10-03T18:32:22Z
updated: 2026-10-03T18:32:54Z
transitions:
  - to: ready
    at: 2026-10-03T18:32:54Z
    by: agent-S-0202
stream: S-0202
tags: []
touches: [flaiover/src/lib/components/CardMenu.svelte, flaiover/src/lib/components/CardMenu.svelte.test.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/routes/board]
after: [T-0764]
---

# T-0765 Right-click, long press, or a menu button on a card opens its menu on the board

## Work

Add `CardMenu.svelte`, a menu like `LaneMenu.svelte` anchored to the card, listing `cardMenu`'s entries and then the lane's entries under a divider, so a right click on a card still reaches what the lane menu offered there (S-0167). Right click on a card, the context-menu key or Shift+F10 on a focused card, a long press on touch, and a menu button beside the focused or hovered card (like `CardReorder`) open it; Escape, focus leaving it, a click elsewhere, or a choice close it. In `routes/board/+page.svelte` each entry runs the write the item page runs: `POST /api/items/:id/finalize`, `.../block` with a reason, `.../unblock`, `CancelConfirm` then the move to cancelled, and `POST /api/items/:id/agent`; then the board and the agents are asked again and a notice says what happened. A read-only board offers no menu. Waits for T-0764, whose `cardMenu` it renders.

## Done when

- Component tests cover opening the menu by right click and by the menu button, Finalize on a draft card posting to the finalize route, no Finalize on a finalized story, and Escape closing it.
- The board's existing tests still pass.

## Notes
