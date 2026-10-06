---
id: TH-0229
title: "Plan for S-0303: the type legend filters the board"
anchor:
  path: wip/kanban/stories/S-0303-the-type-specifiers-act-as-filters-for-the-board.md
  item: S-0303
status: open
participants: [planner-S-0303]
created: 2026-10-06T23:27:06Z
updated: 2026-10-06T23:27:06Z
---

# TH-0229 Plan for S-0303: the type legend filters the board

On wip/kanban/stories/S-0303-the-type-specifiers-act-as-filters-for-the-board.md.

## Entries

### 2026-10-06T23:27:06Z planner-S-0303
Plan for S-0303, four tasks in two layers:

Layer 1, run together, no shared path:
- T-1133 The board shows every card type by default: `boardtypes.svelte.ts` defaults go from stories alone to all three types; tests updated, including board tests that relied on the old default.
- T-1134 The legend's type entries are bordered toggles that show or hide their cards: `BoardLegend.svelte` type entries become `aria-pressed` buttons with the nature entries' border, toggling the existing `boardTypes` store (which already persists to localStorage); highlighted colouring from a new `typeHighlight` map in `cardcolour.ts`.

Layer 2, both after T-1133 and T-1134, no shared path between them:
- T-1135 The board drops its type checkboxes for the legend's toggles: remove `BoardTypes.svelte` and its test, drop it from the board page, add a board test that a legend click hides and shows cards.
- T-1136 The user guide and dashboard design say the type legend filters the board: `docs/users/flaiover.md` § Board and the `/board` row in `design/system/flaiover-dashboard.md`.

Assumptions:
1. The localStorage key `flaiover-board-types` stays, so a browser that already chose (stories alone, for anyone who has used the board) keeps its choice; only a browser with no stored choice starts with every type on. If you want everyone reset to all types on, the key changes instead; say so and T-1133 changes.
2. "Highlighted" means the entry filled with a tint of its own type colour (Tailwind opacity modifier on the existing `type-*` tokens, no new theme tokens), and "off" is the legend's look today plus the border. A different highlight, such as a solid fill, is a change to T-1134 only.
3. The nature entries stay plain labels; only types filter.
4. Forecast 20m (flai's 26m less the two declared folders, checked against S-0141's 8 minutes); cost of delay 25 USD a week from your 10m-per-cycle input.

Proposal: drop the declared touch `flai/cmd` and the `cli` tag. Nothing in the Go CLI knows the legend or filter (cards already carry their type from `/api/board`), no task touches it, and as a folder touch it holds every ready story touching `flai/cmd` while S-0303 is in progress. I left both in place since they are yours.
