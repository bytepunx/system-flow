---
id: T-0466
type: task
nature: feature
title: wait_for_work offers the first story not held, and moves warn on a held story
status: done
parent: S-0128
owner: alex
created: 2026-09-26T08:05:56Z
updated: 2026-09-26T08:14:53Z
transitions:
  - to: ready
    at: 2026-09-26T08:06:04Z
    by: agent-S-0128
  - to: in-progress
    at: 2026-09-26T08:12:56Z
    by: agent-S-0128
  - to: done
    at: 2026-09-26T08:14:53Z
    by: agent-S-0128
stream: S-0128
tags: []
touches: [flai/internal/mcpserver, flai/internal/workitem]
---
# T-0466 wait_for_work offers the first story not held, and moves warn on a held story

## Work

- `wait_for_work` (one project and a folder of projects) offers the first ready story that is not held; when every ready story is held it waits, and says so with `waiting_for: held` and each held story's reason.
- `flai move <story> in-progress` and MCP `item_move` warn on a held story and still move it.

## Done when

- Behaviour tests cover both; `go test ./internal/mcpserver ./internal/workitem ./cmd` passes.

## Notes
