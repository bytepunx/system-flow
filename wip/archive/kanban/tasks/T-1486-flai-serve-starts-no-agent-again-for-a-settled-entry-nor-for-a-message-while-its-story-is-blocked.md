---
id: T-1486
type: task
nature: improvement
title: flai serve starts no agent again for a settled entry, nor for a message while its story is blocked
status: done
parent: S-0372
owner: alex
created: 2026-10-12T17:23:23Z
updated: 2026-10-14T21:36:39Z
transitions:
  - to: ready
    at: 2026-10-14T21:13:52Z
    by: agent-S-0372
  - to: in-progress
    at: 2026-10-14T21:13:53Z
    by: agent-S-0372
  - to: done
    at: 2026-10-14T21:36:39Z
    by: agent-S-0372
stream: S-0372
tags: [messages, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-1484]
usage:
  source: log
  seconds: 1366
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 102
      output: 36253
      cache_read: 5020678
      cache_write: 232572
      cost: 3.0376
---
# T-1486 flai serve starts no agent again for a settled entry, nor for a message while its story is blocked

## Work

In `flai/internal/serve/agents.go`, change when an agent that ended asking (`asked`) is started again, as the ADR of T-1482 says:

- `messaged`, the check for a newer message to the story, passes over a settled entry. A settled reply to a story whose agent ended for a thread does not start it.
- `answered` counts a conversation the run waited on (`run.Conversations`) as answered once it no longer awaits the other story: it is closed, it is settled, or it awaits the run's own story (`owesReply`). Without this, a run waiting on a reply that came settled would never start.
- `resume` and the asked branch of `look` do not start the agent of a blocked story (`it.Blocked` with no `Until`) for a message or a closed conversation. An answer on the thread it asked (`asking`) still starts it. `restartEnded` already skips blocked stories; follow it.
- `unanswered` and `awaitingOther`, read at `judgeRun`, follow from `messages.Awaiting` and need no change. Check that a run that ends with only settled conversations is not recorded as asked on them.

It waits for T-1484, which adds the mark. It shares no path with T-1485, so the two run together.

## Done when

- Tests in `flai/internal/serve/agents_test.go` reproduce I-0125's instances and fail without the change:
  - An agent that ended for a thread is not started again for a settled reply to its story.
  - An agent that ended waiting on a conversation is started again once, when the other story's reply comes settled, and not again when it answers settled in turn.
  - The agent of a blocked story is not started again for a new message or a conversation that closed, and is started when the thread it asked is answered.
- `TestAnAgentThatEndedAskingIsStartedAgainOnANewMessageToItsStory`, `TestAMessageInTheSecondARunEndedIsNew`, and the other tests of the package still pass.
- `flai test` passes on the two files.

## Notes

E-0019's touches name `flai/internal/serve` as a folder. Its open stories S-0350 and S-0355 may change `agents.go` meanwhile. A sync before this task brings what was accepted.
