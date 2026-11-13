---
id: T-3145
type: task
nature: improvement
title: One lookup finds the host's Docker engine and the gateway of the container harness's network, for the automatic value and the container check alike
status: done
parent: S-0732
owner: alex
created: 2026-11-12T15:54:20Z
updated: 2026-11-12T17:41:33Z
transitions:
  - to: ready
    at: 2026-11-12T17:19:40Z
    by: agent-S-0732
  - to: in-progress
    at: 2026-11-12T17:19:40Z
    by: agent-S-0732
  - to: done
    at: 2026-11-12T17:41:33Z
    by: agent-S-0732
stream: S-0732
tags: [cli]
touches: [flai/internal/serve/mcpbind.go, flai/internal/serve/mcpbind_test.go, flai/internal/serve/container.go]
usage:
  source: log
  seconds: 1313
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 19363
      cache_read: 2862312
      cache_write: 104934
      cost: 1.5809
---
# T-3145 One lookup finds the host's Docker engine and the gateway of the container harness's network, for the automatic value and the container check alike

## Work

Move the engine and gateway lookup out of `containerReachable` and `containerEngine` in `flai/internal/serve/container.go` into a new `flai/internal/serve/mcpbind.go`. The lookup asks `docker info --format '{{.OperatingSystem}}'` and, for an engine other than Docker Desktop, `docker network inspect <network>` for the gateway of the network the container harness runs on, bridge when none is set. Export one function that resolves `auto` to addresses: 127.0.0.1, plus that gateway when the engine routes `host.docker.internal` to it. With Docker Desktop, with no `docker` on `PATH`, with an engine `docker info` cannot name, or with `--network host`, it gives 127.0.0.1 alone. Export a second function that says whether an address is a gateway of a container network, for the warning. Leave `containerReachable` calling the moved code, with its behaviour unchanged.

This task waits for none: it is the first layer, with the configuration.

## Done when

- `mcpbind_test.go` shows each case with a fake docker runner: Docker Desktop, a native engine with its gateway, no docker, an unreadable engine, and the host network.
- The existing `container_test.go` tests pass unchanged.
- `flai test flai/internal/serve` passes.

## Notes
