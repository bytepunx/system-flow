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
