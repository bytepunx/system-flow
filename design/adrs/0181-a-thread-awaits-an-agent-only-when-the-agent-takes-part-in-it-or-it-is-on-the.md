---
id: ADR-0181
title: "A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, and any other open thread reads awaiting operator or other"
status: accepted
date: 2026-10-21
supersedes: []
superseded_by: []
refines: [ADR-0090, ADR-0179]
topics: [mcp, inbox, orchestration, flai]
---

# ADR-0181 A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, and any other open thread reads awaiting operator or other

## Context

`inbox.Summarize` in `flai/internal/inbox/inbox.go` marks a thread `awaiting: "you"` whenever its last entry is by anyone but the reader, but for an opener whose recommendation awaits the operator's confirmation ([ADR-0090](0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md)). So every agent counts every thread awaiting the operator as its own.

Measured on 2026-10-20, a new story agent's inbox was 43,279 bytes, 42,403 of them threads: 31 threads, all `awaiting: "you"`, none on the agent's story. They were planners' questions, the orchestrator's escalations, and flai serve's notices, all waiting on the operator. S-0453's own `story_start` answered 90,556 characters, its inbox 77,093 bytes with 38 threads awaiting it, one on its story.

[ADR-0179](0179-task-done-and-story-start-carry-an-inbox-that-lists-the-story-s-threads-and.md) keeps `story_start` and `task_done` under Claude Code's limit by listing only the story's threads and counting the rest in `threads_omitted`. The rule underneath stands: the standalone `inbox` grows with every thread the operator has not answered, `awaiting_you` and `threads_omitted` count threads the agent should not answer, and `wait_for_work` answers `thread` to an idle agent for any new thread awaiting the operator.

## Decision

A thread awaits an agent only when the agent takes part in it or it is on the agent's own story; any other open thread reads `awaiting: "operator"` when it awaits the operator, else `"other"`.

- **Awaiting the agent.** A thread is `awaiting: "you"` when the agent is among its `participants`, or the thread is on the agent's own story or one of its tasks, and its last entry is not the agent's. A reply adds its author to `participants`, so an agent that replied on a thread takes part in it. ADR-0090's rule stands: a thread the agent opened, while a recommendation on it awaits the operator's confirmation, does not await the agent.
- **The agent's own story.** For the MCP server it is the story the message tools write as: `FLAI_STORY`, else `agent-S-nnnn` in the agent's name, else the story branch checked out. For `story_start` and `task_done` it is the story they start or close. An agent with none takes part only in the threads it wrote to.
- **Awaiting the operator.** Any other open thread is `awaiting: "operator"` when a recommendation on it is pending or its last entry is by none of the operator's names, else `"other"`. The operator's names are the manifest's `owner` and the thread's story's `owner`, those [ADR-0097](0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md) takes an answer from, or `designer` when neither is named, as the dashboard's `inbox.designer` reads it.
- **What counts and what is listed.** `awaiting_you` counts only the threads awaiting the agent, and `inbox` lists only those without `all: true`. In the inbox `story_start` and `task_done` carry, `threads_omitted` counts only threads on other items that await the agent. `wait_for_work` answers `thread` only for a thread awaiting the agent.
- **The orchestrator.** Its `inbox`, read in the orchestrator's role, lists the threads awaiting the operator beside those awaiting it, without `all: true`, and its instructions name them by `awaiting: "operator"`.
- **The operator's view.** The dashboard and `flai thread list` are unchanged.

This refines ADR-0090, whose opener now reads the thread as awaiting the operator rather than another, and ADR-0179, whose `threads_omitted` and `awaiting_you` now count only threads awaiting the agent. ADR-0179's 40,000-byte fit stays as the safety net.

## Consequences

- A story's agent sees in `inbox` the threads it takes part in and those on its story, not the operator's backlog of questions. Its inbox no longer grows with every thread the operator has not answered.
- An idle agent's `wait_for_work` is not woken by a thread it cannot answer.
- A thread on another story that an agent never wrote to does not await it, even when it is addressed to that agent by name. It reaches the agent through the operator, or as a message between stories ([ADR-0120](0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md)).
- The orchestrator keeps finding the threads it answers by its role, so its instructions barely change.

## Alternatives considered

- **Shorter summaries everywhere.** Rejected by the operator on 2026-10-20: the lists would shrink, but every agent would still count the operator's threads as its own, and `wait_for_work` would still wake for them.
- **A general size cap on `inbox`.** Rejected with it: a cap cuts what the agent needs as readily as what it does not.
- **The orchestrator asks for `all: true`.** Rejected: `all` lists threads awaiting other agents too, and the operator chose that the orchestrator see the threads awaiting the operator by its role.
