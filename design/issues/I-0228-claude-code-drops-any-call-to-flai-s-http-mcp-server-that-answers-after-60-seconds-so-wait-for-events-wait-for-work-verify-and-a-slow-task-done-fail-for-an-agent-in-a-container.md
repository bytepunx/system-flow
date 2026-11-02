---
id: I-0228
title: Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-29T00:56:36Z
last_reported: 2026-10-29T00:56:36Z
updated: 2026-10-29T00:56:36Z
---

# I-0228 Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container

## Description
Claude Code drops any call to flai's HTTP MCP server that answers after 60 seconds, so wait_for_events, wait_for_work, verify, and a slow task_done fail for an agent in a container

## Instances

### 2026-10-29T00:56:36Z
Story: S-0519.
Measured by S-0519 T-2227 with a scratch MCP server answering in JSON response mode, as flai/internal/mcphttp answers: Claude Code 2.1.311 on the host and 2.1.312 in flai-agent-claude-code:local give up on an HTTP MCP tool call after 60 s with "The operation timed out." when MCP_TOOL_TIMEOUT is unset, for a tool the model calls as for the permission prompt tool. With MCP_TOOL_TIMEOUT=600000 a second idle limit cuts in at 300 s (330 s for the permission prompt tool): "sent no response or progress for 300s" (CLAUDE_CODE_MCP_TOOL_IDLE_TIMEOUT). A per-server "timeout" of 600000 in the MCP config held 400 s on the host. Every call carries a progressToken, but JSON response mode delivers keepAlive's progress notifications only when the call ends. So over HTTP wait_for_events (60 s by default, up to 30 minutes), wait_for_work (5 minutes by default), verify, and any task_done or story_start slower than 60 s fail for an agent in a container, while the server goes on running them. Logs in the S-0519 narrative. S-0519 fixes only permission_prompt (I-0217); the remedy for the rest is a per-server timeout in the MCP file pilot and flai serve write, or SSE answers carrying keepAlive's progress.

## Remediation

S-0506 fixes `wait_for_events` and `wait_for_work`: over HTTP the server holds one call at most 45 seconds, the bound `permission_prompt` holds there, and then answers `timed_out` for the agent to call again ([ADR-0244](../adrs/0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md)). `verify` has answered at once since S-0546 made the close-out a `flai serve` job, whose `verified` event the agent waits for within that bound ([ADR-0229](../adrs/0229-the-close-out-runs-as-a-flai-serve-job-the-story-s-agent-asks-for-with-the.md)). Still open: a `task_done` or `story_start` slower than 60 seconds over HTTP still meets Claude Code's limit while the server finishes it. A per-server `timeout` in the MCP file S-0508's container adapter writes may lift that limit; it does not replace the bound.
