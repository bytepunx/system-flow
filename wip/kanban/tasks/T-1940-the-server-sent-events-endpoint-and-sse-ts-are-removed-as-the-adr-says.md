---
id: T-1940
type: task
nature: remediation
title: The server-sent events endpoint and sse.ts are removed, as the ADR says
status: backlog
parent: S-0454
owner: alex
created: 2026-10-20T22:22:45Z
updated: 2026-10-20T22:22:45Z
transitions: []
stream: S-0454
tags: [dashboard, server-side]
touches: [flaiover/src/routes/api/events/+server.ts, flaiover/src/lib/server/sse.ts, flaiover/src/lib/server/sse.test.ts]
after: [T-1939]
---
# T-1940 The server-sent events endpoint and sse.ts are removed, as the ADR says

## Work

Delete `flaiover/src/routes/api/events/+server.ts`, `flaiover/src/lib/server/sse.ts`, and `flaiover/src/lib/server/sse.test.ts`, unless T-1931's ADR keeps one of them for a reason it names. Only the route and its test import `sse.ts`; the streaming responses of `/api/import`, `/api/items/[id]/accept`, and `/api/items/[id]/checks/tail` write `application/x-ndjson` of their own and do not use it. A plain `GET /api/events` then answers 404 from SvelteKit, while the upgrade on that path reaches the handler of T-1934 before SvelteKit sees it; if the ADR wants 426 Upgrade Required there instead, keep the route as that.

Search `flaiover/src` for `sse`, `text/event-stream`, and `EventSource` once more: nothing but tests should name them.

It waits for T-1939, whose run against the server-sent events needs them still in place.

## Done when

- The three files are gone, or kept as the ADR says.
- `flai test --all` on `flaiover` passes, svelte-check included, and `scripts/flaiover-e2e.sh` still passes.

## Notes

The docs of T-1938 already describe the state after this task.
