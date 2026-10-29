---
id: D-0417
title: "permission_prompt over HTTP: hold a write at most 45 seconds, below Claude Code's measured 60-second limit: Recommended for the operator: answer \"as recommended\""
kind: answer
sources: [TH-0726]
topics: [mcp, harness, flai, cli, go, code, all]
paths: [design/adrs, design/adrs/README.md, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/server.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, docs/users/flai.md, docs/operators/index.md, design/system/agent-adapters.md, design/system/flai-cli.md, design/conventions/delegation.md, design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md, design/issues/summary.md, design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0417 permission_prompt over HTTP: hold a write at most 45 seconds, below Claude Code's measured 60-second limit: Recommended for the operator: answer "as recommended"

The operator, alex, answered TH-0726 at 2026-10-29T01:59:12Z.

What was asked:

```text
**Recommended:** over the HTTP MCP server, `permission_prompt` holds a write at most **45 seconds**, then refuses it naming the thread, which stays open, as ADR-0124's refusal does. Over stdio it keeps four minutes. The same write made again takes the answer, as today. An ADR refining ADR-0124 records it, and T-2229 builds it.
```

The answer:

```text
**Recommended for the operator:** answer "as recommended". Over the HTTP MCP server, `permission_prompt` holds a write at most 45 seconds, then refuses it naming the thread, which stays open. Over stdio it keeps four minutes. An ADR refining ADR-0124 records it, and T-2228 and T-2229 build it.

**Why I recommend it.**
- **It keeps ADR-0124's own rule.** § Decision, point 1, sets the bound "below Claude Code's idle timeout on both transports", as "a constant in flai, not a setting". T-2227 measured that timeout over HTTP at 60 s, on the host and in the container alike. So four minutes breaks the ADR's premise on that transport, and 45 s restores it.
- **The refusal path is the one ADR-0124 already settles** (points 2 to 4): the thread stays open, the agent's session ends on its own question, and the same write made again takes the answer. Only the figure changes, and only over HTTP.
- **45 s over 0:** an operator watching the board can still allow the write within the call. The 15 s margin covers opening the thread and the round trip.
- **Over leaning on `MCP_TOOL_TIMEOUT` or a per-server timeout:** those rest on Claude Code's settings. A hand-written MCP config would fall back to 60 s.

**Why it is not my answer.** A per-transport bound refines an accepted ADR. That is a decision not yet recorded, so it is yours. T-2228 and T-2229 wait on it. I-0228 (the same 60 s limit on `wait_for_events`, `verify`, and the other long calls) stays out of this story, as the agent proposes.

Source: design/adrs/0124-permission-prompt-holds-a-write-at-most-four-minutes-then-refuses-it-and-leaves.md § Decision
```

Source: TH-0726.
