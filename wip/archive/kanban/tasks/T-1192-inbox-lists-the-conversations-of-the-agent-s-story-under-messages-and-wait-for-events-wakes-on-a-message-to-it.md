---
id: T-1192
type: task
nature: feature
title: inbox lists the conversations of the agent's story under messages, and wait_for_events wakes on a message to it
status: done
parent: S-0331
owner: alex
created: 2026-10-07T20:14:45Z
updated: 2026-10-07T21:24:44Z
transitions:
  - to: ready
    at: 2026-10-07T21:12:48Z
    by: agent-S-0331
  - to: in-progress
    at: 2026-10-07T21:12:48Z
    by: agent-S-0331
  - to: done
    at: 2026-10-07T21:24:44Z
    by: agent-S-0331
stream: S-0331
tags: [flai]
touches: [flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/inbox_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/items_write_test.go, flai/internal/mcpserver/messages.go, flai/internal/mcpserver/messages_test.go, flai/internal/storystart/start.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go]
after: [T-1190]
usage:
  source: log
  seconds: 716
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 84
      output: 30309
      cache_read: 5167052
      cache_write: 184675
      cost: 2.7568
---
# T-1192 inbox lists the conversations of the agent's story under messages, and wait_for_events wakes on a message to it

## Work

Show messages where an agent already looks. It waits for T-1190, because both change `server.go` and this one reports the conversations those tools write.

- `inbox` gains `messages`: the open conversations of the agent's story, each with its ID, the other story, its `about` paths, its last entry, and `awaiting` (`you` or `other`). `awaiting_you` still counts threads only.
- `wait_for_events` watches the messages folder in `watched()`, and reports a new entry to the agent's story as an event of kind `message`, naming the conversation and the sender's story, once per cursor.
- An agent with no story gets no `messages`.

## Done when

- Tests cover the inbox field, a wait woken within a poll of a message arriving, an entry by the agent itself not reported back, and `awaiting_you` unchanged.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
