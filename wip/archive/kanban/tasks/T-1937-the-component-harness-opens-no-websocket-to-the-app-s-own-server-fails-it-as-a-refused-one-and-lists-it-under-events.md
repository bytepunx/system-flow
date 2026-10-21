---
id: T-1937
type: task
nature: remediation
title: The component harness opens no WebSocket to the app's own server, fails it as a refused one, and lists it under Events
status: done
parent: S-0454
owner: alex
created: 2026-10-20T22:21:42Z
updated: 2026-10-21T01:00:49Z
transitions:
  - to: ready
    at: 2026-10-21T00:53:38Z
    by: agent-S-0454
  - to: in-progress
    at: 2026-10-21T00:53:38Z
    by: agent-S-0454
  - to: done
    at: 2026-10-21T01:00:49Z
    by: agent-S-0454
stream: S-0454
tags: [components, dashboard]
touches: [flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, flai/internal/components/testdata/harness/src/lib/Fetcher.svelte, design/system/component-library.md]
after: [T-1931]
usage:
  source: log
  seconds: 431
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 12125
      cache_read: 1842039
      cache_write: 48930
      cost: 0.9064
---
# T-1937 The component harness opens no WebSocket to the app's own server, fails it as a refused one, and lists it under Events

## Work

In `flai/internal/components/harness/harness.mjs`, `answerTheAppsRequests()` (lines 114 to 163) replaces `window.fetch` and `window.EventSource` (ADR-0172). Replace `window.WebSocket` the same way. A socket is for the app's own server when its host is the dev server's own (`location.host`) and its path passes the rule `forTheApp` (lines 87 to 92) applies; `forTheApp` compares origins, and a `ws:` URL's origin is never the page's, so map `ws:` to `http:` and `wss:` to `https:` before it, or compare hosts. Vite's own HMR socket, under `/@` or on its own path with the `vite-hmr` protocol, goes through. A socket for the app opens no connection: it fires `error`, then `close` with code 1006, and is `readyState` 3, closed. The harness posts `{type: "event", component, name: "WebSocket", args: [path]}`, the path with its query. Bring the header comment (lines 28 to 34) up to date.

Extend the fixture `flai/internal/components/testdata/harness/src/lib/Fetcher.svelte` with a socket to `/api/events?project=x` beside its `EventSource`, and `checkFetches` in `flai/internal/components/integration_test.go` (lines 802 to 873) with its state and the event `WebSocket ["/api/events?project=x"]`, beside `EventSource` at line 868. Keep the `EventSource` case: another project's components may still use one.

In `design/system/component-library.md` § The harness, and the `event` row (line 177), name `WebSocket` beside `fetch` and `EventSource`.

It waits for T-1931, whose ADR refines ADR-0172 with this. It runs beside the server and client tasks: no path in common.

## Done when

- The integration test passes with flaiover's Playwright installed (`scripts/integration.sh`), and a framed component's socket to its own server reaches nothing.
- Vite's HMR still connects in the harness.

## Notes

The harness's runner takes Playwright from `flaiover` (`integration_test.go` lines 183 to 205), so run the test after `scripts/flaiover-install.sh`. The line in `docs/users/flaiover.md` § Components that names `EventSource` is T-1938's.
