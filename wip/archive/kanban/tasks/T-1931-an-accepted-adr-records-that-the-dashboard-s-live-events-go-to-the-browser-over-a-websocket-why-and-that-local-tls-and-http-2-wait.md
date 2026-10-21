---
id: T-1931
type: task
nature: remediation
title: An accepted ADR records that the dashboard's live events go to the browser over a WebSocket, why, and that local TLS and HTTP/2 wait
status: done
parent: S-0454
owner: alex
created: 2026-10-20T22:20:43Z
updated: 2026-10-21T00:53:26Z
transitions:
  - to: ready
    at: 2026-10-21T00:53:19Z
    by: agent-S-0454
  - to: in-progress
    at: 2026-10-21T00:53:19Z
    by: agent-S-0454
  - to: done
    at: 2026-10-21T00:53:26Z
    by: agent-S-0454
stream: S-0454
tags: [dashboard, adr]
touches: [design/adrs/README.md, design/adrs/0182-the-dashboard-s-live-events-reach-the-browser-over-a-websocket.md]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 851
      cache_read: 129283
      cache_write: 3434
      cost: 0.0636
---
# T-1931 An accepted ADR records that the dashboard's live events go to the browser over a WebSocket, why, and that local TLS and HTTP/2 wait

## Work

Write and accept the ADR with `flai adr new` and `flai adr accept`, titled "The dashboard's live events reach the browser over a WebSocket", so that its file is the one this task touches. If another story took 0179 first, keep the number flai gives and widen the touches with `flai touches`.

The ADR refines ADR-0029 and ADR-0172. It records:

- The context: one `EventSource` on `/api/events` per project per tab (S-0154), HTTP/1.1, and Chromium's six connections to one host. On 2026-10-20 Brave at 127.0.0.1:4242 held six streams, and every other request in every tab waited, while the server answered in 0 to 145 ms.
- The decision: the socket's path, `/api/events` upgraded in `flaiover/server.js` and the dev server's plugin beside `/agent`, unless the design finds a reason for another. Its messages are `ready` once, `change` with `{ path, kind }`, and `agent` with `{ story }`, as JSON with a `type`, for the project the URL's `project` parameter names. It pings to keep the connection alive. The upgrade is authenticated and authorized as a request to `/api/events` is, with 401 or 403 before the socket opens. The events path is routed ahead of the agent registry, which keeps refusing any request with an `Origin`, and the events upgrade checks its `Origin` against the request's own host (TH-0538). The client reconnects with a backoff and tells its listeners to refetch. The harness treats the socket as it treats an `EventSource`.
- What goes: the server-sent events endpoint and `flaiover/src/lib/server/sse.ts`, unless something still needs them.
- What waits: local TLS and HTTP/2. The three streaming responses of `/api/items/[id]/accept`, `/api/items/[id]/checks/tail`, and `/api/import` stay, since they hold a connection only while they run.
- The alternatives the operator weighed on 2026-10-20: closing the stream while a tab is hidden, one stream shared across tabs, and HTTP/2 behind local TLS.

It is the first task because every other task builds what it decides.

## Done when

- The ADR is accepted and its row is in `design/adrs/README.md`.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

`design/adrs` is in the manifest's shared paths, so this task holds no other story there.
