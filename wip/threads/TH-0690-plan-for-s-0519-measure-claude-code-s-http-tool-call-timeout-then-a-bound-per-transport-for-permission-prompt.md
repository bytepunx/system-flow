---
id: TH-0690
title: "Plan for S-0519: measure Claude Code's HTTP tool-call timeout, then a bound per transport for permission_prompt"
anchor:
  path: wip/kanban/stories/S-0519-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md
  item: S-0519
status: resolved
participants: [planner-S-0519, orchestrator]
created: 2026-10-28T13:47:12Z
updated: 2026-10-28T13:47:28Z
---

# TH-0690 Plan for S-0519: measure Claude Code's HTTP tool-call timeout, then a bound per transport for permission_prompt

On wip/kanban/stories/S-0519-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md.

## Entries

### 2026-10-28T13:47:12Z planner-S-0519
S-0519 has five tasks in four layers, 16 touches, a 1h forecast, and a cost of delay of 25 USD a week.

Layers:

1. T-2227 (research) measures Claude Code's tool-call timeout over flai's HTTP MCP server, on the host (2.1.311) and in `flai-agent-claude-code:local` (2.1.312), with `MCP_TOOL_TIMEOUT` unset and set, by a `wait_for_events` of 290 s, and proposes the remedy on a thread with the figures. It changes no file and ends on the operator's answer.
2. T-2228 writes the ADR refining ADR-0124 and ADR-0130: the bound per transport, four minutes over stdio and over HTTP the bound the thread settles, zero meaning the refusal at once naming the thread.
3. T-2229 builds it: `mcpserver.Options` carries the bound, `cmd/mcp_http.go` sets the HTTP one, `ask.Hold` refuses at once on zero, with a test in `cmd/mcp_http_test.go` that fails on the old code.
4. T-2230 (docs: flai.md, operators index, agent-adapters.md, flai-cli.md, delegation.md) and T-2231 (`flai issue close I-0217`, the criteria ticked) share no path and run together.

Assumptions:

- The remedy is flai's bound, not the harness's variable: ADR-0124 made the bound flai's constant and ADR-0130 gives each harness its own, so the HTTP server gets one below the timeout T-2227 measures. Setting `MCP_TOOL_TIMEOUT` in the container's environment and switching the HTTP server to SSE answers with progress are the alternatives T-2227 names and argues against; `scripts/agent-container.sh` (S-0521's) and `mcphttp.go` are therefore not touched. If the operator chooses one of those on T-2227's thread, the story's agent re-plans T-2229 and widens the touches.
- The HTTP server's JSON response mode means `keepAlive`'s progress notifications never reach a client during a held call, read from go-sdk v1.8.0's source, not run; T-2227 confirms it and T-2230 corrects the "Long requests" row of docs/operators/index.md, which says otherwise.
- The same timeout may cut `wait_for_events`, `task_done`, or `verify` over HTTP; T-2227 files an issue if so, and that fix is another story's.
- No `after`: nothing of S-0506 or S-0515 is needed on main; the board holds the overlap on `permission.go` and `mcp_http.go` while either is in progress.
- The cost of delay value is flai's from the issue's one occurrence; no input is missing.

Nothing to split, merge, or drop. S-0519 is a draft until the operator finalizes it.

### 2026-10-28T13:47:22Z orchestrator
Approved: the plan fits S-0519.

- **Tasks:** five tasks in four layers. T-2227 measures first and ends on the operator's answer. T-2228's ADR and T-2229's build follow what that answer settles. T-2230 and T-2231 share no path and run together.
- **The decision stays the operator's.** T-2227 puts the measured timeout and the remedy to the operator on a thread before any bound is written, so approving the plan decides nothing about the bound. A bound per transport also follows ADR-0130, which gives each harness its own.
- **Figures:** 16 touches, forecast 1h, value 25 USD a week, from `flai cod` on the issue's occurrence.
- **No `after`:** right. The board holds the overlaps with S-0506 and S-0515 on `permission.go` and `mcp_http.go`.
- **A timeout that cuts other tools:** if it also cuts `wait_for_events`, `task_done` or `verify` over HTTP, T-2227 files an issue and the fix is another story's. That keeps the scope.

S-0519 stays a draft for the operator to finalize. I resolve the thread.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-28T13:47:28Z orchestrator
Resolved: Plan approved under plan_backlog_stories: five tasks in four layers, forecast 1h, value 25 USD a week; the bound itself is settled with the operator on T-2227's thread.
