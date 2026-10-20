---
id: S-0453
type: story
nature: remediation
title: A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, so a story's agent no longer counts the threads awaiting the operator as its own
status: backlog
owner: alex
created: 2026-10-20T21:49:07Z
updated: 2026-10-20T22:12:55Z
transitions: []
tags: [flai]
topics: [mcp, inbox]
after: [S-0450]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-10-20T22:12:55Z
---
# S-0453 A thread awaits an agent only when the agent takes part in it or it is on the agent's own story, so a story's agent no longer counts the threads awaiting the operator as its own

## Goal

`Summarize` in `flai/internal/inbox/inbox.go` marks a thread `awaiting: "you"` whenever its last entry is by anyone but the reader. So every story's agent counts every thread awaiting the operator as its own. Measured on 2026-10-20, a new story agent's inbox was 43,279 bytes, 42,403 of them threads: 31 threads, all `awaiting: "you"`, none on the agent's story. They were planners' questions, the orchestrator's escalations, and flai serve's notices, all waiting on the operator.

S-0450 (ADR-0179) keeps `story_start` and `task_done` under Claude Code's limit by listing only the story's threads and counting the rest in `threads_omitted`. The rule underneath still stands. The standalone `inbox` is near the 50,000-character limit and grows with every thread the operator has not answered. `awaiting_you` and `threads_omitted` count threads the agent should not answer. And `wait_for_work` (`flai/internal/mcpserver/work.go`) answers `thread` to an agent with no story in progress for any new thread awaiting the operator.

This story fixes the rule. The operator chose this on 2026-10-20 over shorter summaries everywhere and a general size cap. ADR-0179's 40,000-byte fit stays as the safety net.

## Acceptance criteria
- [ ] An accepted ADR records the rule. A thread awaits an agent only when the agent is among its participants, or the thread is on the agent's own story or one of its tasks, and its last entry is not the agent's. ADR-0090's rule for an opener whose recommendation awaits confirmation stands. Any other open thread is `awaiting: "operator"` when it awaits the operator, else `"other"`. The ADR refines ADR-0179 and ADR-0090.
- [ ] `inbox`, `wait_for_work`, `story_start`, and `task_done` apply the rule. A test reproduces the 2026-10-20 measurement: open threads awaiting the operator, none on the reading agent's story. The story's agent then gets no thread listed, `awaiting_you` 0, `threads_omitted` 0, and no `thread` answer from `wait_for_work`, while a thread the operator answered on the agent's story still awaits it.
- [ ] The orchestrator still finds the threads awaiting the operator in `inbox` without `all: true`, as its role's default, and its instructions in `flai/internal/harness/harness.go` name them by `awaiting: "operator"`.
- [ ] The operator's view is unchanged: the dashboard and `flai thread list` show the same threads awaiting the operator as before.

## Tasks

## Notes

Decided with the operator on 2026-10-20: the orchestrator sees the threads awaiting the operator by default, as its role, rather than by asking for `all: true`, so that its instructions barely change.

Open for the planner:
- How flai knows the operator's name, to tell "awaiting the operator" from "awaiting another agent". The threads' operator entries are by `alex`, the story's `owner`.
- Whether a reply adds its author to a thread's `participants`.

Seen alongside, and not this story's: 17 of the 31 open threads were flai serve's "host flai is older than …" notices, each replaced by a newer one and never resolved, 18 KB of the 42 KB. Resolving a replaced notice when a newer one opens would be a separate story.
