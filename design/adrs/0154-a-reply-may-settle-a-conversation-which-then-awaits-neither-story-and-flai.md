---
id: ADR-0154
title: "A reply may settle a conversation, which then awaits neither story, and flai serve starts no agent again for a settled reply or for a message while its story is blocked"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0108, ADR-0120]
---

# ADR-0154 A reply may settle a conversation, which then awaits neither story, and flai serve starts no agent again for a settled reply or for a message while its story is blocked

## Context

[ADR-0120](0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md) keeps the messages between two open stories' agents in a conversation, which awaits the story that did not write its last entry. [ADR-0108](0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md) and the launcher before it start an agent that ended asking again when what it waited on is answered, and on a new message to its story.

[I-0125](../issues/I-0125-a-settled-conversation-between-two-stories-agents-keeps-waking-each-in-turn-because-every-reply-awaits-the-other.md) recorded eight instances of agents woken for nothing, in four kinds:

- **Agreement wakes the other side.** A reply that only says "agreed, nothing to change" makes the conversation await the other story, which wakes its agent, whose agreement wakes the first in turn (MS-0021, MS-0022, MS-0037, MS-0041, MS-0109, MS-0110). The convention made it worse: "Reply on the conversation either way, since it awaits your story until you do".
- **A notice for a generated path.** `flai task done` told stories about `design/issues/summary.md`, which every sync regenerates ([ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)), so the notice could break nothing (MS-0021, MS-0031, MS-0053).
- **Wakes while blocked.** S-0337 was blocked on a thread for four days, and flai serve started its agent 12 times for notices, acknowledgements, and closing entries it could not act on.
- **Notices inside a share's split.** MS-0101 woke S-0356 three times for changes inside S-0239's side of the split; each needed only an acknowledgement, which woke S-0239 in turn.

## Decision

A reply may settle a conversation, which then awaits neither story, and flai serve starts no agent again for a settled reply or for a message while its story is blocked.

1. **A settled reply.** `flai message reply --settled`, or the MCP tool `message_reply` with `settled: true`, adds an entry marked settled in its heading, `### <time> <author> <story> (settled)`, as [ADR-0090](0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md) marks a recommendation on a thread. While the newest entry from a story is settled, the conversation awaits neither story: `awaiting` is empty in `flai message list`, `show`, `--json`, `inbox`, and the dashboard. It stays open. The next message from either story, a notice from `flai task done` among them, makes it await the other story again. Only a reply can be settled; `send`, `share`, `escalate`, and the notices write unsettled entries.
2. **No wake for a settled entry.** flai serve does not start an agent again for a settled entry to its story. A run that ended waiting on a conversation is answered once that conversation no longer awaits the other story: it is closed, settled, or awaits the run's own story. So a run waiting on the very reply that came settled is started once, and a settled answer to it starts no one. `wait_for_events` still reports a settled entry as a change of kind `message`, and its text says that no reply is owed.
3. **No wake while blocked.** While a story is blocked, flai serve does not start its agent again for a message or for a conversation that closed. An answer on a thread it asked still starts it. Its first `inbox` after that lists the messages that came meanwhile.
4. **No notice for a generated path.** `flai task done` leaves the paths sync regenerates out of its changed-paths notice, so a commit that changes only those tells no one.

## Consequences

- The baseline rule of `work-management.md` on answering a message changes, in this repository and the template: a reply that asks nothing is sent settled, and a task-done notice that does not break your work is answered settled.
- Notices inside the other side of a share's split, I-0125's fourth kind, still come, but each is answered with a settled reply, which wakes no one.
- An agent can still misjudge and settle a conversation that needed an answer. The other agent reads the entry at its next look, and sends a new message if it disagrees, which awaits the first story again.
- `flai stats`' conversation figures ([ADR-0138](0138-flai-stats-reports-per-week-the-conversations-between-agents-and-how-they-ended.md)) do not count settled entries apart; that would be a decision of its own.

## Alternatives considered

- **Hold a notice to a story with no commits on its branch.** Not adopted: such a story's agent may be about to build on the shape the commit changed.
- **Leave `awaiting` empty after any message that asks nothing, judged by flai.** Not adopted: flai cannot tell from the text whether a reply asks something; the author can.
- **Stop notices for paths inside the other story's side of a share's split.** Not adopted here: the settled reply removes the ping-pong, and a split is prose flai does not parse.
