---
id: T-0787
type: task
nature: feature
title: The MCP tool plan starts a planner run for an epic or a story
status: done
parent: S-0208
owner: alex
created: 2026-10-04T00:44:06Z
updated: 2026-10-04T03:23:49Z
transitions:
  - to: ready
    at: 2026-10-04T00:44:32Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T03:16:40Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T03:23:49Z
    by: agent-S-0208
stream: S-0208
tags: []
touches: [flai/internal/mcpserver, flai/cmd]
after: [T-0786]
usage:
  source: log
  seconds: 429
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 41
      output: 11956
      cache_read: 2992921
      cache_write: 51921
      cost: 1.1526
---
# T-0787 The MCP tool plan starts a planner run for an epic or a story

## Work

- An MCP tool `plan` with `id` (an epic or a story) that starts a planner run through flai serve, as `agent_start` reaches it: a callback in `mcpserver.Options` wired in `flai/cmd` to what `flai plan` and `plan.run` use (T-0786), refused with the reason when the `plan` action is off, the item already has a run, or this flai mcp cannot start runs.
- Registered for a project and in a folder served whole, and the tool-list tests updated.

Waits for T-0786, which makes the serve call it reaches and shares `flai/cmd`.

## Done when

- [x] Tests show `plan` starts a run, and returns the refusals above
- [x] The tool lists in `server_test.go` and `folder_test.go` name `plan`

## Notes
