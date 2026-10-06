---
id: S-0303
type: story
nature: improvement
title: The type specifiers act as filters for the board
status: ready
owner: alex
created: 2026-10-06T23:23:11Z
updated: 2026-10-06T23:33:30Z
transitions:
  - to: ready
    at: 2026-10-06T23:23:12Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, flaiover/src/lib/boardtypes.svelte.ts, flaiover/src/lib/boardtypes.svelte.test.ts, flaiover/src/lib/cardcolour.ts, flaiover/src/lib/components/BoardLegend.svelte, flaiover/src/lib/components/BoardLegend.svelte.test.ts, flaiover/src/lib/components/BoardTypes.svelte, flaiover/src/lib/components/BoardTypes.svelte.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 237
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 202
          output: 6327
          cache_read: 1129819
          cache_write: 68340
          cost: 0.2302
        - model: claude-opus-5-5
          input: 46
          output: 16533
          cache_read: 999390
          cache_write: 52129
          cost: 0.9478
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: alex
    at: 2026-10-06T23:23:11Z
  value: 25
  by: planner-S-0303
  at: 2026-10-06T23:26:54Z
forecast:
  duration: 20m
  delivery: 2026-10-06T23:55:00Z
  basis: "Its own forecast of 20m; 2nd in the pull order with an in-progress limit of 3, behind S-0300 and S-0302."
  by: flai
  at: 2026-10-06T23:33:30Z
---
# S-0303 The type specifiers act as filters for the board

## Goal

The type legend at the top of the board acts as a filter for card types on the board. When toggled on, their coloring changes from the default to a highlighted version. When toggled off, their coloring reverts to its original. Toggled on types are shown on the board while toggled off types are hidden. This replaces the checkbox set that currently provides this filtering.

## Acceptance criteria
- [ ] The type legend provides a border around each type so that it's clearer that it's a clickable entity (in the same way nature's have a border)
- [ ] When a type is clicked, this toggles whether that type is shown or hidden
- [ ] When a type is "on" and that type of card should be shown, change its coloring to a highlighted version of the default
- [ ] Store a user's toggle selections in the browser's local storage so that navigating away from the page does not reset their selections
- [ ] by default, all types are selected

## Tasks
- T-1133 The board shows every card type by default
- T-1134 The legend's type entries are bordered toggles that show or hide their cards
- T-1135 The board drops its type checkboxes for the legend's toggles
- T-1136 The user guide and dashboard design say the type legend filters the board

## Notes

### Planning

Touches, and where each came from:

- `flaiover/src`, `flai/cmd`: declared by the operator, kept. `flaiover/src` is a folder touch; flai narrows it in the story's claim to the files its tasks name (ADR-0096), so it holds nothing beyond them. `flai/cmd` is a folder touch no task names: the board's cards already carry their type from `/api/board` (`flaiover/src/lib/server/board.ts`), and nothing in the Go CLI knows the legend or the filter, so it is expected to stay unchanged. Kept because it is declared; while the story is in progress it holds every ready story touching `flai/cmd`, so the plan thread proposes dropping it with the `cli` tag.
- `flaiover/src/lib/boardtypes.svelte.ts`, `flaiover/src/lib/boardtypes.svelte.test.ts`: layout. The store behind the checkboxes (S-0141), key `flaiover-board-types`, whose default (stories alone) becomes every type shown.
- `flaiover/src/lib/components/BoardLegend.svelte`, `flaiover/src/lib/components/BoardLegend.svelte.test.ts`, `flaiover/src/lib/cardcolour.ts`: layout. The legend (S-0055) whose type entries become toggles, and the colour map it and the cards share, which gains the highlighted colouring.
- `flaiover/src/lib/components/BoardTypes.svelte`, `flaiover/src/lib/components/BoardTypes.svelte.test.ts`, `flaiover/src/routes/board/+page.svelte`, `flaiover/src/routes/board/board.svelte.test.ts`: layout. The checkbox set the goal replaces, the board page that renders it, and the board test that sets the store.
- `docs/users/flaiover.md`, `design/system/flaiover-dashboard.md`: design and co-change. The user guide's Board section and the dashboard design's `/board` row describe the checkboxes and their default; `flai touches suggest` ranks both among the files most often changed with `flaiover/src`.
- Not taken from `flai touches suggest`: its `flai/internal/hostapi`, `mcpserver`, `serve`, and `check` files co-change with the declared `flai/cmd`, not with the board legend, and `design/issues/summary.md` is written by the close-out, not planned.

Figures:

- Forecast: `flai forecast` gave 26m (84 s per unit of size times size 18: 5 criteria, 13 touches), delivery 2026-10-07T01:48Z. Adjusted to 20m: two of the 13 touches are the declared folders, which add no work beyond the files listed, and S-0141, which built the store and checkboxes this story reworks, went from in-progress to review in 8 minutes; this story is about twice its work over two layers. Delivery moved back by the same 6 minutes, to 2026-10-07T01:42Z, keeping flai's queue position (6th in the pull order, in-progress limit 3).
- Cost of delay: `flai cod` gave 25.00 USD a week from the operator's 10m lost per 168h cycle at 150 USD an hour. It stands: the operator's inputs are the whole basis and nothing in the plan changes them.
