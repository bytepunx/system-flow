---
id: ADR-0252
title: "A task-done notice is marked a notice, and flai serve does not start an agent paused on its own unanswered question for a conversation that holds only notices to its story"
status: accepted
date: 2026-11-06
supersedes: []
superseded_by: []
refines: [ADR-0154, ADR-0205]
---

# ADR-0252 A task-done notice is marked a notice, and flai serve does not start an agent paused on its own unanswered question for a conversation that holds only notices to its story

## Context

[ADR-0154](0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md) lets a reply settle a conversation and keeps flai serve from starting an agent again for a settled entry, or for any message while its story is blocked. [ADR-0205](0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md) §5 then made flai serve block a story whose agent ended asking on a thread, with the reason `awaiting the operator on TH-nnnn`, but expressly not as ADR-0154 §3's block: the asked run is still started again for a message to its story, "even while its question to the operator is unanswered".

[I-0232](../issues/I-0232-flai-serve-wakes-a-story-s-agent-that-is-paused-on-a-thread-for-every-task-done-message-from-another-story-each-needing-only-a-settled-reply.md) records what that costs. Every notice `flai task done` sends to a story whose claim the commit covers ([ADR-0128](0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md), S-0333) awaits that story's reply, so `messaged` in `flai/internal/serve/agents.go` starts its paused agent. The agent reads the commit, finds it breaks nothing, replies settled, and ends again on the same question. S-0510, paused on TH-0814, was started again for MS-0323 and MS-0315 on 2026-11-02, each wake needing only a settled reply. Nothing in a conversation tells a notice from a message that asks something: `flai task done` writes it through `messages.Notify` as an ordinary entry, as the sync's conflict message and the grown-claim message are written.

ADR-0154 rejected holding a notice to a story with no commits, since its agent may be about to build on the shape the commit changed. An agent paused on its own question builds nothing until the answer comes.

## Decision

A task-done notice is marked a notice in its heading, and flai serve does not start an agent that ended asking on its own question for a conversation whose entries to its story since the run ended are only notices, while that question is unanswered. This refines ADR-0154 and ADR-0205; the rest of each stands.

1. **The notice mark.** The entry `flai task done` writes to tell another story of its commit is headed `### <time> <author> <story> (notice)`, as ADR-0154 marks a settled reply with `(settled)`. `messages.Entry` carries `notice`, and `flai message show`, `--json`, and `message_get` show it. Only that notice is marked: the message `flai stream sync` sends on a trial-merge conflict, the merge queue's conflict message, the grown-claim message, and flai serve's ask about a held ready story ([ADR-0134](0134-flai-serve-asks-the-agent-of-a-story-in-progress-about-each-ready-story-it.md)) each ask the two agents to agree, so they stay unmarked. A notice is never settled: it awaits the told story's reply, as before.
2. **No wake for notices alone while asking.** A run judged `asked` with a thread whose question is not yet answered, as `threadAnswered` reads it, is not started again for a conversation whose entries from the other story since a second before the run ended are all notices: neither as a message that came after it ended (`messaged`) nor on a conversation it ended waiting on that a notice turned to await its own story. One entry among them that is not a notice, a plain message, reply, or share, starts it as before.
3. **The answer starts it.** When the thread is answered, flai serve starts the agent as before, closing the interval ADR-0205 §5 opened. Its first `inbox` lists the notices under `messages`, awaiting its reply, and it answers them before it goes on, settled where the change breaks nothing.
4. **Unchanged.** A run that ended waiting only on another story's reply, with no question of its own open, is started on a notice as before: it may be about to build. `wait_for_events` still reports a notice as a change of kind `message`.

## Consequences

- An agent paused on a question spends no turn on notices until the answer, then answers them all in one turn.
- The told story's agent learns of a breaking change later, at its answer, rather than at the notice. It builds nothing meanwhile, so nothing it builds is lost; the other story's agent waits for the reply that long, as it waits for any reply.
- A run paused on a question is still started for a message that asks something, a share, and flai's conflict and overlap messages.
- An older flai parses a `(notice)` heading as an entry with no story, as it parsed `(settled)` before ADR-0154: it reads the entry as a closing one, so its conversation does not await the told story and its serve starts nobody for it. That holds only while a host runs a flai older than this change.

## Alternatives considered

- **Judge a notice by its text.** Not adopted: the text is prose a later change may reword, and the heading mark is how ADR-0090 and ADR-0154 mark entries already.
- **Hold every message to a story paused on a question.** Not adopted: a message that asks something, such as a conflict at sync or a request to share, needs an answer the operator's does not give, and ADR-0205 §5 kept that wake on purpose.
- **Send no notice to a story paused on a question.** Not adopted: the notice is how its agent learns of the change at all, and it needs it once the answer starts it.
