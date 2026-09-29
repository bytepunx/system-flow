---
id: T-0594
type: task
nature: feature
title: The create screen starts an item in the lane it was asked from
status: done
parent: S-0167
owner: alex
created: 2026-09-29T23:33:31Z
updated: 2026-09-29T23:42:45Z
transitions:
  - to: ready
    at: 2026-09-29T23:33:34Z
    by: agent-S-0167
  - to: in-progress
    at: 2026-09-29T23:40:09Z
    by: agent-S-0167
  - to: done
    at: 2026-09-29T23:42:45Z
    by: agent-S-0167
stream: S-0167
tags: []
touches: [flaiover/src/routes/new, flaiover/src/lib/components/NewItemForm.svelte]
usage:
  source: log
  seconds: 156
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 31
      output: 9987
      cache_read: 2663625
      cache_write: 31000
      cost: 0.9806
---
# T-0594 The create screen starts an item in the lane it was asked from

## Work

- `/new?status=<lane>` preselects the starting lane in the form: `backlog`, `ready`, or `in-progress`; any other lane, or none, is `backlog`.
- After flai creates the item in backlog, move it along to the lane chosen through the existing move route; a refused move is shown with the created item's ID and nothing is undone.
- Component tests for the lane chosen, the fallback to backlog, and a refused move.

## Done when

- A story created from the ready lane's menu lands in ready; one from review, done, or cancelled lands in backlog.
- `npm run check`, lint, and the unit tests pass in `flaiover/`.

## Notes
