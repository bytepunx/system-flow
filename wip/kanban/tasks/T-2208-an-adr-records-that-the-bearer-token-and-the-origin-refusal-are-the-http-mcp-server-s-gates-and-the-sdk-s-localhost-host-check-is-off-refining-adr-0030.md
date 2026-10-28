---
id: T-2208
type: task
nature: remediation
title: An ADR records that the bearer token and the Origin refusal are the HTTP MCP server's gates and the SDK's localhost Host check is off, refining ADR-0030
status: backlog
parent: S-0515
owner: alex
created: 2026-10-28T13:02:35Z
updated: 2026-10-28T13:02:35Z
transitions: []
stream: S-0515
tags: [design, mcp]
touches: [design/adrs/README.md]
---
# T-2208 An ADR records that the bearer token and the Origin refusal are the HTTP MCP server's gates and the SDK's localhost Host check is off, refining ADR-0030

## Work

Write an ADR with `flai adr new`, refining [ADR-0030](../../../design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md), that records the decision the plan on S-0515 proposes. Context: the go-sdk's `StreamableHTTPHandler` refuses with 403 `invalid Host header` a request on a loopback listener whose `Host` is not loopback, its protection against DNS rebinding; an agent in a container reaches the server at `127.0.0.1:4243` through Docker's host gateway with `Host: host.docker.internal:4243`, and was refused (I-0213, S-0505). Decision: flai's HTTP MCP server has two gates, the refusal of any request with an `Origin` header and the bearer token compared in constant time, both checked before the SDK sees the request, and the SDK's localhost protection is off on both handlers, so a request is admitted whatever `Host` it names once it carries the token; no allowed-hosts list, since the token already refuses the rebinding page the `Host` check guards against and a list would make the operator name each runtime's name for the host. Consequences: a container reaches the server through any name for the host; the token stays the one secret, and its exposure is as ADR-0030 and the operator guide state it; a listener beyond loopback is warned about as before. Give the ADR the topics `mcp` and `docker` with `flai adr topics`, and accept it with `flai adr accept`.

Waits for nothing: the ADR records the decision the code task follows, and the two touch no path in common. `flai adr new` names the ADR's file, which is why the story keeps the folder `design/adrs` in its touches; this task names only the index it updates.

## Done when

- A new ADR under `design/adrs/` names the two gates, the SDK's localhost protection off, and the rejected allowed-hosts list, with `refines` ADR-0030, topics `mcp` and `docker`, and status accepted
- `design/adrs/README.md` lists it, and `flai check --strict` is clean

## Notes

The dashboard's `dashboard.allowed_hosts` (ADR-0193) widens which `Origin` a browser may open the events socket from; it is the pattern I-0213 names as one alternative, and the ADR says why the MCP server takes the other.
