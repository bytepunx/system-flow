---
id: T-2407
type: task
nature: feature
title: flai verify --job and the MCP verify tool ask flai serve for the close-out job and return its id at once, and wait_for_events lists verified among its kinds
status: backlog
parent: S-0546
owner: alex
created: 2026-10-31T07:44:18Z
updated: 2026-10-31T07:44:18Z
transitions: []
stream: S-0546
tags: [flai]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/internal/mcpserver/server.go]
after: [T-2400, T-2403]
---
# T-2407 flai verify --job and the MCP verify tool ask flai serve for the close-out job and return its id at once, and wait_for_events lists verified among its kinds

## Work

Add `--job` to `flai verify` (`flai/cmd/verify.go`). It takes the commit options the close-out takes, writes the request with the serve function T-2403 exports, and prints the job id. When serve is not running for the project, it says so and exits non-zero. Without `--job`, `flai verify` runs in process as today.

Make the MCP `verify` tool (`flai/internal/mcpserver/verify.go`) ask for the job the same way and return the job id, with a description that tells the agent to wait for `verified` with `wait_for_events`. It no longer calls `verify.Verify` in the MCP server's process. Pass what it needs, such as the serve directory, through `mcpserver.Options` in `flai/cmd/mcp.go` and `mcp_http.go`. Add `verified` to `eventKindsDescription` in `flai/internal/mcpserver/server.go`. Both paths answer from the record when T-2400's rule finds a verdict at the commit, and say so.

Per S-0552's split on MS-0261, this story changes only `verify.go` and the job's new files in `flai/cmd`. `mcp.go` and `mcp_http.go` are outside that list, so tell S-0552's agent on that conversation before changing them.

Waits for T-2400 (the answer from the record) and T-2403 (the request it writes).

## Done when

- [ ] A test in `flai/cmd/verify_test.go` runs `flai verify --job` against a test serve, gets an id without waiting, and gets the refusal with no serve
- [ ] A test in `flai/internal/mcpserver/verify_test.go` calls `verify` and gets a job id, and gets the record's verdict at a commit it holds
- [ ] `TestVerifyAnswersAPassAndStoresIt`'s description check matches the new description
- [ ] `go test -race ./cmd ./internal/mcpserver -run 'Verify|WaitForEvents'` passes
- [ ] Story criteria 1 and 2 for both entry points

## Notes
