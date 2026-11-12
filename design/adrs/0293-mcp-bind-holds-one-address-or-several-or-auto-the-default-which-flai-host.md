---
id: ADR-0293
title: "mcp.bind holds one address or several, or auto, the default, which flai host resolves at each server start to 127.0.0.1 and, on a Docker engine other than Docker Desktop, the container harness network's gateway, and every address listens on the project's one port"
status: proposed
date: 2026-11-12
supersedes: []
superseded_by: []
refines: [ADR-0269]
---

# ADR-0293 mcp.bind holds one address or several, or auto, the default, which flai host resolves at each server start to 127.0.0.1 and, on a Docker engine other than Docker Desktop, the container harness network's gateway, and every address listens on the project's one port

## Context

The container harness gives a story's agent the HTTP MCP server's URL with its host as `host.docker.internal`. Docker Desktop routes that name to the host's loopback. A native engine routes it to the gateway of the container's network, `172.17.0.1` for the default bridge. `mcp.bind` held one address, `127.0.0.1` by default, so on a native engine the agent could not reach the server ([I-0254](../issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md)). S-0619 made `flai serve` refuse such a start and name the gateway to bind ([ADR-0269](0269-flai-serve-refuses-a-container-start-whose-http-mcp-server-listens-on-loopback.md)).

On 2026-11-12 the operator's host moved from Docker Desktop to the native engine after a restart. A planner's start was refused. Each way out with one address cost something:

- Binding the gateway alone takes the server off loopback, so clients on the host lose it.
- Binding `0.0.0.0` puts the server on the LAN, with a token flai does not encrypt.
- `--network host` gives up the container's network isolation.

The question is how the server serves loopback clients and containers at once, and whether flai finds the addresses for the engine by default.

## Decision

`mcp.bind` holds one address or several. Each is an IP address or a host name without a port, or `auto`. The host config writes one address as a JSON string, as before, and several as an array. `flai serve mcp bind <address> [<address>...]` sets it; `flai serve mcp clear` unsets it.

Unset is `auto`. `auto` beside explicit addresses adds them to its own.

`flai host` resolves `auto` at each server start, and `flai host restart mcp` resolves it again:

1. `127.0.0.1` always.
2. The gateway of the container harness's network, as `docker network inspect <network>` gives it, when `docker info --format '{{.OperatingSystem}}'` names an engine other than Docker Desktop and the network is not `host`. The network is the one `flai serve agent harness container --network` names, `bridge` by default.

On Docker Desktop, with no Docker, or when docker cannot be read, `auto` is `127.0.0.1` alone, and the host log warns `mcp auto resolved to loopback alone`.

Every address listens on the project's one port, kept across restarts. The port moves only to one free on every address. An address this machine cannot listen on is dropped with a warning while another remains. `flai mcp http` and `flai mcp start` take `--addr` more than once, all on one port.

The server warns `listening beyond this machine` for each address that is neither loopback nor a Docker network's gateway. A gateway gets an info line.

ADR-0269's refusal stands, on the addresses the server listens on. On a native engine and a network other than `host`, a start passes when the server listens on the network's gateway or on an unspecified address, or its URL is not loopback. Otherwise it is refused. The refusal names `flai serve mcp bind auto` and `flai host restart mcp` first, then `flai serve mcp bind 127.0.0.1 <gateway>`, then `--network host`.

A running server keeps the addresses it started on. It does not restart when the engine changes. The next container start it cannot serve is refused, naming `flai host restart mcp`.

## Consequences

- On a native engine an agent in a container reaches the server with no setting, and clients on the host still reach it on `127.0.0.1`.
- On such a host the server also listens on the gateway. Every container on that network can reach it there; the bearer token and the `Origin` refusal guard it ([ADR-0216](0216-the-http-mcp-server-s-gates-are-the-bearer-token-and-the-origin-refusal.md)).
- Docker Desktop is unchanged: `auto` is `127.0.0.1`, as the default was.
- Under `auto`, each server start asks docker once for the engine and the gateway. `flai serve mcp show` asks docker to say what `auto` resolves to now. `flai serve settings` and the dashboard's settings page ask docker nothing: they show the addresses as written and, while the server runs, what it listens on.
- A host config with one address reads and writes as before. A host config that set `127.0.0.1` keeps it and does not get the gateway.
- After an engine change the operator runs `flai host restart mcp`; the refusal says so.
- Not measured on a native engine end to end: the tests use a fake docker that answers `info`, `network inspect`, and `network ls`.

## Alternatives considered

- `127.0.0.1` as the default, with `auto` opt-in. Not chosen: a native engine host would still be refused until the operator set it. Thread TH-1099 asks the operator to confirm `auto` as the default; the default stands on that confirmation, and this alternative is the one it offers.
- The gateway alone as the default. Rejected by ADR-0269: S-0506 measured that Docker Desktop under WSL2 cannot reach a bind to WSL's `docker0` address, and loopback clients would lose the server.
- Restart the servers when the engine changes. Left for a story of its own: it needs flai to watch the engine. The refusal names the restart meanwhile.
- Probe the server from a throwaway container before each start. ADR-0269's rejection stands: it costs a container start per agent start, and `docker info` names the engine kind, which is the cause.
