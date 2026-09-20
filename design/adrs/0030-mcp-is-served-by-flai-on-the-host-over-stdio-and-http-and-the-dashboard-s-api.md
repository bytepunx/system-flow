---
id: ADR-0030
title: "MCP is served by flai on the host, over stdio and HTTP, and the dashboard's API still names its project"
status: accepted
date: 2026-09-20
supersedes: [ADR-0024]
superseded_by: []
refines: [ADR-0020, ADR-0029]
---

# ADR-0030 MCP is served by flai on the host, over stdio and HTTP, and the dashboard's API still names its project

## Context

[ADR-0024](0024-mcp-over-http-and-project-identity.md) put MCP over HTTP in flaiover: `/mcp` started one `flai mcp` process per session inside the dashboard's container and carried JSON-RPC to it. That needed a flai binary in the image and the project's files mounted into the container.

[ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md) takes both away. The dashboard reaches a project only by asking flai on the host; since S-0075 the image holds no flai, and the mount is next to go. `/mcp` has answered 503 since then.

Relaying `/mcp` to the host over the channel was tried in S-0071 and works. The operator decided against it on 2026-09-20: "lets move the MCP integration down to flai as a process again (flai MCP can start a server process) since we want flai and agents working together on the host and not having agents interacting with a project through the dashboard."

The protocol is also moving. Revision 2026-07-28 removes sessions, the GET stream, resumption, and requests from server to client. The Go SDK (1.8.0) serves that revision over HTTP only from a handler without sessions, and serves the revisions before it, which today's clients speak, with sessions.

## Decision

**MCP is flai's, on the host, on two transports.** `flai mcp` serves stdio as before, and `.mcp.json` does not change. `flai mcp http` serves the same server over Streamable HTTP at `/mcp`: the tools, the rules, and the agent's cursor are one implementation (`internal/mcpserver`), and `internal/mcphttp` decides only who may ask, under what name, and how many at once.

**A process of its own, one per project.** `flai mcp start` runs `flai mcp http` detached, `flai mcp stop` ends it by its recorded PID, and `flai mcp status` says where it listens and prints an agent's configuration. Its state (`mcp-http.json`, rewritten every few seconds so a dead server's file is known to be stale), its log, its token, and the address it last used are files in the main checkout's `.flai-cache`. It is not part of `flai serve`: that process exists to dial dashboards, and an agent's server must work with no dashboard at all.

- It listens on `127.0.0.1:4243` unless `--addr` says otherwise, and the address is remembered per project so an agent's configuration survives a restart. A second project on the same machine chooses another port. An address beyond the machine is allowed and warned about; TLS remains a proxy's or a tunnel's job.
- Only a bearer token authenticates it: `.flai-cache/mcp.token`, mode 0600, created when first needed, printed by `flai mcp token`, replaced by `--rotate`, which restarts a running server. It is not the dashboard's token. A request with an `Origin` header is refused before the token is looked at: agents send none, and a page in a browser must not reach this server.
- The agent is named by the `X-Flai-Agent` header, else by the client's own name (from `initialize`, or from the request's `_meta` where there are no sessions), made safe for a file name, else `agent`.
- Both generations of the transport are served at the one address, chosen by the `Mcp-Protocol-Version` header. Up to 2025-11-25 a client initializes a session, which ends on DELETE or after it idles (30 minutes) and is capped (16). From 2026-07-28 there are no sessions and each request stands alone. Requests in flight are capped (64) on both. Answers are `application/json`, not event streams. The GET stream is whatever the SDK offers for the revision; flai sends nothing unprompted, and an agent waits with `wait_for_events`, which holds its POST open.

**The dashboard serves no MCP.** `/mcp`, the bridge, and their settings are removed from flaiover. Any request to `/mcp` answers 410 with a JSON-RPC error that says where MCP lives now. It needs no token, because it says nothing about the project.

**Project identity, carried forward unchanged from ADR-0024.** Every `/api/*` response from the dashboard names the project from the manifest: JSON object bodies gain `project: { name, key }`, and every response carries the headers `X-Flai-Project-Key` and `X-Flai-Project-Name` (the name URI-encoded). The manifest's `key` is required. flai's MCP server over HTTP sends the same two headers on every answer, refusals included, since its caller already chose the project by choosing the address.

The hub shape sketched in ADR-0024 is not carried forward: ADR-0029 decided how a dashboard and a host reach each other.

## Consequences

- An agent on the host needs nothing new: `.mcp.json` starts `flai mcp` on stdio as before.
- An agent elsewhere configures one URL, one token, and its name, and the operator decides how that address is reached (an SSH forward, a tunnel); the dashboard is not in the path and need not run.
- A connected agent no longer costs a process. Sessions cost memory only; the caps bound held requests instead of processes.
- The dashboard's token no longer grants MCP. Someone who configured an agent against `/mcp` is told by the 410 what to run.
- Two tokens per project instead of one. They open different things, and rotating one does not disturb the other.
- The server reads the manifest when it starts and serves the main checkout, so a document read over HTTP is main's, where stdio in a worktree reads that worktree's.
- When clients stop speaking the older revisions, the session half of `mcphttp` can go without changing anything an agent configured.
- `wait_for_events` still holds a request for up to five minutes; a proxy in front must allow it.

## Alternatives considered

- **MCP inside `flai serve`, one listener for every project.** One port and one token per user, and no port to choose per project. But `flai serve` would gain a listener it does not otherwise need, agents would depend on the process that exists for dashboards, and an agent's address would carry a project key that stdio never needed. Worth revisiting if one dashboard for all projects (S-0080) makes per-project ports a nuisance.
- **Keep `/mcp` and relay it over the channel to the host.** Tried in S-0071; it works. Declined by the operator: agents belong with flai on the host, not behind the dashboard.
- **Sessions only.** What ADR-0024 had. The SDK then refuses 2026-07-28 and clients fall back, which works today and stops working when a client drops the older revisions.
- **No sessions only.** Simplest, and where the protocol is going. But a client on an older revision names itself only in `initialize`, so without the header every such agent would be `agent` and share one cursor.
- **Reuse the dashboard's token.** One secret to hand out, but it would make the dashboard's login able to act as any agent, and rotate both together.
