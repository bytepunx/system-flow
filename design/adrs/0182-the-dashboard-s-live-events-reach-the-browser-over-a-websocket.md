---
id: ADR-0182
title: "The dashboard's live events reach the browser over a WebSocket"
status: accepted
date: 2026-10-21
supersedes: []
superseded_by: []
refines: [ADR-0029, ADR-0172]
---

# ADR-0182 The dashboard's live events reach the browser over a WebSocket

## Context

Each dashboard tab follows a project's live events with one `EventSource` on `/api/events`, one per project per tab since S-0154. The dashboard serves HTTP/1.1, and Chromium-based browsers keep at most six connections to one host. A server-sent event stream holds its connection for as long as the tab is open. Once six streams hold all six, every other request from every tab waits in the browser's queue, with no timeout.

On 2026-10-20 Brave at 127.0.0.1:4242 hung in every tab. The flaiover log showed every request handled in 0 to 145 ms, and no request from the browser after its last `/api/events` at 21:44:05. `ss` showed exactly six browser connections to port 4242, each of which had sent its one request 7 to 42 minutes earlier and since then only received the server's pings.

A WebSocket does not count toward the browser's six HTTP/1.1 connections to a host.

## Decision

The dashboard's live events reach the browser over a WebSocket instead of server-sent events.

- **Path.** The socket is an upgrade of `GET /api/events?project=<key>`, handled in `flaiover/server.js` and in the dev server's upgrade plugin in `flaiover/vite.config.ts`, beside the `/agent` upgrade ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). The one upgrade entry both hand to becomes a router by path, and it routes the events path ahead of the agent registry, which keeps refusing any request that carries an `Origin`.
- **Who may.** Upgrades bypass SvelteKit's hooks, so the events upgrade is authenticated and authorized as a `GET` of `/api/events` is in `flaiover/src/hooks.server.ts`, with the same functions. An unauthenticated upgrade is answered 401 and a forbidden one 403 before the socket opens. Since a browser sends the dashboard's cookie on a cross-site WebSocket too, an upgrade whose `Origin` is not the request's own host is answered 403.
- **Messages.** The server sends JSON text messages, each with a `type`, for the project the URL's `project` parameter names: `ready` once, when it listens to the project's `Repo`; `change` with `{ path, kind }`; and `agent` with `{ story }`. They carry what the server-sent events of the same names carried.
- **Liveness.** The server pings every 25 seconds and closes a socket that missed the pong of the ping before. When the socket closes, it stops listening to the project's `Repo`. When the server stops, it closes every events socket with code 1001, so that an open socket never holds a stopped server (I-0025).
- **Client.** `flaiover/src/lib/events.ts` keeps its `listen` API and one socket per project per tab, opened by the first listener and closed with the last. After a drop it reconnects with a backoff while listeners remain. A WebSocket does not resume as an `EventSource` did, so on reconnecting it tells every listener to refetch, through the callbacks listeners have now.
- **Harness.** The component harness treats a WebSocket to the app's own server as it treats an `EventSource` ([ADR-0172](0172-the-component-harness-answers-a-framed-component-s-requests-to-the-app-s-own.md)): it opens no connection, fails as a refused one does, and posts an `event` named `WebSocket` with the path. Vite's own sockets go through.
- **What goes.** The server-sent events endpoint, `flaiover/src/routes/api/events/+server.ts`, and `flaiover/src/lib/server/sse.ts` are removed: nothing else uses them. A plain `GET /api/events` then answers 404.

## Consequences

- Any number of tabs can follow their projects without holding the browser's connections to the dashboard, so no tab's request waits on another tab's events.
- The streaming responses of `/api/items/[id]/accept`, `/api/items/[id]/checks/tail`, and `/api/import` still hold a connection each, but only while they run.
- The dev server and the production server share one upgrade router; a third upgrade path is a route in it, not another global.
- A proxy in front of the dashboard must pass WebSocket upgrades on `/api/events` as it does on `/agent`.
- A Playwright test in `flaiover/e2e` opens seven pages on one project's events and shows a fetch from the last still answers.
- [flaiover-dashboard.md](../system/flaiover-dashboard.md) and [dashboard-host-channel.md](../system/dashboard-host-channel.md) describe the socket.

## Alternatives considered

The operator weighed these on 2026-10-20.

- **Close the stream while a tab is hidden.** Fewer streams, but six visible tabs, or six windows, still hang, and a hidden tab comes back stale.
- **One stream shared across tabs**, through a `SharedWorker` or a `BroadcastChannel` leader. One connection per browser, but a second browser profile or a private window adds its own, and leader election is code to keep right.
- **HTTP/2 behind local TLS.** HTTP/2 multiplexes every request over one connection, but browsers speak it only over TLS, so the dashboard would need a certificate the browser trusts on 127.0.0.1 and on the LAN. Deferred, not rejected: it would also free the three streaming responses above.
