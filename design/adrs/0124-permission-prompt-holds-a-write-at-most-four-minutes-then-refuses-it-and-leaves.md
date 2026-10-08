---
id: ADR-0124
title: "permission_prompt holds a write at most four minutes, then refuses it and leaves its thread open, and takes the answer when the agent makes the same write again"
status: accepted
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0086, ADR-0097]
---

# ADR-0124 permission_prompt holds a write at most four minutes, then refuses it and leaves its thread open, and takes the answer when the agent makes the same write again

## Context

[ADR-0086](0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md) has flai's `permission_prompt` open a thread on the story before a write to a path Claude Code protects, and hold the call until the owner answers. [ADR-0097](0097-permission-prompt-takes-an-answer-from-the-story-s-owner-or-the-project-s-owner.md) says who may answer. Neither bounds the wait.

[I-0103](../issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md) records the cost. S-0270's agent asked to write `template/root/.claude/agents/verifier.md`, and nobody answered. `permission_prompt` sends no progress while it waits, so after 1800 s Claude Code ended the call ("sent no response or progress for 1800s") and the write failed. The thread was then settled as refused, so an answer given later would have been lost. The agent stood blocked all that time, and could not tell beforehand whether the operator was there. It then staged the files in `.flai-cache/` and asked the operator to copy them in by hand.

Claude Code's idle timeout for an MCP call is 30 minutes over stdio, as the instance shows, and 5 minutes over HTTP.

On TH-0335 (2026-10-08) the planner proposed this remedy, the orchestrator recommended it, and the operator confirmed it.

## Decision

**`permission_prompt` holds a write at most four minutes. Unanswered, it refuses the write and leaves its thread open, and it takes the owner's answer when the agent makes the same write again.** This refines how long ADR-0086's call waits and what an unanswered thread becomes. Who may answer, and which words allow, stay as ADR-0097 says.

1. The call waits for an answer at most four minutes, below Claude Code's idle timeout on both transports. The bound is a constant in flai, not a setting.
2. When the four minutes pass, or the session ends, with no answer, it refuses the write. The refusal names the thread, says it stays open, and says the answer is taken when the agent makes the same write again. The thread is not settled.
3. Before it opens a thread, a request looks for an open thread on the story that the same agent opened for the same tool, path, and input. An allow already on it lets the write through at once and settles the thread. A refusal refuses the write and settles it. With no answer yet, the call waits again, bounded the same way. A request with other input opens a thread of its own.
4. The open thread is the agent's own question on its story. So `wait_for_events` ends an agent `flai serve` started once nothing else is left, and `flai serve` starts it again when the thread is answered. The agent then makes the same write again.
5. The agent goes on, meanwhile, with work that does not need the write. It no longer stages protected files in `.flai-cache/` for the operator to copy in.

## Consequences

- No permission call outlives Claude Code's idle timeout, so a write never fails on it. An agent is blocked four minutes at most.
- An answer given after the four minutes is kept and used on the retry, not lost.
- A thread stays open until the agent makes the write again. A story cancelled meanwhile leaves it open until archiving resolves it ([ADR-0109](0109-archiving-an-item-resolves-the-threads-still-open-on-it.md)).
- An agent that changes the content between attempts opens a second thread. The first stays open with the content it showed.
- `auto-approve`, the paths covered, and the refusal of `.git` are unchanged.

## Alternatives considered

- **Send progress notifications while waiting** (`keepAlive` in `flai/internal/mcpserver/timing.go`). They help only when Claude Code sends a progress token, and they keep the agent blocked for as long as the operator is away.
- **A manifest setting for the bound.** No case for another value is known. Four minutes already sits below both timeouts, and a longer one would only block the agent longer.
- **Keep the thirty-minute wait, and stage files in `.flai-cache/` with `cp` commands.** The operator copies each file by hand, and the agent cannot tell in advance whether to wait or stage.
