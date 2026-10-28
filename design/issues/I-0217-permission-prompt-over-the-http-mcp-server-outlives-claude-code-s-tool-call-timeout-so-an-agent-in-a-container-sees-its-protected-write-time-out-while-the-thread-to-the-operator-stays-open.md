---
id: I-0217
title: permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-28T12:30:56Z
last_reported: 2026-10-28T12:30:56Z
updated: 2026-10-28T12:53:05Z
---

# I-0217 permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open

## Description
permission_prompt over the HTTP MCP server outlives Claude Code's tool-call timeout, so an agent in a container sees its protected write time out while the thread to the operator stays open

## Instances

### 2026-10-28T12:30:56Z
Story: S-0505.
S-0505 T-2166: agent-S-0513, in flai-agent-claude-code:local with --permission-prompt-tool mcp__flai__permission_prompt over the project's HTTP MCP server, wrote .claude/s0513-probe.md once on purpose. permission_prompt opened TH-0677 "Allow Write .claude/s0513-probe.md?" on S-0513 and held the write, as ADR-0124 says, for up to four minutes; Claude Code 2.1.312 gave up first with `Error calling tool (Write): The operation timed out.` No file was written, the agent did not learn of the thread, and the thread stayed open awaiting the operator. On the host, over stdio, the same hold returns within the four minutes. Remediation for S-0506 (the HTTP MCP server for containers): keep the hold under the HTTP client's timeout, or answer at once with the thread named as ADR-0124's refusal does, and measure Claude Code's MCP tool timeout over HTTP (MCP_TOOL_TIMEOUT).

## Remediation

Story S-0519 remediates this issue, created from it at 2026-10-28T12:53:05Z.
