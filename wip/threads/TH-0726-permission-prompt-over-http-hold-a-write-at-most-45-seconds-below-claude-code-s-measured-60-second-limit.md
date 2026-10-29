---
id: TH-0726
title: "permission_prompt over HTTP: hold a write at most 45 seconds, below Claude Code's measured 60-second limit?"
anchor:
  path: wip/kanban/stories/S-0519-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md
  item: S-0519
status: open
participants: [agent-S-0519, orchestrator]
created: 2026-10-29T00:56:48Z
updated: 2026-10-29T00:57:08Z
---

# TH-0726 permission_prompt over HTTP: hold a write at most 45 seconds, below Claude Code's measured 60-second limit?

On wip/kanban/stories/S-0519-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md.

## Entries

### 2026-10-29T00:56:48Z agent-S-0519
**Recommended:** over the HTTP MCP server, `permission_prompt` holds a write at most **45 seconds**, then refuses it naming the thread, which stays open, as ADR-0124's refusal does. Over stdio it keeps four minutes. The same write made again takes the answer, as today. An ADR refining ADR-0124 records it, and T-2229 builds it.

## What T-2227 measured

A scratch MCP server answered in JSON response mode, as flai's HTTP server does. Each figure is the seconds from the call reaching the server to Claude Code giving up.

| Where | Claude Code | MCP_TOOL_TIMEOUT | Tool the model calls | Permission prompt tool |
|-------|-------------|------------------|----------------------|------------------------|
| host | 2.1.311 | unset | 60 | 60 |
| host | 2.1.311 | 600000 | 300 | 330 |
| container | 2.1.312 | unset | 60 | 60 |
| container | 2.1.312 | 600000 | not run | 330 |

- Unset, every failure reads `The operation timed out.`, which is I-0217's text exactly.
- With `MCP_TOOL_TIMEOUT` set, a second idle limit cuts in: "sent no response or progress for 300s".
- A per-server `"timeout": 600000` in the MCP config held 400 s on the host. It was not tried in the container.
- Every call carries a progress token. JSON response mode delivers no progress while the call runs.

So today's four-minute hold can never answer over HTTP. The agent always sees a timeout and never learns of the thread.

## Why 45 seconds

- It sits 15 s below the measured limit. That covers opening the thread and the round trip.
- An operator watching the board can still allow the write within the call.
- It is a constant in flai, like the four minutes. It needs no harness setting.

## Alternatives

- **Zero:** refuse at once, naming the thread. This is the safest choice, but the operator can never answer within the call.
- **`MCP_TOOL_TIMEOUT` or a per-server timeout in the container's MCP file:** this rests on Claude Code's settings rather than flai's bound. A hand-written MCP config stays at 60 s.
- **SSE answers that carry progress:** these keep the agent blocked for as long as the operator is away, which is what ADR-0124 rejected.

## A related issue I recorded

The same 60 s limit cuts short every long call over HTTP: `wait_for_events`, `wait_for_work`, `verify`, and a slow `task_done`. I recorded this as I-0228 and left it out of this story.

Answer "as recommended", or name another figure, such as 0 or 50. T-2228 and T-2229 wait on your answer.

Source: design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md § Decision

### 2026-10-29T00:57:08Z orchestrator (recommendation)
**Recommended for the operator:** answer "as recommended". Over the HTTP MCP server, `permission_prompt` holds a write at most 45 seconds, then refuses it naming the thread, which stays open. Over stdio it keeps four minutes. An ADR refining ADR-0124 records it, and T-2228 and T-2229 build it.

**Why I recommend it.**
- **It keeps ADR-0124's own rule.** § Decision, point 1, sets the bound "below Claude Code's idle timeout on both transports", as "a constant in flai, not a setting". T-2227 measured that timeout over HTTP at 60 s, on the host and in the container alike. So four minutes breaks the ADR's premise on that transport, and 45 s restores it.
- **The refusal path is the one ADR-0124 already settles** (points 2 to 4): the thread stays open, the agent's session ends on its own question, and the same write made again takes the answer. Only the figure changes, and only over HTTP.
- **45 s over 0:** an operator watching the board can still allow the write within the call. The 15 s margin covers opening the thread and the round trip.
- **Over leaning on `MCP_TOOL_TIMEOUT` or a per-server timeout:** those rest on Claude Code's settings. A hand-written MCP config would fall back to 60 s.

**Why it is not my answer.** A per-transport bound refines an accepted ADR. That is a decision not yet recorded, so it is yours. T-2228 and T-2229 wait on it. I-0228 (the same 60 s limit on `wait_for_events`, `verify`, and the other long calls) stays out of this story, as the agent proposes.

Source: design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md § Decision
