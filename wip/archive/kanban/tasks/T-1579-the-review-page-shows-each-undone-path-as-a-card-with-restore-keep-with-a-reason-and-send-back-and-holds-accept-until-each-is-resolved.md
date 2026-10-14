---
id: T-1579
type: task
nature: feature
title: The review page shows each undone path as a card with Restore, Keep with a reason, and Send back, and holds Accept until each is resolved
status: done
parent: S-0388
owner: alex
created: 2026-10-14T21:09:01Z
updated: 2026-10-14T23:29:08Z
transitions:
  - to: ready
    at: 2026-10-14T23:15:44Z
    by: agent-S-0388
  - to: in-progress
    at: 2026-10-14T23:15:45Z
    by: agent-S-0388
  - to: done
    at: 2026-10-14T23:29:08Z
    by: agent-S-0388
stream: S-0388
tags: [dashboard]
touches: [flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/WorktreeUncommitted.svelte]
after: [T-1577, T-1578]
usage:
  source: log
  seconds: 803
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 70
      output: 21369
      cache_read: 3102336
      cache_write: 114740
      cost: 1.7157
---
# T-1579 The review page shows each undone path as a card with Restore, Keep with a reason, and Send back, and holds Accept until each is resolved

## Work

- In `flaiover/src/lib/components/Review.svelte`, add `undoes_main` to the `Preview` type, in the shape S-0371's T-1479 gives it: each path, main's lines the branch removes, and the branch's version.
- Render one card per path (`data-testid="undoes-main"`) above the other blockers. Show the path, and main's removed lines beside the branch's version.
- **Restore.** POST `/api/items/<id>/restore` with the path, then fetch the preview again. The card clears when the path drops out.
- **Keep.** Ask for a one-line reason on the card. A kept card shows its reason and can be taken back before Accept.
- **Send back.** One action makes two calls, in order:
  1. POST `/api/threads` on the story, with text naming the undone paths.
  2. POST `/api/items/<id>/move` with `to: in-progress` and a reason naming the same paths.
- **Accept.** Keep it disabled (`canAccept`) while any card is neither restored nor kept.
  - Leave the undone paths' blockers out of `otherBlockers` in `WorktreeUncommitted.svelte`, so a kept path does not hold Accept.
  - Have `accept()` send the kept paths and their reasons as `keep_removals`.
- **Refusal.** A refused acceptance whose data carries `undoes_main` sets the same cards, as S-0378's T-1520 does for markdown findings.
- It waits for T-1577 and T-1578: it calls their endpoints, and its tests mock their bodies.

## Done when

- [ ] `Review.svelte.test.ts` shows the cards rendered from a preview's `undoes_main`
- [ ] It shows Restore calling the restore endpoint and fetching the preview again
- [ ] It shows Keep asking for a reason, and Accept sending `keep_removals`
- [ ] It shows Send back opening the thread, then moving the story
- [ ] It shows Accept disabled until every card is restored or kept, and a refused acceptance with `undoes_main` showing the cards
- [ ] `flai test` on the three files passes

## Notes
