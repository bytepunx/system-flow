---
id: T-1190
type: task
nature: feature
title: The MCP tools message_send, message_reply, and message_get work a conversation from the calling agent's story
status: backlog
parent: S-0331
owner: alex
created: 2026-10-07T20:14:33Z
updated: 2026-10-07T20:14:33Z
transitions: []
stream: S-0331
tags: [flai]
touches: [flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go]
---
# T-1190 The MCP tools message_send, message_reply, and message_get work a conversation from the calling agent's story

## Work

Add the tools over S-0330's `flai/internal/messages`. It waits for no task of this story; S-0330 is done before the story starts. It shares no path with T-1191 and runs beside it.

- `message_send` (to, text, about), `message_reply` (id, text), and `message_get` (id), in `flai/internal/mcpserver/messages.go`.
- The sender's story is `FLAI_STORY`, else the story the agent name `agent-S-nnnn` names; a session with neither is refused, saying so.
- Register the tools in `server.go` and route them for a folder that is not a project in `folder.go`.
- Rewrite the server's instructions and the `overlapped` guidance in `server.go`: message the other story's agent rather than open a thread, and open a thread on the operator only when the two do not agree.

## Done when

- Tests cover each tool, the story taken from `FLAI_STORY` and from the agent name, and the refusal.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
