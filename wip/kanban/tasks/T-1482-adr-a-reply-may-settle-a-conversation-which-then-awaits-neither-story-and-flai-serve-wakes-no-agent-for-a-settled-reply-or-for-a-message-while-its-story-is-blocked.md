---
id: T-1482
type: task
nature: improvement
title: "ADR: a reply may settle a conversation, which then awaits neither story, and flai serve wakes no agent for a settled reply or for a message while its story is blocked"
status: done
parent: S-0372
owner: alex
created: 2026-10-12T17:22:31Z
updated: 2026-10-14T20:58:12Z
transitions:
  - to: ready
    at: 2026-10-14T20:57:33Z
    by: agent-S-0372
  - to: in-progress
    at: 2026-10-14T20:57:34Z
    by: agent-S-0372
  - to: done
    at: 2026-10-14T20:58:11Z
    by: agent-S-0372
stream: S-0372
tags: [messages, docs]
touches: [design/adrs/README.md, design/adrs/0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md]
usage:
  source: log
  seconds: 37
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 5
      output: 20
      cache_read: 331293
      cache_write: 4391
      cost: 0.1545
---
# T-1482 ADR: a reply may settle a conversation, which then awaits neither story, and flai serve wakes no agent for a settled reply or for a message while its story is blocked

## Work

Record the remedy for I-0125 as an ADR with `flai adr new`, refining ADR-0120 (a conversation awaits the story that did not write its last entry) and ADR-0108 (restarting an agent that ended asking). Accept it with `flai adr accept` once the plan thread on S-0372 has settled it. It waits for nothing. The code tasks wait for it, because they build what it decides.

The decision to record:

- **A settled reply.** `flai message reply --settled`, or `message_reply` with `settled: true`, adds an entry marked settled in its heading, as ADR-0090 marks a recommendation on a thread. A conversation whose newest entry for a story is settled awaits neither story. `Conversation.Awaiting` returns empty, and `list`, `show`, `inbox`, and `View` say so. The conversation stays open: the next message from either story, a notice from `flai task done` among them, makes it await the other story again.
- **No wake for a settled reply.** flai serve does not start an agent again for a settled entry, unless its run ended waiting on that conversation's reply. That run is answered once the conversation no longer awaits the other story: it is closed, settled, or awaits the run's own story. `wait_for_events` still reports the entry as a `message`, and its text says that no reply is owed.
- **No wake while blocked.** While a story is blocked, flai serve does not start its agent again for a message, or for a conversation that closed. An answer on the thread the agent asked still starts it. The agent's first `inbox` lists the messages that came meanwhile.
- **No notice for a generated path.** `flai task done` leaves out of its changed-paths notice the paths that sync regenerates (`design/issues/summary.md`, ADR-0098), so a commit that changes only those tells no one.
- **Considered and not adopted.** Holding a notice to a story with no commits on its branch. Such a story's agent may still be about to build on the shape the commit changed. Once the story is not blocked, the notice is useful to it.
- **Consequence.** The work-management convention changes. An agent answers a message that needs no reply with a settled reply, including a notice that does not break its work, instead of a reply that awaits the other story.

## Done when

- An ADR under `design/adrs/` records the decision above, with its context, I-0125's four instances, and its consequences, and is accepted.
- `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

`flai adr new` names the ADR's file. Its number depends on which story writes an ADR first, so the touch is the index alone, and `flai task done` widens it with the file.
