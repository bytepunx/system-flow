---
id: T-0374
type: task
nature: feature
title: The host package supervises child processes behind a control API on one machine-wide address
status: done
parent: S-0106
owner: alex
created: 2026-09-24T01:30:19Z
updated: 2026-09-24T01:35:48Z
transitions:
  - to: ready
    at: 2026-09-24T01:30:38Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T01:30:39Z
    by: system-flow
  - to: done
    at: 2026-09-24T01:35:48Z
    by: system-flow
stream: S-0106
tags: []
touches: [flai/internal/host]
---
# T-0374 The host package supervises child processes behind a control API on one machine-wide address

## Work

- New package `flai/internal/host`: `Run(ctx, Options)` binds the control address first (127.0.0.1:4241 by default), which is what makes one host per machine, then writes its state (`host/state.json` beside the config: pid, version, started, updated, addr, children) and a bearer token (`host/token`, 0600).
- A child is a named command line; the host starts it, waits for it, restarts one that ends unasked with a backoff, and stops every child (SIGTERM, a grace period, then SIGKILL) before it returns. The one `serve` child is always wanted; `mcp:<root>` children are wanted when asked for.
- Control API over HTTP, bearer token, requests with an `Origin` refused: `GET /_health` (no token: pid, version, config), `GET /status`, `POST /start|stop|restart {process}` (`serve`, `mcp`, `all`), `POST /mcp/ensure|release {root}`, `POST /upgrade` and `GET /check` through functions the caller supplies.
- Behaviour tests with helper processes: start, restart after a crash, stop leaves no child, a second host on the same address refuses, the token is required.

## Done when

- `go test -race ./internal/host/...` passes and lint is clean.

## Notes
