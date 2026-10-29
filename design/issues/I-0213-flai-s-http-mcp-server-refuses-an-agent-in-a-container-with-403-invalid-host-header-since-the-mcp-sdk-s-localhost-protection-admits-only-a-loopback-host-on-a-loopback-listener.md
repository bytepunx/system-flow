---
id: I-0213
title: flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener
class: blocker
status: closed
count: 1
cost: 20m
first_reported: 2026-10-28T12:22:29Z
last_reported: 2026-10-28T12:22:29Z
updated: 2026-10-29T07:57:04Z
---

# I-0213 flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener

## Description
flai's HTTP MCP server refuses an agent in a container with 403 invalid Host header, since the MCP SDK's localhost protection admits only a loopback Host on a loopback listener

## Instances

### 2026-10-28T12:22:29Z
Story: S-0505.
S-0505 T-2159: Claude Code in flai-agent-claude-code:local on Docker Desktop (WSL2 mirrored networking) reaches the project's HTTP MCP server at 127.0.0.1:4243 through host.docker.internal (192.168.65.254, the host gateway): without a token the server answers 401, so TCP gets through. With the token the go-sdk v1.8.0 StreamableHTTPHandler answers 403 `Forbidden: invalid Host header "host.docker.internal:4243"`: it refuses as DNS rebinding a request on a loopback listener whose Host is not loopback, and flai/internal/mcphttp does not set DisableLocalhostProtection or take an allowed-hosts list. --network host fails on Docker Desktop (its host network is the VM's, connection refused). The spike gets through only by sending `"Host": "127.0.0.1:4243"` in mcp.json's headers, which Claude Code 2.1.312 honours and Node's fetch does not. Remediation: let flai's HTTP MCP server accept the hosts a container dials (an allowed-hosts list, as dashboard.allowed_hosts does for the dashboard, or the token as the only gate), carried by S-0506, the HTTP MCP server for containers.

## Remediation

Story S-0515 remediates this issue, created from it at 2026-10-28T12:52:34Z.
Closed 2026-10-29T07:57:04Z: Fixed by S-0515 (ADR-0216): both of the go-sdk handlers in flai/internal/mcphttp are built with DisableLocalhostProtection, so a request that reaches the loopback listener naming host.docker.internal, or any Host, is served once it carries the bearer token and no Origin; the token and the Origin refusal stay the gates. TestAForeignHostWithTheTokenIsServed reproduces the 403 invalid Host header against main's handler and passes with the fix.
