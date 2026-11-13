---
id: T-3150
type: task
nature: improvement
title: A container start passes when the server listens on an address the container reaches, and the refusal names flai serve mcp bind auto first
status: done
parent: S-0732
owner: alex
created: 2026-11-12T15:55:05Z
updated: 2026-11-12T18:22:52Z
transitions:
  - to: ready
    at: 2026-11-12T17:56:25Z
    by: agent-S-0732
  - to: in-progress
    at: 2026-11-12T17:56:25Z
    by: agent-S-0732
  - to: done
    at: 2026-11-12T18:22:52Z
    by: agent-S-0732
stream: S-0732
tags: [cli]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents_start.go]
after: [T-3145, T-3147]
usage:
  source: log
  seconds: 1587
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 39
      output: 15614
      cache_read: 2308239
      cache_write: 84621
      cost: 1.2749
---
# T-3150 A container start passes when the server listens on an address the container reaches, and the refusal names flai serve mcp bind auto first

## Work

Have `mcpEnv` in `flai/internal/serve/agents_start.go` read `addrs` from the server's state file, falling back to `addr` for a state written by an older flai. Pass `addrs` to `containerReachable` in `flai/internal/serve/container.go`. The start passes when one of those addresses is the gateway T-3145's lookup finds for the container's network, or when the engine reaches loopback. It is refused, as now, when none is. Rewrite the refusal so that its first way out is `flai serve mcp bind auto` (or `flai serve mcp bind 127.0.0.1 <gateway>`) then `flai host restart mcp`. Keep `--network host` as the other way, and drop the advice to bind the gateway alone. Keep the warning and the start for an engine `docker info` cannot name.

It waits for T-3145, whose lookup it calls. It waits for T-3147 too, whose state file carries `addrs`. It shares `container.go` with T-3145 and goes after it.

## Done when

- `container_test.go` shows a native engine with the server on 127.0.0.1 and the bridge gateway starting, with neither `--network host` nor a gateway-only bind. This is the criterion 3 test's container side.
- The server on 127.0.0.1 alone is refused, with `flai serve mcp bind auto` first in the message.
- A state with `addr` only reads as before.
- `flai test flai/internal/serve` passes.

## Notes
