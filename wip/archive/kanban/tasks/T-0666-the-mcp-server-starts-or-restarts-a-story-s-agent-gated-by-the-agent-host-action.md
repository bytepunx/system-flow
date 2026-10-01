---
id: T-0666
type: task
nature: feature
title: The MCP server starts or restarts a story's agent, gated by the agent host action
status: done
parent: S-0177
owner: alex
created: 2026-10-01T10:07:04Z
updated: 2026-10-01T10:22:41Z
transitions:
  - to: ready
    at: 2026-10-01T10:07:15Z
    by: agent-S-0177
  - to: in-progress
    at: 2026-10-01T10:19:45Z
    by: agent-S-0177
  - to: done
    at: 2026-10-01T10:22:41Z
    by: agent-S-0177
stream: S-0177
tags: []
touches: [flai/internal/mcpserver, flai/cmd, flai/internal/guard]
usage:
  source: log
  seconds: 176
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 81
      output: 28406
      cache_read: 5870634
      cache_write: 123122
      cost: 2.5214
---
# T-0666 The MCP server starts or restarts a story's agent, gated by the agent host action

## Work

- MCP tools `agent_start` and `agent_restart` take `story`. They run `serve.Start` and `serve.Restart` with the options `flai serve agent` uses, so the `agent` action gates them as it gates the CLI and dashboard. They are journalled with the calling agent as who acted.
- `flai guard` refuses them to a sub-agent if its rules do not already.
- Tests: a start, and the refusal while the action is off.
- `design/system/flai-cli.md` and `docs/users/flai.md` list them.

## Done when

The tests pass and the change is committed.

## Notes
