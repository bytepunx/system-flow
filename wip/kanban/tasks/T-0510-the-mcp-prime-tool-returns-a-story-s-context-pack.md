---
id: T-0510
type: task
nature: feature
title: The MCP prime tool returns a story's context pack
status: done
parent: S-0138
owner: alex
created: 2026-09-29T01:08:51Z
updated: 2026-09-29T01:11:29Z
transitions:
  - to: ready
    at: 2026-09-29T01:08:55Z
    by: agent-S-0138
  - to: in-progress
    at: 2026-09-29T01:08:56Z
    by: agent-S-0138
  - to: done
    at: 2026-09-29T01:11:29Z
    by: agent-S-0138
stream: S-0138
tags: [cli]
touches: [flai/internal/context, flai/internal/mcpserver, flai/cmd/prime.go]
---
# T-0510 The MCP prime tool returns a story's context pack

## Work

Move what `flai prime --story` does to build a pack into `ctxpack.ForStory` in `flai/internal/context`, so the command and the MCP server build it the same way. Add a `prime` tool to the MCP server (single project and folder) that takes a story ID and returns the pack as `flai prime --story --json` prints it. Mention it in both servers' instructions. Update the tool list tests and add a test that the tool returns the pack and refuses an ID that is not a story.

## Done when

- `flai prime --story` prints what it printed before, from `ctxpack.ForStory`.
- The `prime` tool returns the pack; the tool list tests name it; `make test` passes.

## Notes
