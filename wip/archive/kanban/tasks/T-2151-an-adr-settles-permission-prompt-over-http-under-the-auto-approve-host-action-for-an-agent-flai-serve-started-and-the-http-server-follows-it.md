---
id: T-2151
type: task
nature: feature
title: An ADR settles permission_prompt over HTTP under the auto-approve host action for an agent flai serve started, and the HTTP server follows it
status: done
parent: S-0506
owner: alex
created: 2026-10-26T04:47:21Z
updated: 2026-11-02T05:30:45Z
transitions:
  - to: ready
    at: 2026-11-02T05:15:26Z
    by: agent-S-0506
  - to: in-progress
    at: 2026-11-02T05:15:27Z
    by: agent-S-0506
  - to: done
    at: 2026-11-02T05:30:45Z
    by: agent-S-0506
stream: S-0506
tags: [cli, mcp]
touches: [design/adrs/README.md, flai/cmd/mcp_http.go, flai/internal/mcpserver/permission.go, flai/internal/mcpserver/permission_test.go, flai/internal/mcphttp/mcphttp.go, flai/cmd/mcp_http_permission_test.go, design/adrs/0243-over-http-permission-prompt-honours-the-auto-approve-host-action-only-for-the.md, flai/internal/mcpserver/server.go]
usage:
  source: log
  seconds: 918
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 82
      output: 25698
      cache_read: 4772506
      cache_write: 133463
      cost: 2.2437
---
# T-2151 An ADR settles permission_prompt over HTTP under the auto-approve host action for an agent flai serve started, and the HTTP server follows it

## Work

Write an ADR with `flai adr new`, refining ADR-0086 and ADR-0097, that settles whether `permission_prompt` over HTTP honours the `auto-approve` host action for an agent `flai serve` started and named, as the stdio server does today (`cmd/mcp.go` passes `AutoApprove`; `cmd/mcp_http.go` does not, so `mcpserver/permission.go` treats it as off). Recommended: honour it for an agent whose name `flai serve` registered for the run, since `flai serve` knows the action and started the agent, and keep asking for any other agent, so that the dashboard token still cannot grant it (`docs/users/flai.md` says the shared HTTP server always asks; the ADR changes that sentence's scope). Make `cmd/mcp_http.go` read the host action as `cmd/mcp.go` does and `mcpserver/permission.go` apply it only to the agents the ADR names, with a test in `permission_test.go` for an agent flai serve started and one it did not. This task waits for nothing.

## Done when

The ADR is accepted and indexed in `design/adrs/README.md`, `permission_prompt` over HTTP behaves as it says for both kinds of agent in `permission_test.go`, and `flai test flai/cmd flai/internal/mcpserver flai/internal/mcphttp` passes.

## Notes

The ADR's file is named by `flai adr new`; the story's touches keep the `design/adrs` folder for it.
