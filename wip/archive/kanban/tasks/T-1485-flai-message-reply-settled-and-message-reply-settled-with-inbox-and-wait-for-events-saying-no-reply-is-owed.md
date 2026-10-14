---
id: T-1485
type: task
nature: improvement
title: flai message reply --settled and message_reply settled, with inbox and wait_for_events saying no reply is owed
status: done
parent: S-0372
owner: alex
created: 2026-10-12T17:23:12Z
updated: 2026-10-14T21:30:54Z
transitions:
  - to: ready
    at: 2026-10-14T21:13:51Z
    by: agent-S-0372
  - to: in-progress
    at: 2026-10-14T21:13:52Z
    by: agent-S-0372
  - to: done
    at: 2026-10-14T21:30:54Z
    by: agent-S-0372
stream: S-0372
tags: [messages, cli, mcp]
touches: [flai/cmd/message.go, flai/cmd/message_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/mcpserver/folder.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
after: [T-1484]
usage:
  source: log
  seconds: 1022
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 80
      output: 28495
      cache_read: 3946264
      cache_write: 182802
      cost: 2.3876
---
# T-1485 flai message reply --settled and message_reply settled, with inbox and wait_for_events saying no reply is owed

## Work

Expose the settled reply that T-1484 adds.

- `flai/cmd/message.go`: `flai message reply` takes `--settled`. Its output and `--json` say that the conversation awaits no one. `list` prints `settled` in place of `awaiting S-nnnn`, and `show` marks the entry.
- `flai/internal/mcpserver/messages.go`: `message_reply` takes `settled`, a boolean. Its description and `message_get`'s say what it means: use it for a reply that asks nothing, such as "agreed" or "the change does not break this story", so that the other story's agent is not woken to answer it.
- `flai/internal/mcpserver/folder.go`: the server's instructions on messages say the same.
- `flai/internal/inbox/inbox.go`: a settled conversation's `MessageSummary.Awaiting` is empty, not `other`. The `message` event of a settled entry (`Describe`) says to read it and that no reply is owed, instead of "answer it with message_reply before you go on". `endWhy` in `mcpserver/server.go` reads `Awaiting` and needs no change. Check that a settled conversation counts neither as awaiting the story's reply nor as awaiting the other's.

It waits for T-1484, which adds the mark. It shares no path with T-1486, so the two run together.

## Done when

- `flai/cmd/message_test.go` checks `reply --settled`: its output, `list`, `show`, and `--json`.
- `flai/internal/mcpserver/messages_test.go` checks the following:
  - `message_reply` with `settled` leaves the conversation awaiting neither story in `inbox`.
  - `wait_for_events` reports the settled entry with the no-reply text.
  - An agent whose only open conversation is settled is not told it owes a reply.
- `flai/internal/inbox/inbox_test.go` checks the summary's `awaiting` and the event text of a settled entry.
- `flai test` passes on the changed files.

## Notes

`docs/users/flai-reference.md` is generated from the commands' flags. T-1487 regenerates it with `make flai-reference`.
