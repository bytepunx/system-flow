---
id: ADR-0216
title: "The HTTP MCP server's gates are the bearer token and the Origin refusal, whatever Host a request names, and the MCP SDK's localhost Host check is off"
status: accepted
date: 2026-10-29
supersedes: []
superseded_by: []
refines: [ADR-0030]
topics: [mcp, docker]
---

# ADR-0216 The HTTP MCP server's gates are the bearer token and the Origin refusal, whatever Host a request names, and the MCP SDK's localhost Host check is off

## Context

[ADR-0030](0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md) serves flai's MCP server over HTTP for agents that cannot start a process on the host, behind two gates of flai's own. It refuses any request that carries an `Origin` header before it looks at the token, and it admits nothing without the bearer token in `.flai-cache/mcp.token`, mode 0600.

The MCP go-sdk adds a third check of its own to every `StreamableHTTPHandler`, its localhost protection against DNS rebinding. A request that reaches a loopback listener with a `Host` header that is not loopback is refused with `403 Forbidden: invalid Host header`.

[I-0213](../issues/I-0213-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md) records an agent in a container refused so. The agent dials `host.docker.internal:4243`, Docker Desktop forwards the connection to the listener on `127.0.0.1`, and the request still names `host.docker.internal`. The token was right, and the SDK refused it before flai's tools ran.

The rebinding attack the SDK guards against is a page in a browser that resolves its own name to `127.0.0.1`. Such a page sends an `Origin` and has no token, so flai's own gates refuse it before the SDK sees the request.

## Decision

The HTTP MCP server's gates are the bearer token and the `Origin` refusal, whatever `Host` a request names, and the MCP SDK's localhost `Host` check is off.

1. A request is served whatever `Host` it names once it carries no `Origin` and the right token.
2. `flai/internal/mcphttp` builds both of its SDK handlers, the one with sessions and the one without, with `DisableLocalhostProtection`, so the SDK's `Host` check never runs.
3. There is no allowed-hosts setting for the HTTP MCP server.

## Consequences

- An agent in a container reaches the server through whatever name its runtime gives the host, such as `host.docker.internal`, with no setting to keep in step.
- A DNS-rebinding page is refused by the `Origin` check, and would be by the token: it holds neither. The token remains the one secret that admits a caller, so `flai mcp token --rotate` stays the answer to a leaked token.
- `TestAForeignHostWithTheTokenIsServed` in `flai/internal/mcphttp` pins the decision: a foreign `Host` with the token is served in both revisions, and without the token, or with an `Origin`, it is refused.

## Alternatives considered

- **An allowed-hosts list**, as `dashboard.allowed_hosts` is for the dashboard ([ADR-0193](0193-the-dashboard-accepts-a-browser-origin-naming-one-of-the-hosts-in-dashboard.md)). Rejected: it asks the operator to name each runtime's name for the host, the container adapter would have to keep it in step with the URL it hands a run, and it adds no gate the token does not already give.
- **Recording the decision in `design/system/flai-cli.md` alone.** Rejected: switching off a protection the MCP specification recommends for local servers is a security posture, which belongs in an ADR.
