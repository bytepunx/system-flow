---
id: T-2229
type: task
nature: remediation
title: The HTTP MCP server's permission_prompt holds a write at most the HTTP bound and answers before Claude Code's tool-call timeout, with a test that reproduces I-0217
status: done
parent: S-0519
owner: alex
created: 2026-10-28T13:45:00Z
updated: 2026-10-29T02:40:37Z
transitions:
  - to: ready
    at: 2026-10-29T02:00:30Z
    by: agent-S-0519
  - to: in-progress
    at: 2026-10-29T02:00:31Z
    by: agent-S-0519
  - to: done
    at: 2026-10-29T02:40:37Z
    by: agent-S-0519
stream: S-0519
tags: [cli, mcp]
touches: [flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcpserver/server.go, flai/cmd/mcp_http.go, flai/cmd/mcp_http_test.go, flai/internal/ask/ask.go, flai/internal/ask/ask_test.go, flai/internal/mcpserver/folder.go]
after: [T-2228]
usage:
  source: log
  seconds: 2406
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 75
      output: 24035
      cache_read: 5264148
      cache_write: 214387
      cost: 2.8021
---
# T-2229 The HTTP MCP server's permission_prompt holds a write at most the HTTP bound and answers before Claude Code's tool-call timeout, with a test that reproduces I-0217

## Work

`permissionWait` in `flai/internal/mcpserver/permission.go` is one value for both transports, and nothing in the server knows which transport a call came over. Give `mcpserver.Options` in `flai/internal/mcpserver/server.go` the hold bound, set by `flai/cmd/mcp_http.go` where it builds the Options for the HTTP server to the bound the ADR settles, while `flai/cmd/mcp.go` leaves it to the stdio default of four minutes. `askOperator` passes it to `ask.Asker.Hold`; when the bound is zero, `Hold` in `flai/internal/ask/ask.go` opens or finds the thread and refuses at once with the refusal naming it, without a wait, with a test for that in `ask_test.go`. The tool description `permissionPromptDescription` and the refusal say the bound in force. The test that reproduces I-0217 goes in `flai/cmd/mcp_http_test.go`, or `permission_test.go` against the Options: the HTTP server's `permission_prompt`, called by a client whose context deadline is the measured timeout, returns the refusal naming the open thread before that deadline; it fails on the code before this task, where the hold outlives the deadline, and the tests of the four-minute hold over stdio keep passing. It waits for the ADR task: the bound is what the ADR settles.

## Done when

Over HTTP, an unanswered `permission_prompt` returns the refusal naming the thread within the HTTP bound, over stdio within four minutes as before, the new test fails on the old code and passes now, and `flai test flai/cmd flai/internal/mcpserver flai/internal/ask` passes.

## Notes

`flai/internal/mcphttp/mcphttp.go` stays as it is: the remedy recommended changes the bound, not the response mode.
