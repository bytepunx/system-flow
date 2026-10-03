---
id: S-0202
type: story
nature: feature
title: A card's right-click menu offers Finalize on a draft story and the card's other actions
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T19:13:17Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:10Z
    by: alex
  - to: in-progress
    at: 2026-10-03T18:31:16Z
    by: agent-S-0202
  - to: review
    at: 2026-10-03T18:50:46Z
    by: agent-S-0202
  - to: done
    at: 2026-10-03T19:13:17Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/lib/components/BoardCard.svelte, flaiover/src/routes/board, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/CardMenu.svelte, flaiover/src/lib/components/CardMenu.svelte.test.ts, flaiover/src/lib/components/Menu.svelte, flaiover/src/lib/components/LaneMenu.svelte, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/lib/cardmenu.ts, flaiover/src/lib/cardmenu.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/StoryAgent.svelte, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md, design/issues/summary.md]
after: [S-0201]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1197
  models:
    - model: claude-opus-5-5
      input: 228
      output: 84034
      cache_read: 13177640
      cache_write: 337691
      cost: 6.5882
    - model: claude-sonnet-5
      input: 46
      output: 12854
      cache_read: 1348618
      cache_write: 82555
      cost: 0.6047
---
# S-0202 A card's right-click menu offers Finalize on a draft story and the card's other actions

## Goal

The board has no per-card menu; a draft story has to be opened to be finalized. The designer asked for a right-click menu on cards with **Finalize** on draft cards.

## Acceptance criteria
- [x] Right-clicking a card (and a long press on touch) opens a menu anchored to the card with: Open, Finalize (draft stories only), Start agent / Retry (as the story page offers them), Block / Unblock, and Cancel; each runs the same hostapi write the page runs and reflects the result on the card
- [x] The menu is keyboard-reachable (a menu button on the focused card) and closes on Escape, blur, or a choice
- [x] Finalize is only offered when the operator may finalize (the dashboard's user, not an agent view)
- [x] `design/system/flaiover-dashboard.md` and the user guide describe the menu; tests cover opening, Finalize on a draft, and its absence on a finalized story

## Tasks
- T-0764 The card menu's entries come from one pure function, and the story page's agent actions share its rule
- T-0765 Right-click, long press, or a menu button on a card opens its menu on the board
- T-0766 The dashboard design and the user guide describe the card menu

## Notes

- The long press is verified with simulated touch pointer events in jsdom (`flaiover/src/routes/board/cardmenu.svelte.test.ts`), not on a touch device: a browser that starts an HTML drag on a long press may cancel the press first.
- "The dashboard's user, not an agent view" is a writable board: the dashboard writes as the designer (ADR-0077), and a read-only board offers no menu.
