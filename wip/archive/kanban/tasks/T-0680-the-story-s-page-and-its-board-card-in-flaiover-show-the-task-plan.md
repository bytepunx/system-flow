---
id: T-0680
type: task
nature: experiment
title: The story's page and its board card in flaiover show the task plan
status: done
parent: S-0176
owner: arobson
created: 2026-10-01T11:42:11Z
updated: 2026-10-01T12:06:26Z
transitions:
  - to: ready
    at: 2026-10-01T11:57:27Z
    by: agent-S-0176
  - to: in-progress
    at: 2026-10-01T11:57:27Z
    by: agent-S-0176
  - to: done
    at: 2026-10-01T12:06:26Z
    by: agent-S-0176
stream: S-0176
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 539
  estimated: true
  models:
    - model: claude-haiku-4-5-20251001
      input: 123
      output: 3701
      cache_read: 680656
      cache_write: 35273
      cost: 0.1308
    - model: claude-opus-5-5
      input: 176
      output: 56551
      cache_read: 12497189
      cache_write: 178328
      cost: 4.7354
---
# T-0680 The story's page and its board card in flaiover show the task plan

## Work

- The story's page in flaiover shows each task as ready to start, waiting with the tasks it waits for (linked), in progress, done, or cancelled, and the plan as its layers, from the item's `plan` (contract in the notes).
- A story's card on the board shows its `tasks` summary when it has tasks: how many are in progress, ready, and waiting, and how many layers the plan has.

## Done when

- Component or unit tests cover the page's plan and the card's summary from fixtures in the contract's shape, and `npm run check`, the lint, and the tests touched pass in `flaiover/`.

## Notes

Built against the contract fixed in the narrative's `## Decisions`, while the flai task builds the data:

- Item (`item.get`): beside `item` and `children`, a story with tasks has `plan`: `{"tasks": [{"id": "T-0001", "state": "waiting", "after": ["T-0002"], "waiting_for": ["T-0002"]}], "layers": [["T-0002", "T-0003"], ["T-0001"]]}`. `state` is one of `ready`, `waiting`, `in-progress`, `done`, `cancelled`; `after` and `waiting_for` are omitted when empty; a task in a cycle is in no layer.
- Board card (`board.get`) of a story with tasks: `"tasks": {"ready": 1, "waiting": 2, "in_progress": 1, "done": 3, "layers": 3}`, omitted when it has none.
