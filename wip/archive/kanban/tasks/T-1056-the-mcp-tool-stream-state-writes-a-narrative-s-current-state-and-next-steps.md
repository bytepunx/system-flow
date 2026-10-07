---
id: T-1056
type: task
nature: improvement
title: The MCP tool stream_state writes a narrative's Current state and Next steps
status: done
parent: S-0271
owner: alex
created: 2026-10-06T22:49:51Z
updated: 2026-10-07T07:03:54Z
transitions:
  - to: ready
    at: 2026-10-07T06:57:07Z
    by: agent-S-0271
  - to: in-progress
    at: 2026-10-07T06:57:08Z
    by: agent-S-0271
  - to: done
    at: 2026-10-07T07:03:54Z
    by: agent-S-0271
stream: S-0271
tags: [mcp]
touches: [flai/internal/mcpserver/folder.go, flai/internal/mcpserver/stream.go, flai/internal/mcpserver/stream_test.go, flai/internal/mcpserver/folder_test.go, flai/internal/mcpserver/server_test.go]
after: [T-1054]
usage:
  source: log
  seconds: 406
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 15213
      cache_read: 2263177
      cache_write: 78416
      cost: 1.2167
---
# T-1056 The MCP tool stream_state writes a narrative's Current state and Next steps

## Work

Criterion 2, over MCP.

- Add `stream_state` in a new `flai/internal/mcpserver/stream.go`, modelled on `criteria.go`. It takes `story`, `current`, `next`, and `project`, and calls T-1054's function as the connected agent.
- Register it in `folder.go` beside `criteria_tick`.
- Answer with the story, the narrative's path, and the sections written.
- Return each refusal as a tool error that names the reason and what to do: no narrative, the wrong state, or a lint finding.
- Write the description so that it says the tool replaces the two sections, leaves the rest alone, and appends nothing to the log.

## Done when

- [ ] Tests in `flai/internal/mcpserver/stream_test.go` show the following:
  - A write that leaves the other sections unchanged.
  - Each text alone.
  - Each refusal.
- [ ] The tool is listed by the server in the existing tool-list test, if there is one.
- [ ] `scripts/flai-test.sh` passes.

## Notes

Drafted by the planner for S-0271. The tool name follows the story's goal. The tick half of the story is the existing `criteria_tick`, unless the plan thread decides otherwise.
