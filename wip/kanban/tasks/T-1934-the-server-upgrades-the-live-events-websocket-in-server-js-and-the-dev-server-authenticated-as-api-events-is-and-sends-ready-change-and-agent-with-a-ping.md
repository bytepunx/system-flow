---
id: T-1934
type: task
nature: remediation
title: The server upgrades the live-events WebSocket in server.js and the dev server, authenticated as /api/events is, and sends ready, change, and agent with a ping
status: backlog
parent: S-0454
owner: alex
created: 2026-10-20T22:21:05Z
updated: 2026-10-20T22:24:37Z
transitions: []
stream: S-0454
tags: [dashboard, server-side]
touches: [flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/hooks.server.ts, flaiover/server.js, flaiover/vite.config.ts]
after: [T-1931]
---
# T-1934 The server upgrades the live-events WebSocket in server.js and the dev server, authenticated as /api/events is, and sends ready, change, and agent with a ping

## Work

Add `flaiover/src/lib/server/events-socket.ts`, the handler of the live-events upgrade, with `ws`'s `WebSocketServer` in `noServer` mode.

- **Routing.** `flaiover/server.js` (lines 27 to 31) and the plugin `flaiover-agent-upgrade` in `flaiover/vite.config.ts` (lines 34 to 47) hand every upgrade to `globalThis.__flaioverAgentUpgrade`, which `exposeAgentUpgrade()` in `flaiover/src/lib/server/agent.ts` sets and `init` in `flaiover/src/hooks.server.ts` calls. Make that one entry a router by path: the events path goes to the new handler, `/agent` to `AgentRegistry.handleUpgrade` as now. The router must run ahead of the registry, because `handleUpgrade` answers 404 to any other path (`agent.ts` line 749) and 403 to any request with an `Origin` header (line 751), and a browser always sends one. The plugin's filter on `/agent` (line 39) widens to both paths.
- **Who may.** Upgrades bypass SvelteKit's `handle`, so the handler builds a `Headers` from `req.headers` and calls `authenticate` and `decide` from `flaiover/src/lib/server/auth.ts` with the path `/api/events` and `GET`, as `hooks.server.ts` does (lines 102 to 114). `unauthorized` is answered 401 and `forbidden` 403 on the raw socket, before `handleUpgrade`, and the socket is destroyed. Check the `Origin` against the request's own host, since a cookie rides along on a cross-site WebSocket and there is no CSRF rule for a GET.
- **What it sends.** Inside `withProject` for the URL's `project` parameter, call `repo().watch()` and subscribe to the `Repo` as `flaiover/src/lib/server/sse.ts` does: `ready` once, `change` with `{ path, kind }`, `agent` with `{ story }`, each a JSON message with a `type`. Ping every 25 s, as `sse.ts` does, and close a socket that misses its pong. When the socket closes, `r.off` both listeners and clear the timer.
- **Shutting down.** `shutdown()` in `server.js` closes the HTTP connections after `SHUTDOWN_TIMEOUT` so an open event stream cannot keep a stopped server alive (I-0025). An upgraded socket is no longer the HTTP server's to close, so close every events socket with code 1001 when the server stops, and say so in the comment there.
- **Comment.** `hooks.server.ts` line 50 names `/api/events`; keep it true.

Test it in `flaiover/src/lib/server/events-socket.test.ts` on a real HTTP server, as `serve(registry)` in `flaiover/src/lib/server/agent.test.ts` (lines 50 to 56) does: an unauthenticated upgrade is refused with 401 and a forbidden one with 403 before the socket opens; an authenticated one gets `ready`, then a `change` and an `agent` the `Repo` emits; closing the socket removes the listeners; `/agent` still reaches the registry.

It waits for T-1931, whose ADR fixes the path and the messages. It runs beside the client and harness tasks: no path in common.

## Done when

- The upgrade works under `node server.js` and under `vite dev`.
- The tests above pass in `flai test` on the changed paths.
- `/agent` behaves as before; `agent.test.ts` passes unchanged.
- A server asked to stop with an events socket open exits without waiting out `SHUTDOWN_TIMEOUT`.

## Notes

The server-sent events endpoint stays until the last task removes it, so the dashboard works between tasks.
