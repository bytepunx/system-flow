---
id: S-0202
type: story
nature: feature
title: A card's right-click menu offers Finalize on a draft story and the card's other actions
status: in-progress
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T18:32:26Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:10Z
    by: alex
  - to: in-progress
    at: 2026-10-03T18:31:16Z
    by: agent-S-0202
tags: [dashboard]
touches: [flaiover/src/lib/components/BoardCard.svelte, flaiover/src/routes/board, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/lib/components/CardMenu.svelte.test.ts, flaiover/src/lib/cardmenu.ts, flaiover/src/lib/cardmenu.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/StoryAgent.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0201]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0202 A card's right-click menu offers Finalize on a draft story and the card's other actions

## Goal

The board has no per-card menu; a draft story has to be opened to be finalized. The designer asked for a right-click menu on cards with **Finalize** on draft cards.

## Acceptance criteria
- [ ] Right-clicking a card (and a long press on touch) opens a menu anchored to the card with: Open, Finalize (draft stories only), Start agent / Retry (as the story page offers them), Block / Unblock, and Cancel; each runs the same hostapi write the page runs and reflects the result on the card
- [ ] The menu is keyboard-reachable (a menu button on the focused card) and closes on Escape, blur, or a choice
- [ ] Finalize is only offered when the operator may finalize (the dashboard's user, not an agent view)
- [ ] `design/system/flaiover-dashboard.md` and the user guide describe the menu; tests cover opening, Finalize on a draft, and its absence on a finalized story

## Tasks
- T-0764 The card menu's entries come from one pure function, and the story page's agent actions share its rule
- T-0765 Right-click, long press, or a menu button on a card opens its menu on the board
- T-0766 The dashboard design and the user guide describe the card menu

## Notes
