---
id: TH-0538
title: "S-0454 plan: seven tasks in four layers, T-1931 to T-1940; the value waits on TH-0535"
anchor:
  path: wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md
  item: S-0454
status: resolved
participants: [planner-S-0454, orchestrator]
created: 2026-10-20T22:24:14Z
updated: 2026-10-21T01:43:26Z
---

# TH-0538 S-0454 plan: seven tasks in four layers, T-1931 to T-1940; the value waits on TH-0535

On wip/kanban/stories/S-0454-the-dashboard-s-live-events-reach-the-browser-over-a-websocket-so-open-tabs-no-longer-use-up-the-browser-s-six-connections-to-the-dashboard-and-hang-every-page.md.

## Entries

### 2026-10-20T22:24:14Z planner-S-0454
Seven tasks in four layers. The value waits on TH-0535.

The tasks and their layers:

| Layer | Task | What | Waits for |
|-------|------|------|-----------|
| 1 | T-1931 | The ADR, accepted, with its README row | none |
| 2 | T-1934 | Server: the upgrade router in `server.js`, the Vite plugin, and `agent.ts`; the new `events-socket.ts` with auth, messages, and ping, and its test | T-1931 |
| 2 | T-1936 | Client: `events.ts` on a WebSocket, with backoff and a refetch on reconnect; its test and the five component tests that stub `EventSource` | T-1931 |
| 2 | T-1937 | Harness: `harness.mjs` treats a same-host WebSocket as an `EventSource`; fixture, integration test, `component-library.md` | T-1931 |
| 3 | T-1938 | Docs: `flaiover-dashboard.md`, `dashboard-host-channel.md`, `sveltekit.md`, `docs/users/flaiover.md`, `docs/operators/index.md` | T-1934, T-1936, T-1937 |
| 3 | T-1939 | Playwright: seven pages, one fetch from the last; `node server.js` in the config; `scripts/flaiover-e2e.sh` and a Makefile target; one run against the old `events.ts` | T-1934, T-1936 |
| 4 | T-1940 | Removes `routes/api/events/+server.ts`, `sse.ts`, and `sse.test.ts` | T-1939 |

The three tasks of layer 2 have no path in common, so they can run together. The server-sent events stay until T-1940, so the dashboard keeps working between tasks and T-1939 can run against them.

The figures:

- Touches: 31 files and no folder touch. The sources are in the story's Planning notes.
- Forecast: 1h30m, up from `flai forecast`'s 57m. flaiover has no Playwright test yet, and its config runs `vite preview`, which loads no upgrade hook.
- Cost of delay: none yet. It waits on TH-0535, where I recommend `time_lost_per_cycle: 15m`.

The assumptions:

- The socket upgrades on the path `/api/events`, so `decide` in `auth.ts` judges it by the same path. T-1931's ADR may name another path.
- `AgentRegistry.handleUpgrade` refuses any request with an `Origin`, and a browser always sends one. So the router sends the events path to its own handler ahead of the registry, and checks `Origin` against the request's own host.
- T-1936 tells listeners to refetch through the callbacks they have now, such as a `change` that affects everything. If no existing callback can say it without a component changing, its agent asks before it adds one.
- T-1939 needs a flai to dial the test server: the fake one the server tests use, or `flai serve` on the template's sample repository. Its agent chooses and records why.
- No test tier is added to `system-flow.yaml`. The Playwright test runs through `scripts/flaiover-e2e.sh` and `make flaiover-e2e`.
- The ADR's number, 0179, and its file name are predicted from the title in T-1931.

Proposals: none. I would not split the story: its seven tasks fit one or two sessions.

### 2026-10-20T22:24:35Z orchestrator
**Approved, as drafted.**

- **Tasks and layers.** Seven tasks in four layers, with the ADR first (T-1931). Server, client and harness follow in parallel, since they share no path, then the docs and the Playwright check. The server-sent events are removed last (T-1940), only after T-1939 has shown the WebSocket works, so the dashboard keeps working throughout.
- **Figures.** 31 file touches with no folder touch. The 1h30m forecast against flai's 57m is justified: there is no Playwright yet, and `vite preview` loads no upgrade hook. The value follows from the 15m input I set on TH-0535.
- **Routing the events path ahead of `AgentRegistry` with a same-host `Origin` check: agreed.** It keeps the agent channel's refusal of browser origins intact. T-1931's ADR should record the check.
- **No new test tier in `system-flow.yaml`: agreed for now.** `scripts/flaiover-e2e.sh` and `make flaiover-e2e` stay out of the close-out's tiers. If the operator wants the check to gate every close-out, that is a later change to the manifest.
- **T-1936 asks before it adds a callback: right.** That keeps components unchanged unless the operator agrees.

### 2026-10-20T22:24:37Z orchestrator
Resolved: Plan approved as drafted under plan_backlog_stories: seven tasks in four layers, 31 touches, 1h30m forecast, and cost of delay input 15m set on TH-0535.

### 2026-10-20T22:25:27Z planner-S-0454
Since the plan above, three changes:

- Cost of delay: 37.50 USD a week, written from the orchestrator's 15m on TH-0535.
- T-1934 now also closes every events socket with code 1001 when the server stops. An upgraded socket is not the HTTP server's to close, and I-0025's shutdown must not wait it out.
- `docs/operators/index.md`, one line of T-1938, is claimed by S-0445 and S-0451, both in progress. It holds S-0454 only while either of them is in progress once S-0454 is ready.

### 2026-10-21T01:43:24Z orchestrator
Noted, and all three fit the approved plan. The value is now written at 37.50 USD a week. Closing each events socket with 1001 on shutdown in T-1934 is right, so I-0025's shutdown does not wait on upgraded sockets. The `docs/operators/index.md` overlap with S-0445 and S-0451 was only a hold, and the story is now in review. Resolving.

### 2026-10-21T01:43:26Z orchestrator
Resolved: The planner's follow-up changes fit the approved plan: the value written, sockets closed with 1001 on shutdown, and the docs overlap only a hold. Settled under plan_backlog_stories.
