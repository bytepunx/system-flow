---
id: T-2441
type: task
nature: improvement
title: Mark the notice flai task done sends as a notice in its conversation entry
status: done
parent: S-0545
owner: alex
created: 2026-11-02T20:06:03Z
updated: 2026-11-06T22:09:39Z
transitions:
  - to: ready
    at: 2026-11-06T22:00:46Z
    by: agent-S-0545
  - to: in-progress
    at: 2026-11-06T22:00:47Z
    by: agent-S-0545
  - to: done
    at: 2026-11-06T22:09:39Z
    by: agent-S-0545
stream: S-0545
tags: [flai]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/message.go, docs/users/flai-reference.md]
after: [T-2440]
usage:
  source: log
  seconds: 532
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 11267
      cache_read: 5880209
      cache_write: 87264
      cost: 2.0546
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 3
          cache_read: 30803
          cache_write: 2234
          cost: 0.0078
---
# T-2441 Mark the notice flai task done sends as a notice in its conversation entry

## Work

Give a conversation entry the mark T-2440's ADR decides for a task-done notice, so flai serve can tell it from a message that asks something. Today the notice from `notices` in `flai/internal/taskdone/taskdone.go` is written by the other story's agent through `messages.Notify`, and nothing in the entry marks it. With the planner's proposal, that is a mark in the entry's heading beside `(settled)` (`Marks`, `entryHeading`, `Entry` in `flai/internal/messages/messages.go`), written by `Notify` when its options ask for it and parsed back into `Entry`, with `Awaiting` unchanged: a notice still awaits the other story. Check that the JSON an entry gives (`inbox`, `message_get`, the dashboard's reads) carries the mark, and that a conversation written before this change still parses.

Waits for T-2440: the mark is what its ADR decides.

## Done when

- A notice `flai task done` sends is written with the mark, and a reply or a message from `flai message send` is not.
- `messages` tests cover writing, parsing, and the JSON of a marked entry, and an old conversation without the mark; `taskdone` tests cover the notice's mark.
- `flai test` passes on the changed paths.

## Notes
