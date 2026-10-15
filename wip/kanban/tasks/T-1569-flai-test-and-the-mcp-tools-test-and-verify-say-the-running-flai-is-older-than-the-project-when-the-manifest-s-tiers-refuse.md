---
id: T-1569
type: task
nature: improvement
title: flai test and the MCP tools test and verify say the running flai is older than the project when the manifest's tiers refuse
status: in-progress
parent: S-0387
owner: alex
created: 2026-10-14T20:54:47Z
updated: 2026-10-15T01:43:14Z
transitions:
  - to: ready
    at: 2026-10-15T01:43:14Z
    by: agent-S-0387
  - to: in-progress
    at: 2026-10-15T01:43:14Z
    by: agent-S-0387
stream: S-0387
tags: [flai, mcp, test]
touches: [flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, flai/internal/mcpserver/test.go, flai/internal/mcpserver/test_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/cmd/test.go, flai/cmd/verify.go]
after: [T-1567]
---
# T-1569 flai test and the MCP tools test and verify say the running flai is older than the project when the manifest's tiers refuse

## Work

Reproduce I-0141's first instance and give it the message T-1567's ADR decides; it waits for T-1567 for that message. It shares no path with T-1568 and runs beside it.

- `verify.CheckoutTiers` (`flai/internal/verify/manifest.go`) refuses a tier whose command holds a placeholder this flai does not fill and tells the agent to fix the manifest with `flai manifest set tests=`. When `release.FlaiOutdated` answers for the checkout, say instead that the running flai is older than the project, with both versions and the upgrade command, and keep the manifest advice only when it does not.
- The MCP tools `test` (`flai/internal/mcpserver/test.go`) and `verify` (`verify.go`) return that message as their error.
- Tests: a manifest with an unknown placeholder in a project whose history has a `flai/v*` tag newer than the running flai gives the version message from `CheckoutTiers` and from both MCP tools; with no newer tag it gives the manifest message as today.

## Done when

- The three refusals name the running and newest versions and the upgrade when the flai is behind, and the manifest fix otherwise, each with a test.
- `flai test` passes on the paths this task changes.

## Notes
