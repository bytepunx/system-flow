---
id: T-1208
type: task
nature: improvement
title: flai serve records an agent that ended on a conversation as asked, and starts it again on the reply or a new message to its story
status: backlog
parent: S-0335
owner: alex
created: 2026-10-07T20:16:36Z
updated: 2026-10-07T20:16:36Z
transitions: []
stream: S-0335
tags: [flai]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
---
# T-1208 flai serve records an agent that ended on a conversation as asked, and starts it again on the reply or a new message to its story

## Work

Extend `judge`, `asked`, and `answered` in `flai/internal/serve/agents.go`. It waits for nothing in this story. It shares no path with T-1207 and runs beside it.

- A run whose story has a conversation awaiting the other story's reply when it ends is `asked`, as a run that ended on a question is.
- It is answered when the other story replies, or when a new message to its story arrives, even while a question to the operator is still open; it is started again in its session, as on a thread's answer.
- `Activity()` reports such an agent `waiting` with a `why` naming the story it waits on, and a field saying it waits on a story's agent rather than the operator.
- Check the race S-0317 records: an answer that comes before the next look still starts the agent again.

## Done when

- Tests cover asked on a conversation, a restart on the reply, a restart on a new message with a question still open, and the activity's reason.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
