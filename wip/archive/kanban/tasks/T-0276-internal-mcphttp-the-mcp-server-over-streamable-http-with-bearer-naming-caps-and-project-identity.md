---
id: T-0276
type: task
nature: feature
title: "internal/mcphttp: the MCP server over Streamable HTTP, with bearer, naming, caps, and project identity"
status: done
parent: S-0076
owner: alex
created: 2026-09-20T12:16:18Z
updated: 2026-09-20T12:25:24Z
transitions:
  - to: ready
    at: 2026-09-20T12:18:07Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:18:07Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:25:24Z
    by: system-flow
stream: S-0076
tags: []
---
# T-0276 internal/mcphttp: the MCP server over Streamable HTTP, with bearer, naming, caps, and project identity

## Work
One http.Handler around the SDK's StreamableHTTPHandler, the tools from mcpserver.New unchanged. Sessions for revisions up to 2025-11-25 (idle timeout, capped); stateless for 2026-07-28, which the SDK serves over HTTP only that way, chosen by the Mcp-Protocol-Version header. Bearer token compared in constant time; a request with an Origin header is refused; the agent is X-Flai-Agent, else the client's name (initialize, or the request's _meta), made safe as the bridge made it; requests in flight are capped; every answer carries X-Flai-Project-Key and X-Flai-Project-Name.

## Done when
- Tests with the SDK's client and with raw HTTP: a full session, the cursor kept per agent name, 401, 403, 404 for an unknown session, 503 at the cap, DELETE, identity headers
- What the SDK does for each revision is written in the narrative for design/tech

## Notes
