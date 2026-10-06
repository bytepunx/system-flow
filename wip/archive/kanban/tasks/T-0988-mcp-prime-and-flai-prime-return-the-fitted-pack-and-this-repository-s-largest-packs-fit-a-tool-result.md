---
id: T-0988
type: task
nature: improvement
title: MCP prime and flai prime return the fitted pack, and this repository's largest packs fit a tool result
status: done
parent: S-0261
owner: alex
created: 2026-10-05T05:52:22Z
updated: 2026-10-06T23:18:24Z
transitions:
  - to: ready
    at: 2026-10-06T23:18:23Z
    by: agent-S-0261
  - to: in-progress
    at: 2026-10-06T23:18:23Z
    by: agent-S-0261
  - to: done
    at: 2026-10-06T23:18:24Z
    by: agent-S-0261
stream: S-0261
tags: [flai]
touches: [flai/internal/mcpserver, flai/cmd/prime.go, flai/cmd/prime_test.go]
after: [T-0987]
usage:
  source: log
  seconds: 1
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 16493
      cache_read: 1370312
      cache_write: 57016
      cost: 0.9544
---
# T-0988 MCP prime and flai prime return the fitted pack, and this repository's largest packs fit a tool result

## Work

Wire T-0987's pack into its callers:

- the MCP `prime` tool, in `prime` in `flai/internal/mcpserver/server.go`, with any new parameter on `PrimeIn` and the tool's description in `folder.go`;
- `flai prime` in `flai/cmd/prime.go`, so that `--json` and the text output agree with what the MCP tool returns.

Extend `TestPrimeOnThisRepository` in `flai/internal/mcpserver/prime_test.go`, or add a test beside it, so that it fails if I-0068 comes back. It encodes the tool's result as the server sends it, for this repository's planner pack for a story and for a story whose briefs exceed the budget, and checks that each stays under the limit the ADR sets.

It waits for T-0987 because it calls what that task builds.

## Done when

- The MCP `prime` result for S-0205, S-0210, S-0211, and the planner's pack for S-0261 is under the limit the ADR sets, measured as the server encodes it.
- The test reproducing I-0068 fails on the code as it was and passes now.
- The tool's description and `flai prime --help` say what the ADR changed.
- `scripts/flai-test.sh` passes.

## Notes
