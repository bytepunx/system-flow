---
id: T-2312
type: task
nature: improvement
title: The MCP prime tool takes reasons with role orchestrate and returns the reason-primed pack in parts, refusing an unknown reason naming the known ones
status: done
parent: S-0537
owner: alex
created: 2026-10-28T22:07:35Z
updated: 2026-10-29T16:12:12Z
transitions:
  - to: ready
    at: 2026-10-29T15:52:19Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T15:52:20Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T16:12:12Z
    by: agent-S-0537
stream: S-0537
tags: [flai, mcp, prime]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go]
after: [T-2308, T-2309]
usage:
  source: log
  seconds: 1192
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 38
      output: 14117
      cache_read: 2821002
      cache_write: 91152
      cost: 1.3647
---
# T-2312 The MCP prime tool takes reasons with role orchestrate and returns the reason-primed pack in parts, refusing an unknown reason naming the known ones

## Work

Criterion 1, the MCP half. Add `Reasons []string` to `PrimeIn` in `flai/internal/mcpserver/server.go`, with a jsonschema description naming the forms T-2308 parses, and route it in `pack` to the reason-primed pack in place of `ForStrategic` when the role is `orchestrate`. Parts (ADR-0104) work as for any pack: the same `reasons` with `part` returns the next part. `reasons` with another role, or with `epic` or `story`, is refused as the CLI refuses it. Update the tool's description where it is registered, `flai/internal/mcpserver/folder.go`, to say what `reasons` loads.

Waits for T-2308 and T-2309, which build the pack for every reason the tests cover.

## Done when

- `prime_test.go` calls `prime` with each reason kind once, asserting the item each loads, and covers the refusals: an unknown reason, `reasons` with role `plan`, and `reasons` with `story`.
- A pack primed for two reasons returns them in the header, and part 2 of a pack over 40,000 bytes comes back with the same `reasons`.
- `flai test` passes on the changed paths.

## Notes
