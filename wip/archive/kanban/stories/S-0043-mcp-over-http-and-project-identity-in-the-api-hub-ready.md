---
id: S-0043
type: story
nature: feature
title: MCP over HTTP and project identity in the API, hub-ready
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-19T05:29:27Z
transitions:
  - to: ready
    at: 2026-09-19T01:56:20Z
    by: alex
  - to: in-progress
    at: 2026-09-19T05:15:03Z
    by: system-flow
  - to: review
    at: 2026-09-19T05:27:19Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:29:27Z
    by: alex
tags: [dashboard, cli]
---

# S-0043 MCP over HTTP and project identity in the API, hub-ready

## Goal
Remote agents and the future multi-project hub reach the same MCP server through flaiover over HTTP with the bearer token, and every API response identifies the project, so a web front end spanning projects can route and label without guessing.

## Acceptance criteria
- [x] flaiover exposes MCP over Streamable HTTP at `/mcp`, authenticated with the bearer token, delegating to `flai mcp` per ADR-0016; sessions map to one process per client
- [x] Every `/api/*` response carries `project: { name, key }` from the manifest; the manifest's `key` becomes required by `flai check` (warning first, error after the template minor that ships it)
- [x] A design note in flaiover-dashboard.md describes the hub shape: flaiover dials out to the hub over a websocket with its token, no inbound ports; nothing in this story builds the hub
- [x] Tests for the HTTP transport with a client; docs/operators documents `/mcp` and the tunnel expectation (TLS terminated by the tunnel or proxy)

## Tasks
- T-0197 ADR-0024 and living design: MCP over HTTP at /mcp, project identity in the API, and the hub shape
- T-0198 Project identity on every API response, and flai check asks for the manifest key
- T-0199 The /mcp endpoint: Streamable HTTP bridged to one flai mcp process per session, bearer only
- T-0200 An official MCP client test against a running dashboard
- T-0201 Operator and user documentation for /mcp, the tunnel expectation, and project identity
- T-0202 Verify: all tiers, and the official client against a real container

## Notes
Last in the epic. Shape agreed 2026-09-17 so the hub is cheap later: bearer first, project identity in responses, outbound connections.

Decided when pulled, 2026-09-19 (ADR-0024). `/mcp` is a bridge in flaiover's server between Streamable HTTP and one `flai mcp` stdio process per session; responses are plain JSON, a GET stream is not offered (405), sessions idle out and are capped, and only the bearer token opens it, never the cookie, with a cross-origin `Origin` refused. Project identity is a `project` member in JSON object bodies and two headers on every `/api` and `/mcp` response. Array responses (`/api/items`, `/api/threads`, `/api/docs/tree`, `/api/docs/adrs`) keep their shape and are identified by the headers only: that is where this story stops short of the second criterion's wording, deliberately, and the narrative's open question proposes deciding envelopes with the hub. Refusals for lack of a token carry no identity. `manifest.key` is a `flai check` warning now; making it an error waits for the template minor after this one, as the criterion says.

Verification, 2026-09-19. `make flai-test` passed (golangci-lint 0 issues; behavior, integration, smoke; markdown lint); flaiover lint, svelte-check 0 errors, 28 files and 174 tests, and a production build passed. The official Go SDK client (`StreamableClientTransport`, `flai/tests/integration/mcp_http_test.go`) connected with the bearer token, listed the tools, and called `inbox` and `board`, first against the dev server and then against a container built from the branch on a scratch project; without the token it got 401. Against that container: a session opened with `X-Flai-Agent: remote-claude` ran as that agent and saw the ready story through `inbox`; exactly one `flai mcp` process existed in the container while the session was open and none after `DELETE` (204); the log had `mcp session started` and `mcp session ended`; a foreign `Origin` got 403, `GET /mcp` 405, the cookie alone 401; `/api/board` carried both headers and `project` in its body, `/api/items` carried the headers with its array unchanged. The operator's dashboard was not touched. Not exercised: the idle timeout firing on its own (thirty minutes by default; the timer's code path is the same `close` the `DELETE` test covers), and a client behind a TLS-terminating proxy, which is documented for operators and not something this repository sets up.
