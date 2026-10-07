---
id: T-1203
type: task
nature: feature
title: flai message share and the MCP tool message_share record a share of paths with a held story, and flai can ask about a hold
status: backlog
parent: S-0334
owner: alex
created: 2026-10-07T20:16:02Z
updated: 2026-10-07T20:16:08Z
transitions: []
stream: S-0334
tags: [flai]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-1202]
---
# T-1203 flai message share and the MCP tool message_share record a share of paths with a held story, and flai can ask about a hold

## Work

Build the share in the messages package and its two front ends. It waits for T-1202, whose ADR settles who may share and when a share ends.

- `messages.AskHold(held, holding, paths, goal)` opens the pair's conversation for flai, with a ready story on one side, once while the hold lasts.
- `messages.Share(id, paths, split, by)` records the share on the conversation; it refuses anyone but the holding story's agent and the operator, and paths outside the overlap.
- `messages.Shares(items)` returns the shares in force: both stories still where the ADR says, paths, and split.
- `flai message share <id> --paths <path>… "<split>"` and the MCP tool `message_share` call `Share`.
- `flai guard` refuses a sub-agent `message_share` and `flai message share`.

## Done when

- Tests cover asking once, sharing, each refusal, the guard's refusal of a sub-agent, and a share that ends when a story leaves.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
