---
id: T-1131
type: task
nature: improvement
title: The board hides the cards of a nature toggled off
status: done
parent: S-0302
owner: alex
created: 2026-10-06T23:20:13Z
updated: 2026-10-06T23:44:17Z
transitions:
  - to: ready
    at: 2026-10-06T23:42:09Z
    by: agent-S-0302
  - to: in-progress
    at: 2026-10-06T23:42:09Z
    by: agent-S-0302
  - to: done
    at: 2026-10-06T23:44:17Z
    by: agent-S-0302
stream: S-0302
tags: [dashboard]
touches: [flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts]
after: [T-1129]
usage:
  source: log
  seconds: 128
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 7258
      cache_read: 573197
      cache_write: 44604
      cost: 0.5417
---
# T-1131 The board hides the cards of a nature toggled off

## Work

In `flaiover/src/routes/board/+page.svelte`, extend `cards(state)` so a lane shows a card only when its type is ticked (S-0141) and `boardNatures.isShown(card.nature)`. As for the type filter, the WIP count, reordering, and the lane's counts by type (S-0256) still count every card the lane holds, so hiding a nature never changes a limit or the pull order; update the comment above `cards` to say so.

Add tests to `flaiover/src/routes/board/board.svelte.test.ts`: with nothing stored every nature's cards show; with a nature stored as off its cards are absent and the others present, and the WIP count is unchanged; and the filter combines with the type checkboxes. Waits for the store task, whose `boardNatures` it reads; it has no path in common with the legend task, so the two run together.

## Done when

- [ ] Cards of a nature toggled off are hidden in every lane, and shown again when it is toggled on
- [ ] WIP counts and lane counts are unchanged by the filter
- [ ] The board tests cover it and pass with `scripts/flaiover-test.sh`

## Notes
