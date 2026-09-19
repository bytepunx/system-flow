---
id: T-0181
type: task
nature: remediation
title: "MCP: inbox reports ready work and the designer's changes, wait_for_events reports events from the cursor, and a board tool"
status: done
parent: S-0058
owner: alex
created: 2026-09-19T02:46:21Z
updated: 2026-09-19T02:52:49Z
transitions:
  - to: ready
    at: 2026-09-19T02:50:01Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:50:01Z
    by: alex
  - to: done
    at: 2026-09-19T02:52:49Z
    by: alex
stream: S-0058
tags: []
touches: [flai/internal/mcpserver]
---

# T-0181 MCP: inbox reports ready work and the designer's changes, wait_for_events reports events from the cursor, and a board tool

## Work
In `flai/internal/mcpserver`: a cursor file per agent under `.flai-cache/mcp/<agent>.json` holding the last look, the boundary keys, and the board order last seen, written atomically; with no cursor, changes from the last 24 hours are reported. `inbox` gains `ready` (stories in the ready column in pull order, each with title, nature, and blocked; and `can_pull` from the in-progress limit), `changes` (from `workitem.Changes`, plus "pull order changed" when `order` differs), and advances the cursor; the `story` filter applies to threads as now. `wait_for_events` first returns any changes already behind the cursor, otherwise polls as now, and returns `events` with the changed paths, advancing the cursor. Add a `board` tool returning the shared board view with `all` to include epics and tasks. Update the server instructions and tool descriptions. Tests with the in-memory transports: a move to ready made between two calls is in the next `inbox` and not the one after; `ready` lists it until it is pulled; the agent's own `item_move` is not reported; a held `wait_for_events` returns the event within a poll; a `wait_for_events` called after the change returns at once; `board` matches the CLI's JSON.

## Done when
- The new tests pass with `-race`
- The existing MCP tests pass unchanged except where the output gained fields

## Notes
