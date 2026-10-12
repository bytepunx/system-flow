---
id: T-1484
type: task
nature: improvement
title: A settled entry makes a conversation await neither story until the next message
status: backlog
parent: S-0372
owner: alex
created: 2026-10-12T17:23:00Z
updated: 2026-10-12T17:23:00Z
transitions: []
stream: S-0372
tags: [messages]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go]
after: [T-1482]
---
# T-1484 A settled entry makes a conversation await neither story until the next message

## Work

In `flai/internal/messages/messages.go`, give an entry a settled mark, written in its heading after the story, as the ADR of T-1482 says. `Entry` and `Entries()` read the mark, and `add` writes it.

- `Reply` takes a `Settled` option.
- `Conversation.Awaiting` returns empty when the newest entry for a story is settled. A closing entry still has no story and is still skipped.
- `AwaitingOther`, `ToSince`, and `View` follow from `Awaiting`. Check each, and give `View`'s entries a `settled` field.
- `Notify`, `Send`, `Escalate`, and `Share` write unsettled entries, so a new message makes the conversation await the other story again.

It waits for T-1482, which decides the mark and the rule. The CLI, MCP, inbox, and serve tasks wait for it, because they read the mark.

## Done when

- Tests in `flai/internal/messages/messages_test.go` cover the following, and the first fails without the change:
  - A settled reply leaves `Awaiting` empty, and the conversation open.
  - A later reply or notice from either story makes it await the other story again.
  - A settled entry is read back from the file with its mark.
  - `View` gives `awaiting` empty and the entry's `settled` true.
- `TestReplyAppendsAndAwaitsTheOther` and the other tests of the package still pass.
- `flai test` passes on the two files.

## Notes

`flaiover/src/lib/components/Messages.svelte` already shows no "awaits" label when `awaiting` is empty, so the dashboard needs no change for a settled conversation.
