---
id: T-2735
type: task
nature: remediation
title: flai serve refuses a container start whose loopback MCP bind a native engine's bridge cannot reach, naming flai serve mcp bind and the bridge's gateway
status: done
parent: S-0619
owner: alex
created: 2026-11-09T17:03:31Z
updated: 2026-11-09T17:17:24Z
transitions:
  - to: ready
    at: 2026-11-09T17:08:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T17:08:11Z
    by: system-flow
  - to: done
    at: 2026-11-09T17:17:24Z
    by: agent-S-0619
stream: S-0619
tags: [flai]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go]
usage:
  source: log
  seconds: 553
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 6880
      cache_read: 4078382
      cache_write: 80047
      cost: 1.578
---
# T-2735 flai serve refuses a container start whose loopback MCP bind a native engine's bridge cannot reach, naming flai serve mcp bind and the bridge's gateway

## Work

I-0254 names no remedy. Its one instance, from S-0508, proposes one: `flai serve` refuses a container start, saying why, when a container on the bridge cannot reach the MCP bind. This task builds that refusal. It does not change the default bind.

- Add a preflight beside `containerAllowed` in `flai/internal/serve/container.go`, and call it from `startOf` in `flai/internal/serve/agents.go` after `containerAllowed` and before `containerMCP` rewrites the URL's host.
- Refuse only when all three hold:
  - the run's network is not `host`;
  - the host of `FLAI_MCP_URL` is a loopback address, `127.0.0.1` by default (`mcp.bind` unset);
  - the engine is not Docker Desktop.
- Find the engine's kind with `docker info --format '{{.OperatingSystem}}'`, which answers `Docker Desktop` there. Look it up once per `flai serve` process and docker client, not once per start.
- Write the refusal to say three things, as the other container refusals do:
  - the container cannot reach the server on `127.0.0.1` through `host.docker.internal`;
  - the command that fixes it, `flai serve mcp bind <gateway>`, where `<gateway>` is the bridge's gateway from `docker network inspect bridge --format '{{(index .IPAM.Config 0).Gateway}}'`, or `docker0`'s address when that cannot be read;
  - `0.0.0.0` with its warning, as the other choice.
- When `docker info` cannot answer, start the container as today, and log the reason at `warn` in the serve journal. Never refuse a start on a guess.
- Strategic agents also run in containers (S-0509). When their start goes through another path than `startOf`, give it the same preflight, and widen the touches with that file.
- Tests in `flai/internal/serve/container_test.go`, with the fake docker the other container tests use:
  - reproduce I-0254: a native engine (`OperatingSystem` other than `Docker Desktop`) with a loopback bind is refused, and the message names `flai serve mcp bind` and the gateway;
  - Docker Desktop with a loopback bind starts;
  - a native engine with `mcp.bind` set to the gateway, or to `0.0.0.0`, starts;
  - the `host` network starts;
  - a `docker info` that fails starts and logs.

This task waits for none.

## Done when

- A container start on a native engine with the default loopback bind is refused with a message that names `flai serve mcp bind` and the bridge's gateway, and a test in `flai/internal/serve/container_test.go` reproduces it.
- Docker Desktop, a set bind, the `host` network, and a `docker info` that fails all still start, each with a test.
- `flai test` on the changed paths passes.

## Notes
