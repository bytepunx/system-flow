---
id: T-1208
type: task
nature: improvement
title: flai serve records an agent that ended on a conversation as asked, and starts it again on the reply or a new message to its story
status: done
parent: S-0335
owner: alex
created: 2026-10-07T20:16:36Z
updated: 2026-10-07T21:58:19Z
transitions:
  - to: ready
    at: 2026-10-07T21:45:15Z
    by: agent-S-0335
  - to: in-progress
    at: 2026-10-07T21:45:16Z
    by: agent-S-0335
  - to: done
    at: 2026-10-07T21:58:19Z
    by: agent-S-0335
stream: S-0335
tags: [flai]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/commit.go, flai/internal/serve/commit_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go]
usage:
  source: log
  seconds: 783
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 101
      output: 40990
      cache_read: 5618280
      cache_write: 181447
      cost: 3.06
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
