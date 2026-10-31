---
id: T-2419
type: task
nature: remediation
title: wait_for_events and wait_for_work over HTTP answer within 45 seconds, below Claude Code's 60-second tool-call timeout, as an ADR refining ADR-0215 settles, and I-0228 says what is left
status: backlog
parent: S-0506
owner: alex
created: 2026-10-31T07:59:10Z
updated: 2026-10-31T07:59:10Z
transitions: []
stream: S-0506
tags: [cli, mcp]
touches: [flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/work.go, flai/internal/mcpserver/timing.go, design/adrs/README.md, docs/users/flai.md, design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md, design/issues/summary.md]
after: [T-2150, T-2160]
---
# T-2419 wait_for_events and wait_for_work over HTTP answer within 45 seconds, below Claude Code's 60-second tool-call timeout, as an ADR refining ADR-0215 settles, and I-0228 says what is left

## Work

I-0228 measured that Claude Code drops an MCP call over HTTP after 60 seconds, since flai's HTTP server answers as JSON and `keepAlive`'s progress reaches the client only when the call ends. So `wait_for_events` (60 s by default) and `wait_for_work` (5 minutes by default) fail for an agent in a container, and criterion 1 names `wait_for_events` among the tools that must answer over HTTP as on stdio. Write an ADR with `flai adr new` that refines ADR-0215: over HTTP, one wait holds at most the 45-second bound `permission_prompt` holds, and returns `timed_out` for the agent to call again. Name what stays open: `verify`, which S-0546 makes a `flai serve` job waited on by event, and a `task_done` or `story_start` slower than 60 seconds. In `flai/cmd/mcp_http.go`, set the server options' `MaxWait` to that bound, with a constant beside `LongestWait` in `flai/internal/mcpserver/timing.go`. Make the descriptions of `wait_for_events` and `wait_for_work` (`server.go`, `work.go`) state the bound in force, as `permission_prompt`'s states its own. Add a sentence on the bound to `docs/users/flai.md` § MCP over HTTP. Record on I-0228, with `flai issue`, what this fixes and what S-0546 carries, and close it if nothing is left.

## Done when

A test in `flai/cmd/mcp_http_test.go` calls `wait_for_events` and `wait_for_work` over HTTP with `timeout_seconds` 600 and gets `timed_out` within the bound. The same calls over stdio still hold as they did. The ADR is in `design/adrs/README.md`. `flai test flai/cmd flai/internal/mcpserver` passes.

## Notes

Added by planner-S-0506 on 2026-10-31. I-0228 was filed by S-0519 after this story was first planned. No story remediates it, and its remedy is the HTTP server's. A per-server timeout in the MCP file the container adapter writes is S-0508's, and does not replace this bound.
