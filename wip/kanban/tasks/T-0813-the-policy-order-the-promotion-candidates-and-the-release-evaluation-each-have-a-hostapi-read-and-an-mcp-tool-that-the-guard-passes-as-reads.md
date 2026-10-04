---
id: T-0813
type: task
nature: feature
title: The policy order, the promotion candidates, and the release evaluation each have a hostapi read and an MCP tool that the guard passes as reads
status: backlog
parent: S-0217
owner: alex
created: 2026-10-04T19:09:29Z
updated: 2026-10-04T19:09:29Z
transitions: []
stream: S-0217
tags: [flai]
touches: [flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flai/cmd/hostapi_reads_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/mcpserver/server.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-0810, T-0811, T-0812]
---
# T-0813 The policy order, the promotion candidates, and the release evaluation each have a hostapi read and an MCP tool that the guard passes as reads

## Work

Add three reads to `readMethods()` in `flai/internal/hostapi/reads.go`, as `accept.preview` and `stream.diff` are added. Each returns the same JSON as its command's `--json`:

- `order.by` takes a policy and returns the computed order. It never applies it, since applying is a write.
- `promote.candidates` takes an optional limit.
- `release.evaluate` takes nothing.

Add three MCP tools in a new `flai/internal/mcpserver/orchestrate.go`, registered in `server.go`: `order_by_policy`, `promote_candidates`, and `release_evaluate`. Each describes its arithmetic in its description.

In `flai/internal/guard/guard.go`, add the three tools to `MCPReads`, so that sub-agents, the planner, and later the orchestrator pass them. Add their command forms to `cliReads`: `order --by` without `--apply`, `promote --candidates`, and `release --evaluate`. `order --by --apply` stays a write and is refused where writes are. `.claude/settings.json`'s matcher `mcp__flai__.*` already covers the tools, so it does not change.

This task waits for T-0810, T-0811, and T-0812, whose functions the reads and tools call.

## Done when

- a hostapi test for each read pins its JSON against its command's `--json` on a fixture
- an MCP test calls each tool on a fixture
- guard tests pass the three tools and the three read commands for a sub-agent and the planner, and refuse `order --by --apply` for the planner
- `go test ./internal/hostapi/ ./internal/mcpserver/ ./internal/guard/ ./cmd/` passes

## Notes
