---
id: T-0197
type: task
nature: feature
title: "ADR-0024 and living design: MCP over HTTP at /mcp, project identity in the API, and the hub shape"
status: done
parent: S-0043
owner: alex
created: 2026-09-19T05:17:07Z
updated: 2026-09-19T05:17:59Z
transitions:
  - to: ready
    at: 2026-09-19T05:17:08Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:17:08Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:17:59Z
    by: system-flow
stream: S-0043
tags: []
touches: [design/adrs, design/system]
---

# T-0197 ADR-0024 and living design: MCP over HTTP at /mcp, project identity in the API, and the hub shape

## Work
Write ADR-0024, refining ADR-0020 (files plus MCP) and ADR-0018 (the token): flaiover exposes the same MCP server over Streamable HTTP at `/mcp`; each session is its own `flai mcp` process, bridged over stdio, so the rules stay in flai (ADR-0016); only the bearer token authenticates it, never the session cookie, and a cross-origin `Origin` is refused; responses are plain JSON, a GET stream is not offered (405), sessions end on DELETE or after idling, and their number is capped. Every `/api/*` and `/mcp` response names the project: `project: { name, key }` in JSON object bodies and `X-Flai-Project-Key` and `X-Flai-Project-Name` headers on every response, because array, streamed, and error responses have no place for a member; the manifest's `key` becomes required, a `flai check` warning now and an error after the template minor that ships it. In `design/system/flaiover-dashboard.md` add `/mcp`, the identity, and a "Hub shape" note: flaiover dials out to the hub over a websocket with its token, no inbound ports, the hub routes by project key; nothing here builds it. Update `flai-cli.md` for the check rule.

## Done when
- ADR-0024 is accepted and indexed
- The design documents describe `/mcp`, the identity, and the hub shape
- `flai check --strict` is clean

## Notes
