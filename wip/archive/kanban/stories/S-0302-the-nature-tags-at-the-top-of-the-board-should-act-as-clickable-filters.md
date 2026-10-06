---
id: S-0302
type: story
nature: improvement
title: The nature tags at the top of the board should act as clickable filters
status: done
owner: alex
created: 2026-10-06T23:17:17Z
updated: 2026-10-06T23:47:21Z
transitions:
  - to: ready
    at: 2026-10-06T23:17:18Z
    by: alex
  - to: in-progress
    at: 2026-10-06T23:39:58Z
    by: agent-S-0302
  - to: review
    at: 2026-10-06T23:46:32Z
    by: agent-S-0302
  - to: done
    at: 2026-10-06T23:47:21Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 416
  models:
    - model: claude-opus-5-5
      input: 122
      output: 42432
      cache_read: 3351099
      cache_write: 260767
      cost: 3.1667
    - model: claude-sonnet-5-5
      input: 10
      output: 2415
      cache_read: 113663
      cache_write: 38503
      cost: 0.1432
  strategic:
    - kind: planner
      seconds: 218
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 178
          output: 6732
          cache_read: 1021019
          cache_write: 78781
          cost: 0.2344
        - model: claude-opus-5-5
          input: 44
          output: 14058
          cache_read: 872416
          cache_write: 47925
          cost: 0.8392
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: alex
    at: 2026-10-06T23:17:17Z
  value: 25
  by: planner-S-0302
  at: 2026-10-06T23:20:39Z
forecast:
  duration: 20m
  delivery: 2026-10-06T23:54:00Z
  basis: "Its own forecast of 20m; 1st in the pull order with an in-progress limit of 3, behind S-0300."
  by: flai
  at: 2026-10-06T23:32:51Z
---
# S-0302 The nature tags at the top of the board should act as clickable filters

## Goal

The nature tags along the top of the board should work like filters. Their default coloring should represent 'filtered out' while a brighter color set should represent 'shown' with clicking each toggling them on or off. All natures should be toggled on (shown) by default.

## Acceptance criteria
- [x] When viewing a board, all natures are toggled on so that all natures of work items are displayed
- [x] When toggling a nature off, it reverts to its default coloring and cards with that nature are hidden
- [x] When toggling a nature on, it shows a highlighted version of its default color theme, and cards with that nature are shown
- [x] A user's choices are preserved in browser local storage so that navigating away from the board does not reset the nature filtering

## Tasks
- T-1128 Brighter nature tints for the shown state, in both themes
- T-1129 A nature filter store kept in localStorage, every nature shown by default
- T-1130 The legend's nature tags are toggle buttons, lit when shown
- T-1131 The board hides the cards of a nature toggled off
- T-1132 Document the nature filter in the dashboard design and the user guide

## Notes

### Planning

Touches:

- `flaiover/src` (declared, folder, kept): the operator declared it. Every file the work changes under it is named by a task (T-1128 to T-1131), so flai replaces it with those files in the story's claim; the new files are `flaiover/src/lib/boardnatures.svelte.ts` and its test (T-1129).
- `flai/cmd` (declared, folder, kept): the operator declared it, and the `cli` tag with it. No task touches it: the change is the dashboard's alone, and `flai board` in `flai/cmd/board.go` prints a text board with no legend to filter. Kept as declared; the plan's thread proposes dropping it, since as a folder touch it holds every ready story that touches `flai/cmd` while this one is in progress.
- `design/system/flaiover-dashboard.md` (design and co-change): the `/board` row of its views table and the Theme palette table describe the legend and the nature tints; `flai touches suggest` lists it among the files most changed with `flaiover/src` (17%).
- `docs/users/flaiover.md` (co-change and layout): its Board section explains the legend's colours and the type checkboxes; `flai touches suggest` lists it (13%).
- Left out of `flai touches suggest`'s list: `docs/users/flai.md`, `design/system/flai-cli.md`, `docs/users/flai-reference.md`, and the `flai/internal` files, which co-change with `flai/cmd` but describe or build the CLI, which this story does not change.

Forecast: 20m, adjusted from `flai forecast`'s 10m (median 94 s per unit of size over 7 done small improvement stories on claude-opus-5-5, times size 6). Doubled because the work is wider than its four criteria suggest: a shown tint per nature in both themes held by the contrast test, a store, the legend, the board's filter, and two documents. S-0141, the type checkboxes this mirrors, took 7.5 minutes with no new colours. Delivery moved by the same 10 minutes, to 2026-10-07T00:28Z, keeping flai's queue position (2nd in the pull order).

Cost of delay: 25.00 USD a week, as `flai cod` computes it from the operator's input, 10m of time lost per 168h cycle at 150 USD an hour; no reason to adjust it.
