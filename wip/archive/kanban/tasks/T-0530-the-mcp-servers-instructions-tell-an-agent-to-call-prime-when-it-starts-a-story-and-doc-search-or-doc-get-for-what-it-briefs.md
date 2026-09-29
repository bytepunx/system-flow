---
id: T-0530
type: task
nature: feature
title: The MCP servers' instructions tell an agent to call prime when it starts a story and doc_search or doc_get for what it briefs
status: done
parent: S-0148
owner: alex
created: 2026-09-29T05:24:53Z
updated: 2026-09-29T05:26:13Z
transitions:
  - to: ready
    at: 2026-09-29T05:25:01Z
    by: agent-S-0148
  - to: in-progress
    at: 2026-09-29T05:25:25Z
    by: agent-S-0148
  - to: done
    at: 2026-09-29T05:26:13Z
    by: agent-S-0148
stream: S-0148
tags: []
touches: [flai/internal/mcpserver]
---
# T-0530 The MCP servers' instructions tell an agent to call prime when it starts a story and doc_search or doc_get for what it briefs

## Work

Both server instructions, the single-project one in `server.go` and the folder one in `folder.go`, say to call `prime` with the story when starting work on it, that the pack briefs design and ADRs, and to fetch a briefed body with `doc_get` and its heading, found with `doc_search`, before changing what it describes. Tests pin it for both servers.

## Done when

- [x] Both servers' instructions name `prime`, `doc_search`, and `doc_get` with a heading; tests pin it.
- [x] `go test ./internal/mcpserver/...` passes.

## Notes
