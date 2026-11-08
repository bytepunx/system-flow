---
id: T-2628
type: task
nature: improvement
title: Save many work items under one hold of the items lock
status: done
parent: S-0599
owner: alex
created: 2026-11-07T23:47:56Z
updated: 2026-11-08T00:03:57Z
transitions:
  - to: ready
    at: 2026-11-07T23:52:07Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T23:52:10Z
    by: system-flow
  - to: done
    at: 2026-11-08T00:03:57Z
    by: agent-S-0599
stream: S-0599
tags: [cli]
touches: [flai/internal/workitem/itemlock.go, flai/internal/workitem/itemlock_test.go, flai/internal/workitem/repo.go]
usage:
  source: log
  seconds: 707
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 2420
      cache_read: 2536658
      cache_write: 41163
      cost: 0.7955
---
# T-2628 Save many work items under one hold of the items lock

## Work

Add `Repo.SaveAll(items []*Item)` beside `lockItems` in `flai/internal/workitem/itemlock.go`. It takes `.flai-cache/items.lock` once for the whole pass and, for each item, does what `Save` does under the lock: `mergeSaved`, `raiseUpdated`, an atomic write, and `readAs` set to what was written (ADR-0184).

- A part both writers changed refuses that item alone with its `*ConflictError`. The other items are still written.
- It returns which items were written and the error of each that was not.
- Keep `Save` as it is. Share the per-item body with it rather than copying it.

First layer: it waits for no task.

## Done when

- `SaveAll` writes every item it is given under one lock hold, and a test shows the lock is taken once, not once per item.
- A test shows an item another writer saved meanwhile, in a different part, is merged.
- A test shows an item another writer changed in the same part is refused with a `*ConflictError` naming it, and the rest are written.
- `flai test flai/internal/workitem` passes.

## Notes
