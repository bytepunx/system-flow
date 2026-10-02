---
id: T-0677
type: task
nature: experiment
title: Tasks carry after, set by flai and MCP and checked by flai check
status: done
parent: S-0176
owner: arobson
created: 2026-10-01T11:42:10Z
updated: 2026-10-01T11:57:15Z
transitions:
  - to: ready
    at: 2026-10-01T11:42:31Z
    by: agent-S-0176
  - to: in-progress
    at: 2026-10-01T11:42:31Z
    by: agent-S-0176
  - to: done
    at: 2026-10-01T11:57:15Z
    by: agent-S-0176
stream: S-0176
tags: []
touches: [flai/internal/workitem, flai/internal/check, flai/cmd/edit.go, flai/cmd/task.go, flai/internal/mcpserver/items_write.go, design/system/work-hierarchy.md, docs/users/flai.md]
usage:
  source: log
  seconds: 884
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 198
      output: 63682
      cache_read: 14073111
      cache_write: 200815
      cost: 5.3325
---
# T-0677 Tasks carry after, set by flai and MCP and checked by flai check

## Work

- A task's `after` lists tasks of the same story that must be done before it starts. Item validation accepts it on tasks (task IDs only) and keeps a story's to story IDs.
- `flai task new --after`, `flai edit --after` and `--clear-after` on a task, and the MCP `item_new` and `item_edit` `after`, set it.
- `flai check` reports an entry that names no task, a task of another story, the task itself, and every cycle among a story's tasks.
- An older flai rejects `after` on a task, so the field list says so and a release that adds it raises `flai.minimum`. `design/system/work-hierarchy.md` documents a task's `after` and which flai the host must run first; `docs/users/flai.md` and the reference document the flags.

## Done when

- Behaviour tests cover each check finding, the flags, and the MCP parameters, and pass with `go test -race -short` on the packages touched.
- work-hierarchy.md and the user guide describe a task's `after`.

## Notes
