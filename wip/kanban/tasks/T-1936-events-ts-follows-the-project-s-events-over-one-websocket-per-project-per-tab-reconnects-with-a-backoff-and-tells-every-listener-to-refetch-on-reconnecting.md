---
id: T-1936
type: task
nature: remediation
title: events.ts follows the project's events over one WebSocket per project per tab, reconnects with a backoff, and tells every listener to refetch on reconnecting
status: backlog
parent: S-0454
owner: alex
created: 2026-10-20T22:21:24Z
updated: 2026-10-20T22:21:24Z
transitions: []
stream: S-0454
tags: [dashboard, client-side]
touches: [flaiover/src/lib/events.ts, flaiover/src/lib/events.test.ts, flaiover/src/lib/inbox.svelte.test.ts, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte.test.ts]
after: [T-1931]
---
# T-1936 events.ts follows the project's events over one WebSocket per project per tab, reconnects with a backoff, and tells every listener to refetch on reconnecting

## Work

In `flaiover/src/lib/events.ts`, replace the `EventSource` of each stream with a `WebSocket` to the events path the ADR names, its URL made from `projectState.tag(...)` with `ws:` or `wss:` after the page's own protocol. Keep `listen`, `Listener`, `follow`, and `debounced` as they are, and keep one socket per project per tab, opened by the first listener and closed with the last.

- Parse each message's JSON `type`: `change` calls each listener's `change(path, kind)`, `agent` calls `agent(story)`, `ready` marks the socket open.
- After a drop, while listeners remain, reconnect with a backoff: from about 1 s, doubling, capped at about 30 s, reset by a `ready`. The last listener leaving stops it.
- On a `ready` after a reconnect, tell every listener to refetch through the callbacks it has now, since a socket does not resume as an `EventSource` did: such as a `change` that `affects` everything in `flaiover/src/lib/changes.ts`, and an `agent` with no story, if each listener already treats those as "ask again". Read the listeners in the 14 files that import `$lib/events` first. If none of the existing callbacks can say it without a component changing, stop and ask on the story's thread rather than add a callback.
- Where there is no `WebSocket`, `listen` does nothing, as it does now without an `EventSource`.

Rewrite `flaiover/src/lib/events.test.ts` against a fake `WebSocket`: one socket shared by two listeners of a project, `change` and `agent` reach both, the last `unlisten` closes it, a drop reconnects after the backoff, the reconnect tells both to refetch, and no `WebSocket` means no-op. jsdom has a real `WebSocket`, so the component tests that stub `EventSource` now stub `WebSocket` instead: `flaiover/src/lib/inbox.svelte.test.ts`, `flaiover/src/routes/board/board.svelte.test.ts`, `flaiover/src/routes/board/cardmenu.svelte.test.ts`, `flaiover/src/routes/board/lanemenu.svelte.test.ts`, and `flaiover/src/lib/components/StoryAgent.svelte.test.ts`.

It waits for T-1931, whose ADR fixes the path and the messages. It runs beside the server and harness tasks: no path in common.

## Done when

- No `.svelte` file that listens changes.
- `flai test` passes on the changed paths, vitest and lint.
- No file under `flaiover/src` but tests names `EventSource`.

## Notes

The server side comes in T-1934; this task's tests use a fake socket and need no server.
