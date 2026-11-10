---
id: T-2934
type: task
nature: remediation
title: The orchestrator's MCP tool accept runs flai accept --by orchestrator on the host, with a dry run that answers the blockers
status: backlog
parent: S-0687
owner: alex
created: 2026-11-10T08:49:37Z
updated: 2026-11-10T08:49:37Z
transitions: []
stream: S-0687
tags: [flai, mcp, guard]
touches: [flai/internal/mcpserver/accept.go, flai/internal/mcpserver/accept_test.go, flai/internal/mcpserver/folder.go, flai/internal/guard/guard_rules.go, flai/internal/guard/guard.go, flai/internal/guard/orchestrate_test.go]
after: [T-2932]
---
# T-2934 The orchestrator's MCP tool accept runs flai accept --by orchestrator on the host, with a dry run that answers the blockers

## Work

Waits for T-2932, whose ADR decides the route this builds. Add an MCP tool `accept` in `flai/internal/mcpserver/accept.go`, registered in `folder.go` beside `verify` and `queue_add`: it takes the story, the commit verified, the evidence, and `dry_run`, and runs the host's `flai accept <S-nnnn> --by orchestrator --verified <commit>` with the evidence on standard input, or `--dry-run`, as the host API's `accept.run` builds its arguments (`flai/internal/hostapi/writes_items.go`). The MCP server runs on the host, so the acceptance sees the host's actions: it enqueues in place of landing while the merge queue batches, as `acceptance` in `flai/cmd/accept.go` decides. Answer the blockers, the refusal's text, or the acceptance's result as data. Let only the orchestrator call it, and only while `orchestration.permissions.accept_reviews` is on: add it to the orchestrator's MCP tools and its permission check in `flai/internal/guard/guard_rules.go` and `guard.go`, and refuse it to every other role and sub-agent. Test the tool's arguments, its refusal outside the orchestrator and with `accept_reviews` off, and a dry run under a host configuration whose `merge-queue` action is on.

## Done when

- [ ] The MCP tool `accept` runs the acceptance on the host with the host's configuration, and its dry run answers the blockers
- [ ] flai guard lets the orchestrator call it only with `accept_reviews` on and refuses it to every other agent
- [ ] Tests in `accept_test.go` and `orchestrate_test.go` cover the arguments, the refusals, and the queue's answer, and pass with `flai test`

## Notes
