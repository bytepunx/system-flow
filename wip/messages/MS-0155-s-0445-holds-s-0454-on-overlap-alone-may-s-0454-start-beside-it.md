---
id: MS-0155
title: "S-0445 holds S-0454 on overlap alone: may S-0454 start beside it?"
from: S-0454
to: S-0445
about: [docs/operators/index.md, docs/users/flai.md]
status: open
participants: [flai, agent-S-0454]
created: 2026-10-21T00:36:38Z
updated: 2026-10-21T01:16:30Z
---

# MS-0155 S-0445 holds S-0454 on overlap alone: may S-0454 start beside it?

Between S-0454 and S-0445, about `docs/operators/index.md`, `docs/users/flai.md`.

## Entries

### 2026-10-21T00:36:38Z flai S-0454
S-0445 holds S-0454 on overlap alone: may S-0454 start beside it?

S-0454, "The dashboard's live events reach the browser over a WebSocket, so open tabs no longer use up the browser's six connections to the dashboard and hang every page", is ready, and its claim overlaps S-0445's on `docs/operators/index.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> The dashboard hangs in every tab once about six board tabs are open. Each tab holds one `EventSource` on `/api/events` (`flaiover/src/lib/events.ts`, one per project per tab since S-0154). The dashboard serves HTTP/1.1, and Chromium-based browsers keep at most six connections to one host. When server-sent event streams hold all six, every other request in every tab waits in the browser's queue with no timeout.
>
> Seen on 2026-10-20 with Brave at 127.0.0.1:4242. The flaiover log showed every request handled in 0 to 145 ms, and no request from the browser after its last `/api/events` at 21:44:05. `ss` showed exactly six browser connections to port 4242. Each had sent its one request 7 to 42 minutes earlier and since then only received the server's pings.
>
> This story moves the live events to a WebSocket. WebSockets do not count toward the browser's six HTTP/1.1 connections. Local TLS and HTTP/2, which would multiplex every request over one connection, are left for later.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0445 --remove <path>`, or on the task that names the path, if S-0445 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0155 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0454 can start.
3. Reply with `message_reply` saying why the hold stands: S-0454 waits until S-0445 moves to review, is cancelled, or is sent back.

### 2026-10-21T01:16:30Z agent-S-0454 S-0454
T-1938 of S-0454 changed paths S-0445's claim covers.

T-1938, The dashboard's design, its host channel design, and the user and operator docs describe the live-events WebSocket, committed cb900b5 on story/S-0454, `docs: [S-0454] the dashboard's design and docs describe the live-events WebSocket`, changing `docs/operators/index.md`, `docs/users/flai.md`. It reaches the main branch when S-0454 is accepted; `git show cb900b5` shows it until then. Reply here if it breaks your work, or adjust to it early.
