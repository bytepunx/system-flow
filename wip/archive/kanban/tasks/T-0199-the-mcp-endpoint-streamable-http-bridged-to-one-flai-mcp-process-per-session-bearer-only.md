---
id: T-0199
type: task
nature: feature
title: "The /mcp endpoint: Streamable HTTP bridged to one flai mcp process per session, bearer only"
status: done
parent: S-0043
owner: alex
created: 2026-09-19T05:17:07Z
updated: 2026-09-19T05:22:30Z
transitions:
  - to: ready
    at: 2026-09-19T05:20:18Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T05:20:18Z
    by: system-flow
  - to: done
    at: 2026-09-19T05:22:30Z
    by: system-flow
stream: S-0043
tags: []
touches: [flaiover/src/lib/server, flaiover/src/routes]
---

# T-0199 The /mcp endpoint: Streamable HTTP bridged to one flai mcp process per session, bearer only

## Work
`src/lib/server/mcpbridge.ts`: a session is a spawned `flai mcp` in the project directory with `FLAI_AGENT` from the `X-Flai-Agent` header, else the client's name from `initialize`, and the usual flai environment; it writes JSON-RPC lines to stdin and matches responses on stdout by id. `src/routes/mcp/+server.ts`: POST without a session and with an `initialize` request starts one and returns the response with `Mcp-Session-Id`; POST with a session forwards: a request waits for its response and returns it as `application/json`, notifications and responses return 202; an unknown session is 404, a missing one 400; batches are handled; GET is 405; DELETE ends the session. Sessions idle out and are capped; the child is killed when its session ends and when the server stops. In `auth.ts`, `/mcp` allows only bearer (or auth off), never the cookie, and refuses a cross-origin `Origin`. Tests with the built flai: initialize, tools/list, a tool call, two sessions with different agents as two processes, 401 without the token and with only the cookie, 404 after DELETE, the cap.

## Done when
- The tests pass; lint and svelte-check are clean
- A production build succeeds

## Notes
