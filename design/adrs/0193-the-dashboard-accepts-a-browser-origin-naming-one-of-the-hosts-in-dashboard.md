---
id: ADR-0193
title: "The dashboard accepts a browser Origin naming one of the hosts in dashboard.allowed_hosts, which flai keeps in the host configuration and passes to the container as FLAIOVER_ALLOWED_HOSTS on each start"
status: proposed
date: 2026-10-23
supersedes: []
superseded_by: []
---

# ADR-0193 The dashboard accepts a browser Origin naming one of the hosts in dashboard.allowed_hosts, which flai keeps in the host configuration and passes to the container as FLAIOVER_ALLOWED_HOSTS on each start

## Context

The dashboard refuses one thing by host: the `/api/events` WebSocket upgrade answers 403 when the browser's `Origin` host is not the request's `Host` header (`flaiover/src/lib/server/events-socket.ts`). That check stops another site from opening an events socket with the dashboard's cookie. A tunnel or a reverse proxy can rewrite `Host` to the upstream address, or serve the dashboard under a domain of its own. Behind either, every page's live events are refused, and the page shows "Reconnecting…". The operator asked for a way to admit such hosts (S-0489).

## Decision

The dashboard accepts a browser `Origin` whose host is the request's own `Host` or one of the hosts in `dashboard.allowed_hosts`, which flai keeps in the host configuration and gives the container as `FLAIOVER_ALLOWED_HOSTS` each time it starts it.

- `dashboard.allowed_hosts` is in `~/.flai/config.json`: a comma-separated list of host names, each with an optional port, empty by default. It is a fact about one machine's shared container, so it is never in a project's manifest.
- `flai config set` refuses an entry that is not a host name with an optional port.
- `flai dashboard` passes the list, trimmed and lower-cased, as `FLAIOVER_ALLOWED_HOSTS` when it names any. That variable is the contract with the dashboard: a comma-separated list, lower-case, with no spaces and no scheme.
- A restart, an upgrade, and flai host's watch start the container with the list as the configuration has it then. A change therefore applies on the next start, and a running container keeps the list it was started with.
- The container record keeps the list. `flai dashboard status` reports it, and `settings.get` reports the configured one.
- The dashboard adds no refusal: a request with no `Origin` is served as before.
- The Host page edits the list through a host-wide settings write, behind the `settings` host action.

## Consequences

- An operator admits a tunnel or a domain with `flai config set dashboard.allowed_hosts <host>`, or from the Host page, then restarts the dashboard.
- An allowed host widens only which origins may open the events socket. Every request still needs the dashboard's token.
- Applying a change without a restart would need host settings carried over the channel, which nothing needs yet.

## Alternatives considered

- **Keep the list in the manifest.** Refused: one container serves every project on the host, so the list belongs to the host.
- **Drop the origin check behind a proxy.** Refused: it reopens cross-site socket use with the operator's cookie.
- **Trust `X-Forwarded-Host`.** Refused: a client can send that header itself. An explicit list states what the operator meant.
