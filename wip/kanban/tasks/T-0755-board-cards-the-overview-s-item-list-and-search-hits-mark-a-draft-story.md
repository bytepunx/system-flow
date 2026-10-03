---
id: T-0755
type: task
nature: feature
title: Board cards, the overview's item list, and search hits mark a draft story
status: done
parent: S-0201
owner: alex
created: 2026-10-03T07:36:59Z
updated: 2026-10-03T07:51:10Z
transitions:
  - to: ready
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: in-progress
    at: 2026-10-03T07:44:51Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T07:51:10Z
    by: agent-S-0201
stream: S-0201
tags: []
touches: [flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/server/board.ts, flaiover/src/lib/server/search.ts, flaiover/src/routes/search, flaiover/src/routes/+page.svelte]
after: [T-0753]
usage:
  source: log
  seconds: 379
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5566
      cache_read: 1226317
      cache_write: 31837
      cost: 0.561
---
# T-0755 Board cards, the overview's item list, and search hits mark a draft story

## Work

The dashboard marks a draft story wherever it lists one, in text, not colour alone: a board card (`BoardCard.svelte`, fed by `lib/server/board.ts`), the overview's list of items (`routes/+page.svelte`, fed by `/api/items`), and a search hit (`routes/search`, fed by `lib/server/search.ts`). It waits for T-0753, which puts `draft` on flai's board cards and search hits.

## Done when

- [ ] A draft story's board card, overview row, and search hit show a `Draft` marker in text; a finalized story's show none
- [ ] Vitest tests cover the card marker, the overview, and the search hit

## Notes
