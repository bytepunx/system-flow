---
id: T-2671
type: task
nature: improvement
title: messages.List and workitem's Conversations read conversations through the kept store
status: backlog
parent: S-0602
owner: alex
created: 2026-11-07T23:50:04Z
updated: 2026-11-07T23:50:04Z
transitions: []
stream: S-0602
tags: [cli]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go]
after: [T-2660]
---
# T-2671 messages.List and workitem's Conversations read conversations through the kept store

## Work

Two readers parse every conversation file again on each call:

- `List`, `Get`, and `For` in `flai/internal/messages/messages.go`, which `messages.list`, `agent.state`'s `AwaitingOther`, and `stats.get` call. Move them onto a package-level `workitem.Kept[*Conversation]`, with `Read` as its parse and a deep copy as its clone.
- `Repo.Conversations` in `flai/internal/workitem/share.go`, which `Repo.Holds` calls for `agent.state` and `board.get`. Move it onto a second `Kept` of the conversation's front matter, its shares, and its asks, and keep its rule that a file that does not parse is skipped.

`workitem` cannot import `messages`, so the two keep apart. Both stay correct after a write, because `messages`' writes go through `atomicfile`.

It waits for T-2660, whose store it uses.

## Done when

- `messages_test.go` and `share_test.go` each show that a reply written, a conversation added, and one removed are seen on the next read, and that a second read with nothing changed parses no file.
- Each shows that two reads at once return equal answers.
- `flai test` on the four files passes.

## Notes
