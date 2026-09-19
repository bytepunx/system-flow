---
id: ADR-0024
title: flaiover serves MCP over HTTP, and every response names its project
status: accepted
date: 2026-09-19
supersedes: []
superseded_by: []
---

# ADR-0024 flaiover serves MCP over HTTP, and every response names its project

## Context

[ADR-0020](0020-files-plus-mcp.md) gives agents a view of the repository through `flai mcp`, a server on standard input and output. That reaches an agent running on the machine that holds the repository and nobody else. An agent on another machine, and the multi-project hub E-0006 is shaped for, need the same server over the network. flaiover is already the project's network face: it listens on one port and authenticates every request with the project token ([ADR-0018](0018-dashboard-token.md)), and it already hands every write to the bundled `flai` ([ADR-0016](0016-dashboard-delegates-to-flai.md)).

A front end that spans projects also has to know which project an answer came from. Today nothing in a response says so; a client knows only the address it called.

The shape was agreed on 2026-09-17 so that a hub is cheap to add later: bearer first, project identity in responses, outbound connections.

## Decision

**MCP over HTTP.** flaiover exposes the MCP server at `/mcp` using the Streamable HTTP transport. It does not implement the tools: each session is its own `flai mcp` process in the project directory, and flaiover carries JSON-RPC messages between the HTTP request and that process's standard input and output. The rules, the tools, and the agent's cursor are therefore flai's, the same as on stdio.

- A session starts with a POST carrying an `initialize` request and no session header; the response carries `Mcp-Session-Id`, which every later request repeats. A request waits for its response, which is returned as `application/json`; notifications and responses are acknowledged with 202. An unknown session is 404, and the client starts again.
- A GET stream is not offered (405, which the transport allows): `flai mcp` sends nothing unprompted that a client needs, and an agent waits with `wait_for_events`, which holds its POST open.
- A session ends on DELETE, or after it has idled, and the process is killed with it. The number of sessions is capped, so a client that leaks sessions cannot fill the host with processes.
- The agent's name for the session is the `X-Flai-Agent` header on `initialize`, else the client's name from `initialize`.
- Only the bearer token authenticates `/mcp`. The session cookie does not, so a page in the designer's browser cannot drive it, and a request whose `Origin` is another site is refused.
- TLS is not flaiover's job. Beyond a trusted network the dashboard sits behind a tunnel or proxy that terminates TLS, because the token travels in the request.

**Project identity.** Every `/api/*` and `/mcp` response names the project from the manifest: JSON object bodies gain `project: { name, key }`, and every response, whatever its body, carries the headers `X-Flai-Project-Key` and `X-Flai-Project-Name` (the name URI-encoded). Array, streamed, and empty responses have no place for a member, and the headers cover them without changing their shape. The manifest's `key` therefore becomes required: `flai check` warns with `manifest.key` when it is missing, and the warning becomes an error after the template minor that ships this.

**Hub shape, not built here.** A hub that fronts several projects is reached by flaiover, not the other way round: each flaiover dials out to the hub over a websocket with its token and keeps the connection, so a project needs no inbound port, and the hub routes requests to a connection by project key and labels answers from the identity above.

## Consequences

- A remote agent configures one URL and one token and gets the tools it would have had locally, including `inbox` and `wait_for_events`.
- Each connected agent costs a process on the host for as long as its session lives. The cap and the idle timeout bound that; both are settings.
- `wait_for_events` can hold a request open for up to five minutes. A proxy in front of the dashboard must allow that, which the operator documentation says.
- A client that wants server-initiated messages gets none. If `flai mcp` ever needs to push, the GET stream is the place, and this decision is revisited.
- Array responses identify their project only in headers. A client that reads bodies alone does not see it there. Turning those arrays into objects would break the API for its one consumer today; it is left for the hub story to decide.
- The dashboard token now also grants everything the MCP tools can do, which is less than the API already allows (moves, thread replies, reads); acceptance stays refused over MCP.

## Alternatives considered

- `flai mcp --http`, one process serving every session itself: fewer processes, but a second listener and a second place that authenticates, and the criterion of one process per client, which keeps an agent's identity and cursor apart, is lost.
- Implement the tools again in flaiover's server: two implementations of the rules, which ADR-0016 exists to prevent.
- Accept the session cookie on `/mcp`: lets any page the designer's browser loads attempt tool calls with their session.
- Wrap every array response as `{ project, items }` now: a breaking change to every list endpoint for a consumer that does not exist yet.
- Have the hub dial in to each project: needs an inbound port and a reachable address per project, which is what the tunnel expectation avoids.
