---
id: T-1207
type: task
nature: improvement
title: wait_for_events tells a story's agent to end when its story awaits only another story's reply
status: backlog
parent: S-0335
owner: alex
created: 2026-10-07T20:16:32Z
updated: 2026-10-07T20:16:32Z
transitions: []
stream: S-0335
tags: [flai]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
---
# T-1207 wait_for_events tells a story's agent to end when its story awaits only another story's reply

## Work

Extend `endWhy` in `flai/internal/mcpserver/server.go`. It waits for nothing in this story; S-0331 is done before it starts. It shares no path with T-1208 and runs beside it.

- Under the same conditions as for a question to the operator (flai serve started it for its story, no task of the story in progress), answer `end: true` when a conversation of its story awaits the other story's reply.
- `why` names each conversation and the story it waits on, beside any thread.
- Update the tool's description to say so.

## Done when

- Tests cover ending on a conversation alone, on a thread and a conversation, holding while a task is in progress, and holding for an agent run by hand.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
