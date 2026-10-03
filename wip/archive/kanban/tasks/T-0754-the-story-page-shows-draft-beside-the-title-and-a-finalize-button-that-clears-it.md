---
id: T-0754
type: task
nature: feature
title: "The story page shows [Draft] beside the title and a Finalize button that clears it"
status: done
parent: S-0201
owner: alex
created: 2026-10-03T07:36:49Z
updated: 2026-10-03T07:44:36Z
transitions:
  - to: ready
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: in-progress
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T07:44:36Z
    by: agent-S-0201
stream: S-0201
tags: []
touches: [flaiover/src/routes/items, "flaiover/src/routes/api/items/[id]/finalize"]
usage:
  source: log
  seconds: 424
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 41
      output: 11011
      cache_read: 2426169
      cache_write: 62987
      cost: 1.1098
---
# T-0754 The story page shows [Draft] beside the title and a Finalize button that clears it

## Work

On `/items/:id`, a story with `draft: true` shows a `[Draft]` indicator beside its title, in text, and a **Finalize** button. The button posts to a new route, `POST /api/items/:id/finalize`, which runs the hostapi write `item.finalize` with the item's id, and the page reloads the item from the API, so the indicator goes without a page reload. A refusal shows its error as the page's other writes do. The epic page's list of its stories marks a draft story the same way. It waits for nothing: the route names the op `item.finalize`, which T-0756 adds to flai, and the tests mock the API.

## Done when

- [ ] A draft story's page shows `[Draft]` beside the title and a Finalize button; a finalized story's shows neither
- [ ] Finalize posts to `/api/items/:id/finalize`, which writes `item.finalize`, and the indicator and button go without a reload
- [ ] An epic's page marks its draft stories in text
- [ ] Vitest tests cover the indicator, the button, the route, and the epic's list

## Notes
