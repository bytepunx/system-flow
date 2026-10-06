---
id: T-1136
type: task
nature: improvement
title: The user guide and dashboard design say the type legend filters the board
status: done
parent: S-0303
owner: alex
created: 2026-10-06T23:26:18Z
updated: 2026-10-06T23:54:22Z
transitions:
  - to: ready
    at: 2026-10-06T23:50:46Z
    by: agent-S-0303
  - to: in-progress
    at: 2026-10-06T23:50:46Z
    by: agent-S-0303
  - to: done
    at: 2026-10-06T23:54:22Z
    by: agent-S-0303
stream: S-0303
tags: [dashboard, docs]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-1133, T-1134]
usage:
  source: log
  seconds: 216
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 23
      output: 7792
      cache_read: 746458
      cache_write: 43539
      cost: 0.5852
---
# T-1136 The user guide and dashboard design say the type legend filters the board

## Work

In `docs/users/flaiover.md` § Board, replace the paragraph's account of the checkboxes ("Checkboxes above the columns choose which types of card are shown…", "whichever boxes are ticked") with the legend: click a type in it to hide or show its cards, a highlighted type is shown, every type is shown until you change it, and the choice is remembered in this browser. Mention under **Colours** that the type entries of the legend are also the filter. In `design/system/flaiover-dashboard.md`, rewrite the `/board` row's clause on the checkboxes (`$lib/boardtypes.svelte.ts`, `BoardTypes.svelte`, S-0141) to name the legend's type toggles, `aria-pressed`, the `typeHighlight` map in `$lib/cardcolour.ts`, the default of every type shown, and S-0303, and its legend clause (`BoardLegend.svelte`, S-0055) to match.

Waits for T-1133 and T-1134, whose defaults and colouring it describes. It shares no path with T-1135, so the two form the second layer together.

## Done when

- [ ] Neither document mentions type checkboxes or stories alone as the default.
- [ ] Both describe the legend toggles, the highlighted colouring, the default of every type shown, and the choice kept per browser, as built.
- [ ] The markdown lint and `flai check --strict` pass.

## Notes

Drafted by the planner for S-0303.
