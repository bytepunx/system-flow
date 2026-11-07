---
id: I-0228
title: Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container
class: defect
status: closed
count: 4
cost: 4m
first_reported: 2026-10-29T00:56:36Z
last_reported: 2026-11-07T06:52:33Z
updated: 2026-11-07T07:35:06Z
---

# I-0228 Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container

## Description
Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container

## Instances

### 2026-10-29T00:56:36Z
Story: S-0519.
Measured by S-0519 T-2227 with a scratch MCP server answering in JSON response mode, as flai/internal/mcphttp answers: Claude Code 2.1.311 on the host and 2.1.312 in flai-agent-claude-code:local give up on an HTTP MCP tool call after 60 s with "The operation timed out." when MCP_TOOL_TIMEOUT is unset, for a tool the model calls as for the permission prompt tool. With MCP_TOOL_TIMEOUT=600000 a second idle limit cuts in at 300 s (330 s for the permission prompt tool): "sent no response or progress for 300s" (CLAUDE_CODE_MCP_TOOL_IDLE_TIMEOUT). A per-server "timeout" of 600000 in the MCP config held 400 s on the host. Every call carries a progressToken, but JSON response mode delivers keepAlive's progress notifications only when the call ends. So over HTTP wait_for_events (60 s by default, up to 30 minutes), wait_for_work (5 minutes by default), verify, and any task_done or story_start slower than 60 s fail for an agent in a container, while the server goes on running them. Logs in the S-0519 narrative. S-0519 fixes only permission_prompt (I-0217); the remedy for the rest is a per-server timeout in the MCP file pilot and flai serve write, or SSE answers carrying keepAlive's progress.

### 2026-11-02T06:19:02Z
Story: S-0506.
S-0506 T-2419 remediated the waits: over HTTP, wait_for_events and wait_for_work now hold a call at most 45 seconds and answer timed_out (ADR-0244), and verify has answered at once since S-0546 made the close-out a flai serve job (ADR-0229). Still open: a task_done or story_start slower than 60 seconds over HTTP, and the per-server timeout in the MCP file the container adapter writes (S-0508).

### 2026-11-07T04:45:40Z
Story: S-0565.
Met again under flai serve on harness container, S-0558's run of 2026-11-07 (session 47f3eec4), worked to review in its container. No call failed, but the slow ones came near the 60 s limit. `story_start` took 41.8 s over HTTP, called at 04:36:34Z and answered at 04:37:16Z. A host run's took 4.9 s over stdio (S-0570, 04:14:43Z). Its three `task_done` calls took 26.2 s, 21.6 s, and 14.2 s. The MCP file the container adapter wrote, `.flai-cache/agents/S-0558/mcp.json`, has no per-server `timeout`, so Claude Code's 60 s default holds. A slower host, or a story with more open branches to trial-merge, would push `story_start` past it. Recorded in S-0565's results document.

### 2026-11-07T06:52:33Z
Story: S-0556.
S-0556's agent, in a container reaching flai's MCP server over HTTP, called the MCP tool test on flai/internal/guard and the orchestrator role files as its start prompt says. The call answered "The operation timed out." The same paths through flai test in Bash passed, with go-test taking 3m59.7s. The start prompt's instruction to run tests with the MCP tool cannot hold for a Go package whose importers take minutes here; the agent used the CLI for the rest of the story, and told its sub-agents to.

## Remediation

S-0506 fixes `wait_for_events` and `wait_for_work`: over HTTP the server holds one call at most 45 seconds, the bound `permission_prompt` holds there, and then answers `timed_out` for the agent to call again ([ADR-0244](../adrs/0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md)). `verify` has answered at once since S-0546 made the close-out a `flai serve` job, whose `verified` event the agent waits for within that bound ([ADR-0229](../adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)). Still open: a `task_done` or `story_start` slower than 60 seconds over HTTP still meets Claude Code's limit while the server finishes it. A per-server `timeout` in the MCP file S-0508's container adapter writes may lift that limit; it does not replace the bound.

Story S-0583 remediates this issue, created from it at 2026-11-07T04:46:50Z.

Closed 2026-11-07T07:35:06Z: S-0583 (ADR-0258): the MCP file flai writes for its HTTP server, pilot's mcp.json in the agent container and the sample flai mcp status prints, names a per-server timeout of 1800000 ms, so task_done, story_start, test, and prime run past Claude Code's 60 seconds; a test call of 126.6 s over HTTP in flai-agent-claude-code:local returned its result, and the timeout held a 660 s call. wait_for_events, wait_for_work, and permission_prompt hold 45 seconds over HTTP (ADR-0244, ADR-0215), and verify answers at once (ADR-0229). A running agent image carries the timeout only once make agent-images rebuilds it; a hand-written MCP file without it stays at 60 s.
