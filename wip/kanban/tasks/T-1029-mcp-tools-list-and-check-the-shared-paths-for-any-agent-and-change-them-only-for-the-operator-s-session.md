---
id: T-1029
type: task
nature: feature
title: MCP tools list and check the shared paths for any agent, and change them only for the operator's session
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:15:55Z
updated: 2026-10-06T12:15:55Z
transitions: []
stream: S-0295
tags: [flai]
touches: [flai/internal/mcpserver/shared.go, flai/internal/mcpserver/shared_test.go, flai/internal/mcpserver/server.go, flai/internal/guard]
after: [T-1025]
---
# T-1029 MCP tools list and check the shared paths for any agent, and change them only for the operator's session

## Work

In `flai/internal/mcpserver/shared.go`, registered in `server.go`, add the MCP tools T-1023's ADR names:

- **A read tool** (proposed `shared_paths`), open to every agent. It lists the patterns and, given paths or a story ID, says for each whether it is shared and which pattern matched.
- **An edit tool** (proposed `shared_paths_edit`). It adds or removes patterns through T-1025's edit function and returns what changed.

The list decides what holds, so an agent that could add to it could free its own story. `flai guard` (`flai/internal/guard`) refuses the edit tool to a story's agent, its sub-agents, the planner, and the analyzer. It is left to the operator's own session, and to the orchestrator only if the ADR grants it. The refusal says to ask the operator on a thread.

This task waits for T-1025, whose matcher and edit function the tools call. It runs beside T-1028, which touches no file of it.

## Done when

- `shared_test.go` covers both tools: list, check of paths and of a story, add, remove, and an invalid pattern refused.
- A guard test shows the edit tool refused to a story's agent and allowed to the operator's session.
- `scripts/flai-test.sh` passes.

## Notes
