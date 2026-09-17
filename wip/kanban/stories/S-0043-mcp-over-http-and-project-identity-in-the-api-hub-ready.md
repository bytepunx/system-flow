---
id: S-0043
type: story
nature: feature
title: MCP over HTTP and project identity in the API, hub-ready
status: backlog
parent: E-0006
owner: alex
created: 2026-09-17T19:46:29Z
updated: 2026-09-17T19:46:29Z
transitions: []
tags: [dashboard, cli]
---

# S-0043 MCP over HTTP and project identity in the API, hub-ready

## Goal
Remote agents and the future multi-project hub reach the same MCP server through flaiover over HTTP with the bearer token, and every API response identifies the project, so a web front end spanning projects can route and label without guessing.

## Acceptance criteria
- [ ] flaiover exposes MCP over Streamable HTTP at `/mcp`, authenticated with the bearer token, delegating to `flai mcp` per ADR-0016; sessions map to one process per client
- [ ] Every `/api/*` response carries `project: { name, key }` from the manifest; the manifest's `key` becomes required by `flai check` (warning first, error after the template minor that ships it)
- [ ] A design note in flaiover-dashboard.md describes the hub shape: flaiover dials out to the hub over a websocket with its token, no inbound ports; nothing in this story builds the hub
- [ ] Tests for the HTTP transport with a client; docs/operators documents `/mcp` and the tunnel expectation (TLS terminated by the tunnel or proxy)

## Tasks

## Notes
Last in the epic. Shape agreed 2026-09-17 so the hub is cheap later: bearer first, project identity in responses, outbound connections.
