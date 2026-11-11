---
id: T-2977
type: task
nature: remediation
title: The design and the operators' guide say a project's MCP server moves off a port another process holds
status: done
parent: S-0697
owner: alex
created: 2026-11-10T10:29:09Z
updated: 2026-11-11T08:08:29Z
transitions:
  - to: ready
    at: 2026-11-11T08:08:09Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T08:08:10Z
    by: system-flow
  - to: done
    at: 2026-11-11T08:08:29Z
    by: agent-S-0697
stream: S-0697
tags: [docs, host, mcp]
touches: [design/system/flai-cli.md, docs/operators/index.md]
after: [T-2974]
usage:
  source: log
  seconds: 19
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 588
      cache_read: 275165
      cache_write: 8279
      cost: 0.133
---
# T-2977 The design and the operators' guide say a project's MCP server moves off a port another process holds

## Work

Waits for T-2974, because it describes the behaviour T-2974 builds, in the words and log line T-2974 settles on.

- In `design/system/flai-cli.md`, in the `flai mcp http` row, where it says a remembered address the host has already given another project is not given twice, add that a port another process holds at a start is not given either. The project then moves to a free port the same way, and the host logs `mcp address moved`. Name I-0271 and this story, as the row names S-0110.
- In `docs/operators/index.md`, in the Address row of the MCP server's table, where it says when a project moves to a free port, add the case of a port another process holds. Say that an agent given the old address over HTTP must be given the new one, as for S-0110's move.

## Done when

- Both documents say what the host does when another process holds a project's port, and agree with T-2974's code and log line.
- `flai test design/system/flai-cli.md docs/operators/index.md` passes the markdown lint.

## Notes
