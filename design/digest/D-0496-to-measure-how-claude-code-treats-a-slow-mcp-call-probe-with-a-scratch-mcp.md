---
id: D-0496
title: To measure how Claude Code treats a slow MCP call, probe with a scratch MCP server in JSON response mode under the worktree's git-ignored .flai-cache, not with flai's own tools
kind: decision
sources: [S-0583, T-2541]
topics: [mcp, harness]
paths: [flai/internal/mcphttp/mcphttp.go, design/adrs/0258-the-mcp-file-flai-writes-for-its-http-server-names-a-per-server-timeout-of.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0496 To measure how Claude Code treats a slow MCP call, probe with a scratch MCP server in JSON response mode under the worktree's git-ignored .flai-cache, not with flai's own tools

S-0583's narrative (wip/archive/agents/S-0583.md, Decisions) records how T-2541 measured Claude Code's per-server MCP `timeout`: "T-2541 measures with a scratch MCP server in JSON response mode, as S-0519's T-2227 did, not flai's own tools: `wait_for_events` returns on any board event and `permission_prompt` opens real threads. The scratch server and its logs live in the worktree's git-ignored `.flai-cache/s0583-probe/`, because the Write tool refuses paths outside the worktree."

So a later story that measures how Claude Code's MCP client handles a slow or held call should answer from a scratch server that sleeps a set time, in the JSON response mode `flai/internal/mcphttp` answers in, rather than time flai's own tools, whose answers depend on board events or open real threads. The story's agent writes it under its worktree's `.flai-cache/`, not the session scratchpad.

The probe sessions ran `claude -p` in `flai-agent-claude-code:local` with `--model haiku`, to keep them cheap, and with the host's container login folder (D-0494). What T-2541 measured, a call held past 60 s only with a per-server `timeout`, is ADR-0258.

Sources: S-0583, T-2541.
