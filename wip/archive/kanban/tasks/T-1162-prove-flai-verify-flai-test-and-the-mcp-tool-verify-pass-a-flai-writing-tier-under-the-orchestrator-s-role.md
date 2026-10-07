---
id: T-1162
type: task
nature: feature
title: Prove flai verify, flai test, and the MCP tool verify pass a flai-writing tier under the orchestrator's role
status: done
parent: S-0311
owner: alex
created: 2026-10-07T14:29:30Z
updated: 2026-10-07T14:42:30Z
transitions:
  - to: ready
    at: 2026-10-07T14:32:44Z
    by: agent-S-0311
  - to: in-progress
    at: 2026-10-07T14:32:45Z
    by: agent-S-0311
  - to: done
    at: 2026-10-07T14:42:30Z
    by: agent-S-0311
stream: S-0311
tags: [cli]
touches: [flai/cmd/verify_test.go, flai/cmd/test_test.go, flai/internal/mcpserver/verify_test.go]
after: [T-1161]
usage:
  source: log
  seconds: 585
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 69
      output: 19286
      cache_read: 2875833
      cache_write: 121700
      cost: 1.7456
---
# T-1162 Prove flai verify, flai test, and the MCP tool verify pass a flai-writing tier under the orchestrator's role

## Work

Show the acceptance criterion end to end: the orchestrator's run of the tiers completes where TH-0260's failed.

- In `flai/cmd/verify_test.go` and `flai/cmd/test_test.go`, set `t.Setenv("FLAI_ROLE", "orchestrate")` and run `flai verify` and `flai test` on a fixture whose manifest tier checks its own role, such as `sh -c 'test "$FLAI_ROLE" = verify'`. Each must pass. Keep the fixture as the existing tests in those files build theirs.
- In `flai/internal/mcpserver/verify_test.go`, do the same through the MCP tool `verify`. The server reads `FLAI_ROLE` when it starts (`mcpserver/server.go`), so set it before the server is made.
- Where a fixture can run a real flai write cheaply, add one tier that does what TH-0260's tests did, a `flai move` to `ready`, and check that it is not refused.

This task waits for T-1161, since these tests fail until the tiers run under `verify`.

## Done when

- [ ] Under `FLAI_ROLE=orchestrate`, `flai verify`, `flai test`, and the MCP tool `verify` each pass a tier that requires `FLAI_ROLE=verify`.
- [ ] `flai test flai/cmd flai/internal/mcpserver` passes.

## Notes
