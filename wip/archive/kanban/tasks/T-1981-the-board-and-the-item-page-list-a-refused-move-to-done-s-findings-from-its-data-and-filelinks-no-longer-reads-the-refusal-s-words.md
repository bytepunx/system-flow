---
id: T-1981
type: task
nature: improvement
title: The board and the item page list a refused move to done's findings from its data, and filelinks no longer reads the refusal's words
status: done
parent: S-0465
owner: alex
created: 2026-10-21T07:37:01Z
updated: 2026-10-21T08:07:49Z
transitions:
  - to: ready
    at: 2026-10-21T07:39:06Z
    by: agent-S-0465
  - to: in-progress
    at: 2026-10-21T08:01:24Z
    by: system-flow
  - to: done
    at: 2026-10-21T08:07:49Z
    by: agent-S-0465
stream: S-0465
tags: [flaiover]
touches: [flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", "flaiover/src/routes/api/items/[id]/move/+server.ts", flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, flaiover/src/lib/components/AcceptFindings.svelte]
after: [T-1979, T-1980]
usage:
  source: log
  seconds: 385
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 97
      output: 23663
      cache_read: 4793944
      cache_write: 197425
      cost: 2.6137
---
# T-1981 The board and the item page list a refused move to done's findings from its data, and filelinks no longer reads the refusal's words

## Work

This is the fix of I-0145. The board's move to done (`move` in `flaiover/src/routes/board/+page.svelte`) and the item page's acceptance (`flaiover/src/routes/items/[id]/+page.svelte`) post to `POST /api/items/:id/move` and show only `body.error` as text. `linkSegments` in `flaiover/src/lib/filelinks.ts` then opens each named file on the story's branch, unless the text holds `MAIN_CHECKOUT_REFUSAL`, flai's wording from `flai/cmd/accept_markdown.go`. With T-1979, `item.move` gives the refused object as the error's data, and `respond` spreads it into the body beside `error`.

- On the board and the item page, read `findings` from a refused move's body. When there are any, show them with `AcceptFindings` (T-1980) in place of the text notice, each file opened on its own checkout. Otherwise keep the notice as it is.
- Drop `MAIN_CHECKOUT_REFUSAL` and the story override it drives from `linkSegments`, and the test of it in `flaiover/src/lib/filelinks.test.ts`. No page picks a checkout from flai's words after this.
- Say in the doc comment of `flaiover/src/routes/api/items/[id]/move/+server.ts` that a refused move to done answers the refused object's `findings` or `undoes_main` beside `error`. Change the route itself only if `respond` does not already pass them.
- Tests: in `board.svelte.test.ts` and `item.svelte.test.ts`, a move to done refused with findings, one on the story's branch and one in the main checkout, lists both and opens each on its own checkout. Add a test to `filelinks.test.ts` that the refusal's words no longer change the checkout.
- Update the `/board`, `/items/:id`, and `POST /api/items/:id/move` rows of `design/system/flaiover-dashboard.md`, and § Fixing markdown that blocks an acceptance in `docs/users/flaiover.md`, so that they say the board and the item page list the findings as the review page does.

It waits for T-1979, which gives `item.move` the refused object as data, and for T-1980, which gives it `AcceptFindings`.

## Done when

- A move to done refused for lint the fix left, from the board or the item page, lists each finding and opens it on the checkout flai names, with no reading of the refusal's text.
- `linkSegments` has no text match on flai's wording.
- The tests above pass, and `flai test` passes on the changed paths.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` say so.

## Notes
