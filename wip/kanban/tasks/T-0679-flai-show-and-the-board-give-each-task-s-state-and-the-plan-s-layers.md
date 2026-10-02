---
id: T-0679
type: task
nature: experiment
title: flai show and the board give each task's state and the plan's layers
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
touches: [flai/internal/taskplan, flai/cmd/items.go, flai/cmd/board.go, flai/internal/workitem/rules.go, flai/internal/workitem/board.go, flai/internal/mcpserver, flai/internal/hostapi, docs/users/flai.md, design/system/flai-cli.md]
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
# T-0679 flai show and the board give each task's state and the plan's layers

## Work

- A story's task plan, computed in one place from its tasks' `after` and status: each task's state, `ready` (to start), `waiting` (with the tasks it waits for), `in-progress`, `done`, or `cancelled`; and the layers, the open and done tasks grouped by how many `after` steps lie before them, so each layer can run at once once the one before it is done.
- `flai show` prints the plan for a story with tasks; `flai show --json`, the MCP `item_get`, and the host's `item.get` carry it as `plan` (the contract below).
- `flai board` and the host's `board.get` give a story's card a `tasks` summary (the contract below).
- `flai move T-nnnn in-progress` warns, and does not refuse, while a task it waits for is not done, as a story's `after` does.

## Done when

- Behaviour tests cover the states, the layers (with a cycle, which puts no task in a layer), the warning, and the JSON, and pass with `go test -race -short` on the packages touched.
- `docs/users/flai.md` and `design/system/flai-cli.md` describe what `flai show` and `flai board` print.

## Notes

The contract with the flaiover task (fixed in the narrative's `## Decisions` so both run at once):

- Item JSON: beside `item` and `children`, a story with tasks has `plan`: `{"tasks": [{"id": "T-0001", "state": "waiting", "after": ["T-0002"], "waiting_for": ["T-0002"]}], "layers": [["T-0002", "T-0003"], ["T-0001"]]}`. `after` and `waiting_for` are omitted when empty; `layers` lists task IDs in ID order within a layer; a task in a cycle is in no layer.
- Board card of a story with tasks: `"tasks": {"ready": 1, "waiting": 2, "in_progress": 1, "done": 3, "layers": 3}`, omitted when it has none.
