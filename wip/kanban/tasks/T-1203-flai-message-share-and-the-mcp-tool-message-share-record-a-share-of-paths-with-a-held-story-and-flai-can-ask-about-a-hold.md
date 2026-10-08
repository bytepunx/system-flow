---
id: T-1203
type: task
nature: feature
title: flai message share and the MCP tool message_share record a share of paths with a held story, and flai can ask about a hold
status: in-progress
parent: S-0334
owner: alex
created: 2026-10-07T20:16:02Z
updated: 2026-10-08T09:49:11Z
transitions:
  - to: ready
    at: 2026-10-08T09:49:11Z
    by: agent-S-0334
  - to: in-progress
    at: 2026-10-08T09:49:11Z
    by: agent-S-0334
stream: S-0334
tags: [flai]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/internal/messages/share.go, flai/internal/messages/share_test.go, flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-1204]
---
# T-1203 flai message share and the MCP tool message_share record a share of paths with a held story, and flai can ask about a hold

## Work

Build the ask and the share in the messages package and its two front ends, as ADR-0134 says. It waits for T-1204, whose `workitem.Share` type the conversation's `shares` front matter holds.

- `messages.AskHold`: flai's message from the held story (ready) to the holding story (in progress), `about` the overlapping paths, with the held story's goal and the three answers; once per pair while the held story stays in ready.
- `messages.Share`: records a share on the conversation's `shares` and adds an entry with the split; refuses anyone but the holding story's agent and the holding story's owner, and paths outside the overlap.
- A reply to a conversation whose other side is ready awaits no one (`AwaitingOther` leaves it out).
- `flai message share <id> --paths <path>… "<split>"` and the MCP tool `message_share` call `Share`; `View` shows the shares.
- `flai guard` refuses a sub-agent `message_share` and `flai message share`.

## Done when

- Tests cover asking once, sharing, each refusal, the guard's refusal of a sub-agent, and a share that ends when a story leaves.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
