---
id: T-0757
type: task
nature: feature
title: The edit form warns that a draft cannot go to ready and offers Finalize
status: done
parent: S-0201
owner: alex
created: 2026-10-03T07:37:09Z
updated: 2026-10-03T07:51:10Z
transitions:
  - to: ready
    at: 2026-10-03T07:37:32Z
    by: agent-S-0201
  - to: in-progress
    at: 2026-10-03T07:44:52Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T07:51:10Z
    by: agent-S-0201
stream: S-0201
tags: []
touches: [flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/routes/edit]
after: [T-0754]
usage:
  source: log
  seconds: 378
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 4034
      cache_read: 888825
      cache_write: 23075
      cost: 0.4066
---
# T-0757 The edit form warns that a draft cannot go to ready and offers Finalize

## Work

The item edit form (`ItemEditor.svelte`, on `/edit/...` and the item page) shows a draft story's state: a warning that the story is a draft and cannot be moved to ready until it is finalized, and a **Finalize** button there that posts to `/api/items/:id/finalize`, after which the warning goes. It waits for T-0754, which adds that route.

## Done when

- [ ] Editing a draft story shows the warning and a Finalize button; editing a finalized story shows neither
- [ ] Finalize from the form clears the draft through the route and the warning goes without a reload
- [ ] Vitest tests cover the warning and the button

## Notes
