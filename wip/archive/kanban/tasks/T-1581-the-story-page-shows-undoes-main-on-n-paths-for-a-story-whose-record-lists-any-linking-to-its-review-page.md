---
id: T-1581
type: task
nature: feature
title: "The story page shows \"undoes main on N paths\" for a story whose record lists any, linking to its review page"
status: done
parent: S-0388
owner: alex
created: 2026-10-14T21:09:09Z
updated: 2026-10-14T23:29:29Z
transitions:
  - to: ready
    at: 2026-10-14T23:15:48Z
    by: agent-S-0388
  - to: in-progress
    at: 2026-10-14T23:15:49Z
    by: agent-S-0388
  - to: done
    at: 2026-10-14T23:29:29Z
    by: agent-S-0388
stream: S-0388
tags: [dashboard]
touches: [flaiover/src/lib/server/repo.ts, "flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts", flaiover/src/lib/server/repo-forget.test.ts]
after: [T-1576]
usage:
  source: log
  seconds: 820
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 36
      output: 11112
      cache_read: 1613211
      cache_write: 59665
      cost: 0.8921
---
# T-1581 The story page shows "undoes main on N paths" for a story whose record lists any, linking to its review page

## Work

- In `flaiover/src/lib/server/repo.ts`, have `itemById` pass `item.get`'s `undoes_main` through when it is not empty, as it passes `expected_cost`.
- In `flaiover/src/routes/items/[id]/+page.svelte`, show "undoes main on N paths" beside the HELD badge (`data-testid="undoes-main"`), as a link to `/review/<id>`. Say "1 path" for one.
- It waits for T-1576, which puts `undoes_main` on `item.get`.

## Done when

- [ ] `item.svelte.test.ts` shows the badge and its link for a story with `undoes_main`, and none without it
- [ ] `flai test` on the three files passes

## Notes
