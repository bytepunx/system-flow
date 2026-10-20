---
id: S-0454
type: story
nature: remediation
title: The dashboard's live events reach the browser over a WebSocket, so open tabs no longer use up the browser's six connections to the dashboard and hang every page
status: backlog
owner: alex
created: 2026-10-20T22:10:33Z
updated: 2026-10-20T22:10:33Z
transitions: []
tags: [dashboard]
topics: [client-side, server-side]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0454 The dashboard's live events reach the browser over a WebSocket, so open tabs no longer use up the browser's six connections to the dashboard and hang every page

## Goal

The dashboard hangs in every tab once about six board tabs are open. Each tab holds one `EventSource` on `/api/events` (`flaiover/src/lib/events.ts`, one per project per tab since S-0154). The dashboard serves HTTP/1.1, and Chromium-based browsers keep at most six connections to one host. When server-sent event streams hold all six, every other request in every tab waits in the browser's queue with no timeout.

Seen on 2026-10-20 with Brave at 127.0.0.1:4242. The flaiover log showed every request handled in 0 to 145 ms, and no request from the browser after its last `/api/events` at 21:44:05. `ss` showed exactly six browser connections to port 4242. Each had sent its one request 7 to 42 minutes earlier and since then only received the server's pings.

This story moves the live events to a WebSocket. WebSockets do not count toward the browser's six HTTP/1.1 connections. Local TLS and HTTP/2, which would multiplex every request over one connection, are left for later.

## Acceptance criteria
- [ ] An accepted ADR records that the dashboard's live events go to the browser over a WebSocket instead of server-sent events, why (the six-connection limit and the 2026-10-20 instance), and that local TLS and HTTP/2 are deferred.
- [ ] The server upgrades the live-events WebSocket in `flaiover/server.js` and in the dev server's upgrade plugin in `flaiover/vite.config.ts`, beside the `/agent` upgrade (ADR-0029). It sends the messages `/api/events` sends now: `ready` once, `change` with `{ path, kind }`, and `agent` with `{ story }`, for the project the URL's `project` parameter names. It keeps the connection alive with a ping, and stops listening to the project's `Repo` when the socket closes.
- [ ] The upgrade is authenticated and authorized as a request to `/api/events` is in `flaiover/src/hooks.server.ts`. Upgrades bypass SvelteKit's hooks, so an unauthenticated or forbidden upgrade is refused with 401 or 403 before the socket opens, and a test shows it.
- [ ] `flaiover/src/lib/events.ts` keeps its `listen` API and still shares one socket per project per tab. After a drop it reconnects with a backoff, and on reconnecting it tells every listener to refetch, since a WebSocket does not resume as an `EventSource` did. No component that listens changes.
- [ ] A Playwright test reproduces the hang: at least seven pages of one browser context following the same project's events, where a fetch from the last page answers. Run against the server-sent events before the change, it hangs.
- [ ] The component harness (ADR-0172) treats a WebSocket to the app's own server as it treats an `EventSource` now: it opens no connection, fails as a refused one does, and lists it under Events.
- [ ] The server-sent events endpoint and `flaiover/src/lib/server/sse.ts` are removed, or kept only if something else still needs them, as the ADR says. `design/system/flaiover-dashboard.md`'s API table, `design/system/dashboard-host-channel.md`, and `docs/users/flaiover.md` describe the WebSocket.

## Tasks

## Notes

Chosen with the operator on 2026-10-20 over closing the stream while a tab is hidden, sharing one stream across tabs, and HTTP/2 behind local TLS.

`ws` is already a dependency, and `server.js` already routes `upgrade` to `globalThis.__flaioverAgentUpgrade` for `/agent`. That handler is the one to extend into a router by path.

Not in this story: the streaming responses of `/api/items/[id]/accept`, `/api/items/[id]/checks/tail`, and `/api/import` hold a connection too, but only while they run.
