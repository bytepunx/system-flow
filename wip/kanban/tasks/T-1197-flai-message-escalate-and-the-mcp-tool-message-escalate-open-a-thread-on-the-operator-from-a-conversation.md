---
id: T-1197
type: task
nature: feature
title: flai message escalate and the MCP tool message_escalate open a thread on the operator from a conversation
status: backlog
parent: S-0332
owner: alex
created: 2026-10-07T20:15:12Z
updated: 2026-10-07T20:15:20Z
transitions: []
stream: S-0332
tags: [flai]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-1194]
---
# T-1197 flai message escalate and the MCP tool message_escalate open a thread on the operator from a conversation

## Work

Give the agents a way out when they do not agree. It waits for T-1194, whose ADR settles what an escalation records. It shares no path with T-1195 or T-1196 and runs beside them.

- `messages.Escalate` opens a thread on the story of the agent escalating, titled for both stories, whose first entry names both, links the conversation, and quotes the reason; the conversation gets an entry naming the thread.
- `flai message escalate <id> "<reason>"` and the MCP tool `message_escalate` call it; only one of the two stories may escalate.
- `flai guard` refuses a sub-agent `message_escalate` and `flai message escalate`, as S-0331 made it refuse the other message writes.

## Done when

- Tests cover escalation from each side, the thread's text, the conversation's entry, a refusal from a third story, and the guard's refusal of a sub-agent.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
