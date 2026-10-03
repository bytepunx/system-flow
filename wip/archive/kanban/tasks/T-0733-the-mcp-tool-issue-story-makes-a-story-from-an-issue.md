---
id: T-0733
type: task
nature: feature
title: The MCP tool issue_story makes a story from an issue
status: done
parent: S-0198
owner: arobson
created: 2026-10-03T01:28:25Z
updated: 2026-10-03T01:58:17Z
transitions:
  - to: ready
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: in-progress
    at: 2026-10-03T01:51:15Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T01:58:17Z
    by: agent-S-0198
stream: S-0198
tags: []
touches: [flai/internal/mcpserver, design/system/flai-cli.md, docs/users/flai.md]
after: [T-0732]
usage:
  source: log
  seconds: 422
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 19368
      cache_read: 3309854
      cache_write: 67840
      cost: 1.4535
---
# T-0733 The MCP tool issue_story makes a story from an issue

## Work

- The flai MCP server gets an `issue_story` tool: given an issue ID and an optional epic, it makes the story as `flai issue story` does and returns the new story's ID and path. The tool sits beside `item_new` in `flai/internal/mcpserver`.
- `design/system/flai-cli.md` and `docs/users/flai.md` list the tool where they list the MCP tools.

Waits for the third task: both edit `design/system/flai-cli.md` and `docs/users/flai.md`.

## Done when

- An mcpserver test makes a story from an issue through the tool, and is refused on a closed issue.
- `go test -race -short ./internal/mcpserver/` passes.

## Notes
