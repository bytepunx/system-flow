---
id: T-1580
type: task
nature: feature
title: "The board's card shows \"undoes main on N paths\" for a story whose record lists any, linking to its review page"
status: done
parent: S-0388
owner: alex
created: 2026-10-14T21:09:06Z
updated: 2026-10-14T23:29:23Z
transitions:
  - to: ready
    at: 2026-10-14T23:15:46Z
    by: agent-S-0388
  - to: in-progress
    at: 2026-10-14T23:15:47Z
    by: agent-S-0388
  - to: done
    at: 2026-10-14T23:29:23Z
    by: agent-S-0388
stream: S-0388
tags: [dashboard]
touches: [flaiover/src/lib/server/board.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/server/repo-channel.test.ts]
after: [T-1576]
usage:
  source: log
  seconds: 801
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 43
      output: 174
      cache_read: 1544157
      cache_write: 61309
      cost: 0.741
---
# T-1580 The board's card shows "undoes main on N paths" for a story whose record lists any, linking to its review page

## Work

- In `flaiover/src/lib/server/board.ts`, add `undoes_main` to `FlaiCard` and `Card`, and copy it in the mapping only when it is not empty, as `shared` is copied.
- In `flaiover/src/lib/components/BoardCard.svelte`, show "undoes main on N paths" in the details row after DRAFT (`data-testid="undoes-main"`). Say "1 path" for one.
- Make the badge lead to `/review/<id>`. The card is already a link, so follow how it links a story in review rather than nest an anchor.
- `flaiover/src/routes/board/+page.svelte` passes the card whole (`card={c}`), so it needs no change unless the card's type moves.
- It waits for T-1576, which puts `undoes_main` on `board.get`'s cards.

## Done when

- [ ] `BoardCard.svelte.test.ts` shows the badge for one path and for three, its link to the review page, and no badge for a card without `undoes_main`
- [ ] `flai test` on the three files passes

## Notes
