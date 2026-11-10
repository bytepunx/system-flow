---
id: T-2949
type: task
nature: remediation
title: The MCP tool stream_state takes context, decisions, and questions, with a test of a container's session writing them
status: in-progress
parent: S-0668
owner: alex
created: 2026-11-10T08:50:54Z
updated: 2026-11-10T09:55:10Z
transitions:
  - to: ready
    at: 2026-11-10T09:55:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T09:55:10Z
    by: system-flow
stream: S-0668
tags: [flai, mcp, narrative]
touches: [flai/internal/mcpserver/stream.go, flai/internal/mcpserver/stream_test.go, flai/internal/mcpserver/stream_sections_test.go]
after: [T-2944]
---
# T-2949 The MCP tool stream_state takes context, decisions, and questions, with a test of a container's session writing them

## Work

In `flai/internal/mcpserver/stream.go`, add to `StreamStateIn`:

- `context`, a string.
- `decisions`, a list of strings.
- `questions`, a list of strings.

Give each a `jsonschema` text that says what it does: context replaces, decisions and questions append, and a question above the threads block holds the move to review until it is answered. Pass them to the widened `SetStreamState`. Refuse a call that gives none of the five, naming all five. Rewrite `streamStateDescription` to say the tool writes these sections too and is how an agent whose `wip/` is read-only, in a container, records them.

This task holds the story's first criterion, I-0277's cause. Add a test in `flai/internal/mcpserver/stream_test.go`, beside `TestInStreamSession`, that reproduces it. A story's agent session writes `decisions`, `context`, and `questions` through `stream_state` and reads them back from the main checkout's narrative. That is the path a container's agent takes, over HTTP to the host's server (`flai/cmd/mcp.go`). Over stdio is enough when the HTTP server's test setup is too heavy, and the test's comment says so. Before T-2944 the test fails, since the tool has no such fields.

Waits for T-2944, whose `SetStreamState` it calls.

## Done when

- `flai test flai/internal/mcpserver/stream.go flai/internal/mcpserver/stream_test.go` passes.
- The new test writes all three sections through `stream_state` and finds them in the narrative, with its threads block unchanged.

## Notes
