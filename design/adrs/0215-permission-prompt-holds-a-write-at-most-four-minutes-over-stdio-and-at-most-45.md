---
id: ADR-0215
title: "permission_prompt holds a write at most four minutes over stdio and at most 45 seconds over HTTP, below Claude Code's tool-call timeout on each transport"
status: accepted
date: 2026-10-29
supersedes: []
superseded_by: []
refines: [ADR-0124, ADR-0130]
---

# ADR-0215 permission_prompt holds a write at most four minutes over stdio and at most 45 seconds over HTTP, below Claude Code's tool-call timeout on each transport

## Context

[ADR-0124](0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md) bounds `permission_prompt`'s hold at four minutes. Its decision 1 says that sits "below Claude Code's idle timeout on both transports", taking the HTTP timeout to be five minutes. [ADR-0130](0130-every-harness-meets-flai-through-neutral-contracts-an-adapter-s-capabilities-a.md) makes the four minutes Claude Code's bound, and says each harness has its own.

[I-0217](../issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md) records an agent in a container whose protected write failed with `The operation timed out.`. It reached `permission_prompt` over flai's HTTP MCP server, and the thread the prompt opened stayed open with no agent waiting on it.

S-0519 measured the timeout with a scratch MCP server. Like flai's HTTP server (`flai/internal/mcphttp`), it answers in JSON response mode, so the client receives nothing until the call ends. Each figure is the seconds from the call reaching the server to Claude Code giving up:

| Where | Claude Code | `MCP_TOOL_TIMEOUT` | Tool the model calls | Permission prompt tool |
|-------|-------------|--------------------|----------------------|------------------------|
| host | 2.1.311 | unset | 60 | 60 |
| host | 2.1.311 | 600000 | 300 | 330 |
| container | 2.1.312 | unset | 60 | 60 |
| container | 2.1.312 | 600000 | not run | 330 |

- Every call carries a progress token. In JSON response mode, though, `keepAlive`'s progress notifications reach the client only when the call ends.
- A per-server `"timeout"` in the MCP config lifted both limits on the host. It was not tried in the container.

Over stdio the timeout is 30 minutes ([I-0103](../issues/I-0103-a-story-agent-s-claude-write-waits-thirty-minutes-on-an-unanswered-permission-thread-then-fails-on-claude-code-s-mcp-idle-timeout.md)), so the four minutes hold there. Over HTTP they can never answer.

The operator confirmed the bound on TH-0726 (2026-10-29).

## Decision

**`permission_prompt` holds a write at most four minutes over stdio and at most 45 seconds over HTTP.** Each bound sits below Claude Code's tool-call timeout on its transport. This refines ADR-0124's decision 1 and gives ADR-0130's per-harness bound one figure per transport.

1. The bound is a constant in flai for each transport, not a setting. The server `flai mcp` runs over stdio holds four minutes, and the one `flai mcp start` runs over HTTP holds 45 seconds, 15 seconds below the measured 60.
2. A bound of zero refuses at once, naming the open thread, without waiting.
3. Everything else stays as ADR-0124 has it. When the bound passes with no answer, the prompt refuses the write, names the thread, and leaves it open. The same write made again takes the answer given since, or waits again within the same bound. The tool's description and the refusal state the bound in force.

## Consequences

- An agent in a container now learns of the thread within 45 seconds, instead of meeting a timeout. When nothing else is left, its session ends on its own question, and `flai serve` starts it again on the answer.
- An operator has 45 seconds, not four minutes, to allow a write within the call over HTTP. Past that, the retry takes the answer.
- The same 60-second limit cuts short every other long call over HTTP, such as `wait_for_events`, `wait_for_work`, `verify`, and a slow `task_done`. [I-0228](../issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md) records this. This decision leaves those calls alone.
- A newer Claude Code that changes the timeout needs this figure measured again.

## Alternatives considered

- **Zero over HTTP, refusing at once.** This is the safest choice, but the operator could never allow a write within the call.
- **`MCP_TOOL_TIMEOUT`, or a per-server `timeout` in the MCP file the container's harness writes.** These rest on Claude Code's settings rather than flai's bound. `MCP_TOOL_TIMEOUT` alone still stops at a 300 to 330 second idle limit. A hand-written MCP config stays at 60 seconds.
- **SSE answers that carry `keepAlive`'s progress.** These keep the agent blocked for as long as the operator is away, which ADR-0124 rejected for stdio.
