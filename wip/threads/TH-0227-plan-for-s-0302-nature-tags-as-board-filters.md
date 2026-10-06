---
id: TH-0227
title: "Plan for S-0302: nature tags as board filters"
anchor:
  path: wip/kanban/stories/S-0302-the-nature-tags-at-the-top-of-the-board-should-act-as-clickable-filters.md
  item: S-0302
status: open
participants: [planner-S-0302]
created: 2026-10-06T23:20:52Z
updated: 2026-10-06T23:20:52Z
---

# TH-0227 Plan for S-0302: nature tags as board filters

On wip/kanban/stories/S-0302-the-nature-tags-at-the-top-of-the-board-should-act-as-clickable-filters.md.

## Entries

### 2026-10-06T23:20:52Z planner-S-0302
Plan for S-0302, five tasks in three layers.

Layer 1 (run together, no paths in common):
- T-1128 Brighter nature tints for the shown state, in both themes: `layout.css`, `cardcolour.ts`, `theme.test.ts`.
- T-1129 A nature filter store kept in localStorage, every nature shown by default: new `lib/boardnatures.svelte.ts` and its test, modelled on S-0141's `boardtypes.svelte.ts`.

Layer 2 (run together):
- T-1130 The legend's nature tags are toggle buttons, lit when shown (after T-1128, T-1129): `BoardLegend.svelte` and its test.
- T-1131 The board hides the cards of a nature toggled off (after T-1129): `routes/board/+page.svelte` and `board.svelte.test.ts`.

Layer 3:
- T-1132 Document the nature filter in the dashboard design and the user guide (after T-1130, T-1131): `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`.

Figures: forecast 20m, raised from flai's 10m because the work adds tints in two themes held by the contrast test, plus a store, the legend, the filter, and docs. Delivery 2026-10-07T00:28Z. Cost of delay 25 USD a week, from your 10m-per-cycle input, unchanged.

Assumptions (correct me where any is wrong):
1. Only the legend's tags change colour. Cards keep their pastel background whichever natures are shown.
2. "Brighter" means a more saturated tint in the light theme and a lighter one in the dark theme, held to the same ink-text contrast as the card tints. The state is also carried by `aria-pressed`, not by colour alone.
3. As with the type checkboxes, WIP counts, lane counts and reordering still count every card, so hiding a nature never changes a limit or the pull order. The nature filter combines with the type filter: a card shows only when both its type and its nature are on.
4. A card whose nature flai does not know is never hidden.

Proposal: drop the declared touch `flai/cmd` and the `cli` tag. No task changes the CLI: `flai board` prints a text board with no legend. As a folder touch, `flai/cmd` holds every ready story that touches `flai/cmd` while this one is in progress. I kept it because it is yours to remove. Say the word, or remove it yourself, and I will not touch it otherwise.
