---
id: S-0454
type: story
nature: remediation
title: The dashboard's live events reach the browser over a WebSocket, so open tabs no longer use up the browser's six connections to the dashboard and hang every page
status: review
owner: alex
created: 2026-10-20T22:10:33Z
updated: 2026-10-21T02:06:07Z
transitions:
  - to: ready
    at: 2026-10-21T00:36:37Z
    by: alex
  - to: in-progress
    at: 2026-10-21T00:51:51Z
    by: agent-S-0454
  - to: review
    at: 2026-10-21T01:43:05Z
    by: agent-S-0454
  - to: in-progress
    at: 2026-10-21T01:52:07Z
    by: flai
  - to: review
    at: 2026-10-21T02:06:07Z
    by: agent-S-0454
tags: [dashboard]
topics: [client-side, server-side, components]
touches: [design/adrs/README.md, flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/hooks.server.ts, flaiover/server.js, flaiover/vite.config.ts, flaiover/src/lib/events.ts, flaiover/src/lib/events.test.ts, flaiover/src/lib/inbox.svelte.test.ts, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/board/lanemenu.svelte.test.ts, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, flaiover/playwright.config.ts, flaiover/e2e/live-events.e2e.ts, scripts/flaiover-e2e.sh, Makefile, flaiover/src/routes/api/events/+server.ts, flaiover/src/lib/server/sse.ts, flaiover/src/lib/server/sse.test.ts, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, design/system/component-library.md, design/tech/sveltekit.md, docs/users/flaiover.md, docs/operators/index.md, design/adrs/0182-the-dashboard-s-live-events-reach-the-browser-over-a-websocket.md, docs/users/flai.md, flaiover/package.json, scripts/README.md, design/issues/I-0171-the-litellm-gateway-smoke-check-fails-because-claude-code-2-1-296-refuses-the-default-haiku-model-name-as-unrecognized-model.md, design/issues/summary.md, design/issues/I-0172-the-merge-queue-s-send-back-thread-quotes-only-the-first-line-of-a-failing-plain-tier-s-finding-and-keeps-no-log-of-the-batch-so-the-step-that-failed-cannot-be-told.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3929
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 1
      test_runs: 4
      hand_edits: 4
      work: 77
  models:
    - model: claude-opus-5-5
      input: 462
      output: 172647
      cache_read: 26229122
      cache_write: 696727
      cost: 12.9058
  strategic:
    - kind: orchestrator
      seconds: 4498
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 112
          output: 2744
          cache_read: 44019605
          cache_write: 57054
          cost: 11.0496
        - model: claude-sonnet-5-5
          input: 18
          output: 122
          cache_read: 357958
          cache_write: 62912
          cost: 0.3445
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: orchestrator
    at: 2026-10-20T22:24:25Z
  value: 37.5
  by: planner-S-0454
  at: 2026-10-20T22:25:11Z
forecast:
  duration: 1h30m
  delivery: 2026-10-21T02:18:00Z
  basis: "Its own forecast of 1h30m; 2nd in the pull order with an in-progress limit of 5, behind S-0445, S-0453 and S-0455."
  by: flai
  at: 2026-10-21T00:41:56Z
finalized:
  by: alex
  at: 2026-10-20T22:12:50Z
---
# S-0454 The dashboard's live events reach the browser over a WebSocket, so open tabs no longer use up the browser's six connections to the dashboard and hang every page

## Goal

The dashboard hangs in every tab once about six board tabs are open. Each tab holds one `EventSource` on `/api/events` (`flaiover/src/lib/events.ts`, one per project per tab since S-0154). The dashboard serves HTTP/1.1, and Chromium-based browsers keep at most six connections to one host. When server-sent event streams hold all six, every other request in every tab waits in the browser's queue with no timeout.

Seen on 2026-10-20 with Brave at 127.0.0.1:4242. The flaiover log showed every request handled in 0 to 145 ms, and no request from the browser after its last `/api/events` at 21:44:05. `ss` showed exactly six browser connections to port 4242. Each had sent its one request 7 to 42 minutes earlier and since then only received the server's pings.

This story moves the live events to a WebSocket. WebSockets do not count toward the browser's six HTTP/1.1 connections. Local TLS and HTTP/2, which would multiplex every request over one connection, are left for later.

## Acceptance criteria
- [x] An accepted ADR records that the dashboard's live events go to the browser over a WebSocket instead of server-sent events, why (the six-connection limit and the 2026-10-20 instance), and that local TLS and HTTP/2 are deferred.
- [x] The server upgrades the live-events WebSocket in `flaiover/server.js` and in the dev server's upgrade plugin in `flaiover/vite.config.ts`, beside the `/agent` upgrade (ADR-0029). It sends the messages `/api/events` sends now: `ready` once, `change` with `{ path, kind }`, and `agent` with `{ story }`, for the project the URL's `project` parameter names. It keeps the connection alive with a ping, and stops listening to the project's `Repo` when the socket closes.
- [x] The upgrade is authenticated and authorized as a request to `/api/events` is in `flaiover/src/hooks.server.ts`. Upgrades bypass SvelteKit's hooks, so an unauthenticated or forbidden upgrade is refused with 401 or 403 before the socket opens, and a test shows it.
- [x] `flaiover/src/lib/events.ts` keeps its `listen` API and still shares one socket per project per tab. After a drop it reconnects with a backoff, and on reconnecting it tells every listener to refetch, since a WebSocket does not resume as an `EventSource` did. No component that listens changes.
- [x] A Playwright test reproduces the hang: at least seven pages of one browser context following the same project's events, where a fetch from the last page answers. Run against the server-sent events before the change, it hangs.
- [x] The component harness (ADR-0172) treats a WebSocket to the app's own server as it treats an `EventSource` now: it opens no connection, fails as a refused one does, and lists it under Events.
- [x] The server-sent events endpoint and `flaiover/src/lib/server/sse.ts` are removed, or kept only if something else still needs them, as the ADR says. `design/system/flaiover-dashboard.md`'s API table, `design/system/dashboard-host-channel.md`, and `docs/users/flaiover.md` describe the WebSocket.

## Tasks
- T-1931 An accepted ADR records that the dashboard's live events go to the browser over a WebSocket, why, and that local TLS and HTTP/2 wait
- T-1934 The server upgrades the live-events WebSocket in server.js and the dev server, authenticated as /api/events is, and sends ready, change, and agent with a ping
- T-1936 events.ts follows the project's events over one WebSocket per project per tab, reconnects with a backoff, and tells every listener to refetch on reconnecting
- T-1937 The component harness opens no WebSocket to the app's own server, fails it as a refused one, and lists it under Events
- T-1938 The dashboard's design, its host channel design, and the user and operator docs describe the live-events WebSocket
- T-1939 A Playwright test opens seven pages of one browser context on one project's events and a fetch from the last page answers, where it hangs on the server-sent events
- T-1940 The server-sent events endpoint and sse.ts are removed, as the ADR says

## Notes

Chosen with the operator on 2026-10-20 over closing the stream while a tab is hidden, sharing one stream across tabs, and HTTP/2 behind local TLS.

`ws` is already a dependency, and `server.js` already routes `upgrade` to `globalThis.__flaioverAgentUpgrade` for `/agent`. That handler is the one to extend into a router by path.

Not in this story: the streaming responses of `/api/items/[id]/accept`, `/api/items/[id]/checks/tail`, and `/api/import` hold a connection too, but only while they run.

### Planning

Touches, 31 files and no folder. The story declared none, so `flai touches suggest` ran from the five files its goal and criteria name; the explorer read the code for the rest.

| Source | Touches |
|--------|---------|
| Criteria | `flaiover/server.js`, `flaiover/vite.config.ts`, `flaiover/src/hooks.server.ts`, `flaiover/src/lib/events.ts`, `flaiover/src/lib/server/sse.ts`, `design/system/flaiover-dashboard.md`, `design/system/dashboard-host-channel.md`, `docs/users/flaiover.md` |
| Co-change | `flaiover/src/lib/server/agent.ts` (33%), which sets the upgrade global and refuses any path but `/agent`; `flaiover/src/lib/server/sse.test.ts` (14%) |
| Layout | the new `flaiover/src/lib/server/events-socket.ts` and its test; `flaiover/src/routes/api/events/+server.ts`; `flaiover/src/lib/events.test.ts` and the five component tests that stub `EventSource`, since jsdom has a real `WebSocket`; the harness's `harness.mjs`, `integration_test.go`, and fixture `Fetcher.svelte`; `flaiover/playwright.config.ts`, the new `flaiover/e2e/live-events.e2e.ts`, `scripts/flaiover-e2e.sh`, and `Makefile` |
| Design | `design/adrs/0179-the-dashboard-s-live-events-reach-the-browser-over-a-websocket.md` (the number and slug are predicted) and `design/adrs/README.md`; `design/system/component-library.md` (ADR-0172's harness rule); `design/tech/sveltekit.md` (the `ws` row); `docs/operators/index.md` ("even with event streams open", which T-1934's shutdown keeps true) |

Left out though co-change lists them: `flaiover/src/lib/server/repo.ts` (38%), `flaiover/package.json`, and the lock file. The socket uses `Repo`'s existing `on` and `off`, and `ws` is already a dependency. `flaiover/src/lib/server/auth.ts` is called, not changed.

`docs/operators/index.md` is claimed by S-0445 and S-0451, both in progress, so this story is held while either is when it reaches ready. The line is one sentence of T-1938; if the hold costs more than it is worth, narrowing it out is the story agent's call.

Forecast 1h30m, delivery 2026-10-21T15:38:00Z. `flai forecast` gave 57m at size 38. Raised by 33m: flaiover has no Playwright test yet, `vite preview` loads no upgrade hook, so T-1939 has to run `node server.js` with a flai to dial it. T-1937's integration test also runs on flaiover's Playwright. The delivery moves by the same 33m.

Cost of delay 37.50 USD a week, as `flai cod` gives it from the orchestrator's input on TH-0535, 15m lost per 168h cycle at 150 USD an hour. Kept: it matches S-0445, a remediation with the same kind of stop and workaround.

Topic `components` added: T-1937 changes the component harness.
- 2026-10-21T01:52:07Z: moved to in-progress: merge queue batch 20261021-014447 fails with S-0454 at 6a4ae3b52fd6: smoke failed: smoke: smoke: template render and check
