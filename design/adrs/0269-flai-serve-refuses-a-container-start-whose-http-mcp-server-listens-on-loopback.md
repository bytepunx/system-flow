---
id: ADR-0269
title: "flai serve refuses a container start whose HTTP MCP server listens on loopback while the container joins a network other than the host's on a Docker engine other than Docker Desktop, naming flai serve mcp bind with the network's gateway"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
---

# ADR-0269 flai serve refuses a container start whose HTTP MCP server listens on loopback while the container joins a network other than the host's on a Docker engine other than Docker Desktop, naming flai serve mcp bind with the network's gateway

## Context

The container harness (S-0508) gives a story's agent the HTTP MCP server's URL with its host rewritten to `host.docker.internal`, which `--add-host host.docker.internal:host-gateway` maps to the engine's host gateway. The server listens on `mcp.bind`, `127.0.0.1` by default. Docker Desktop forwards the host gateway to the host's loopback, so the agent reaches it (S-0505, S-0506). A native Linux engine's host gateway is the network's bridge address, such as `docker0`'s `172.17.0.1`, which a server on `127.0.0.1` does not answer, so the agent's first MCP call fails after the container has started ([I-0254](../issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md)). The documents told the operator to set `mcp.bind`; nothing stopped a start that could not work.

## Decision

`flai serve` refuses a container start that a native engine's bridge cannot reach flai's MCP server from, and names the bind that fixes it.

On the container harness, `startOf` refuses the start, once the MCP server's URL is known and before the URL is rewritten, when all three hold:

1. the URL's host is loopback: `localhost`, an address in `127.0.0.0/8`, or `::1`;
2. the container's network is not `host` (none named is `bridge`, as the adapter runs it);
3. `docker info --format '{{.OperatingSystem}}'` answers something other than `Docker Desktop`.

The refusal is the start's failure. It names the bind, the engine, the network, and its gateway, as `docker network inspect <network>` gives it, and says to run `flai serve mcp bind <gateway>` and then `flai host restart mcp`, or to run the container on the host's network with `flai serve agent harness container --network host`. When the gateway cannot be read, it names that command in place of the address.

The engine's answer is kept by the launcher once docker gives one. A `docker info` that fails or answers nothing is warned of (`container engine unreadable`) and the start goes on: the check never refuses on a guess. The default bind stays `127.0.0.1`.

## Consequences

- On a native engine the agent is not started into a container that cannot reach flai; the operator sees why on the run and in the journal, with the command that fixes it.
- Starts on Docker Desktop, on the host's network, or with a bind that is not loopback are not changed. Each start on the container harness with a loopback URL runs `docker info` once per launcher and, on a native engine, `docker network inspect` once per refusal.
- `scripts/agent-container.sh`, which the operator runs by hand, keeps its own rewrite and gets no check.
- Not measured on a native engine: the host this was built on runs Docker Desktop. The tests use a fake docker that answers `info` and `network inspect`.

## Alternatives considered

- Bind the server to the bridge's gateway by default. Rejected: S-0506 measured that Docker Desktop under WSL2 cannot reach a bind to WSL's `docker0` address, so the default would break the one engine measured here.
- Probe the server from a throwaway container on the network before each start. Rejected for now: it costs a container start per agent start, and `docker info` names the engine kind, which is the cause.
- Leave it to the documents. Rejected: the failure showed only as the agent's first failed MCP call, inside the container.
