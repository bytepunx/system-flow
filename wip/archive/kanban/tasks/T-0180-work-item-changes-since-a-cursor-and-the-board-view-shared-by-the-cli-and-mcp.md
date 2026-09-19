---
id: T-0180
type: task
nature: remediation
title: Work item changes since a cursor, and the board view shared by the CLI and MCP
status: done
parent: S-0058
owner: alex
created: 2026-09-19T02:46:21Z
updated: 2026-09-19T02:50:00Z
transitions:
  - to: ready
    at: 2026-09-19T02:47:57Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:47:57Z
    by: alex
  - to: done
    at: 2026-09-19T02:50:00Z
    by: alex
stream: S-0058
tags: []
touches: [flai/internal/workitem, flai/internal/mcpserver, flai/cmd]
---

# T-0180 Work item changes since a cursor, and the board view shared by the CLI and MCP

## Work
In `flai/internal/workitem`, add the board view `flai board --json` builds today (columns, WIP limits, order, breaches) as a function both `flai/cmd/board.go` and the MCP server call, and a `Changes` function: given the items, a time, and the caller's name, the transitions and block or unblock intervals made after that time by anyone else, each with item ID, type, title, what happened, to what, by whom when known, and when; creation of an item counts as a change. Timestamps have second resolution, so the caller passes the keys already reported for the boundary second and gets keys back. Tests on items built in a temp project: a move by someone else is reported once, the caller's own move is not, a block and an unblock are reported, the boundary second neither loses nor repeats a change.

## Done when
- The tests pass and `flai board --json` output is unchanged (its existing test passes)
- `go test -race ./internal/workitem/... ./cmd/...` passes

## Notes
