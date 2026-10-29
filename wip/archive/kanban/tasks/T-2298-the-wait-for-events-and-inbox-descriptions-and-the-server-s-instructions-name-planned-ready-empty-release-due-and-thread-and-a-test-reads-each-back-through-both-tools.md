---
id: T-2298
type: task
nature: feature
title: The wait_for_events and inbox descriptions and the server's instructions name planned, ready-empty, release-due, and thread, and a test reads each back through both tools
status: done
parent: S-0535
owner: alex
created: 2026-10-28T21:39:01Z
updated: 2026-10-29T00:22:34Z
transitions:
  - to: ready
    at: 2026-10-29T00:08:58Z
    by: system-flow
  - to: in-progress
    at: 2026-10-29T00:08:59Z
    by: system-flow
  - to: done
    at: 2026-10-29T00:22:34Z
    by: agent-S-0535
stream: S-0535
tags: [flai, events, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/folder_test.go]
after: [T-2294, T-2295, T-2296, T-2297]
usage:
  source: log
  seconds: 815
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 58
      output: 19788
      cache_read: 3702487
      cache_write: 163933
      cost: 2.0584
---
# T-2298 The wait_for_events and inbox descriptions and the server's instructions name planned, ready-empty, release-due, and thread, and a test reads each back through both tools

## Work

The `wait_for_events` description (`flai/internal/mcpserver/server.go`), the `inbox` description (`inboxDescription` in `folder.go`), and the server's instructions in `server.go` name the kinds an agent sees and what to do on each. They must name the four new ones, so an agent knows what each carries.

- Name `planned`, `ready-empty`, `release-due`, and `thread` in both descriptions and in the server's instructions, with their fields. Say a story's agent may let the first three pass, and that `thread` is a thread now awaiting it.
- Add a test per kind that writes the event the way its writer does and reads it back through `wait_for_events` and through `inbox`, in the project mode and in the folder mode.

It waits for T-2294, T-2295, T-2296, and T-2297, whose events it reads back. T-2292 changed `server.go` before them.

## Done when

- Both descriptions and the instructions name each new kind and its fields.
- The read-back tests pass for each kind through both tools in both modes.
- `flai test` on the changed paths passes.

## Notes

S-0532 changes the same descriptions for the `agent` kind. Whichever story lands second merges its text with the other's.
