---
id: T-2083
type: task
nature: feature
title: The dashboard's server accepts a request whose Origin names one of FLAIOVER_ALLOWED_HOSTS
status: done
parent: S-0489
owner: alex
created: 2026-10-23T18:07:21Z
updated: 2026-10-23T18:54:54Z
transitions:
  - to: ready
    at: 2026-10-23T18:50:45Z
    by: agent-S-0489
  - to: in-progress
    at: 2026-10-23T18:50:46Z
    by: agent-S-0489
  - to: done
    at: 2026-10-23T18:54:54Z
    by: agent-S-0489
stream: S-0489
tags: [flaiover, dashboard]
touches: [flaiover/src/lib/server/allowed-hosts.ts, flaiover/src/lib/server/allowed-hosts.test.ts, flaiover/src/lib/server/events-socket.ts, flaiover/src/lib/server/events-socket.test.ts, flaiover/server.js, design/system/flaiover-dashboard.md]
after: [T-2082]
usage:
  source: log
  seconds: 248
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 40
      output: 195
      cache_read: 1759870
      cache_write: 81505
      cost: 0.6221
---
# T-2083 The dashboard's server accepts a request whose Origin names one of FLAIOVER_ALLOWED_HOSTS

## Work

The one place the dashboard refuses a browser by its host today is `sameOrigin` in `flaiover/src/lib/server/events-socket.ts`: the `/api/events` upgrade answers `403 Forbidden` when the `Origin` header's host is not the request's `Host` header. Behind a tunnel or a reverse proxy that rewrites `Host` to the upstream address, or that serves the dashboard under a domain of its own, every page's live events are refused and the page shows "Reconnecting…". SvelteKit's own CSRF check (`csrf.checkOrigin`) covers form submissions only, which the dashboard does not use, and `/agent` refuses any `Origin` because flai is no browser; leave both as they are.

Add `flaiover/src/lib/server/allowed-hosts.ts`: it reads `FLAIOVER_ALLOWED_HOSTS` once, a comma-separated list of lower-case host names with optional ports as `flai dashboard` passes it, trimmed and with empties dropped, and answers whether a host is allowed: the request's own `Host`, always, or one of the list. `sameOrigin` takes it: an `Origin` whose host is the request's `Host` or an allowed host passes; a request with no `Origin` passes as before; a malformed `Origin` is still refused. Nothing else refuses by `Host`, so the server adds no new refusal: a request whose `Host` is not on the list and sends no `Origin` is served as today. Log the list once at start, in `flaiover/server.js` beside the listen line or where the registry is made, so an operator can see what the container runs with.

Tests: `allowed-hosts.test.ts` for the parsing and the match (case, port, empty); `events-socket.test.ts` for an upgrade from an allowed origin that passes, one from an unlisted origin that is still refused, and the same-host case unchanged. Record the check and the variable in `design/system/flaiover-dashboard.md` § Runtime, or beside the events socket's description.

This task waits for no other; it agrees with the flai task on the variable's name and shape.

## Done when

- With `FLAIOVER_ALLOWED_HOSTS=flow.example.com`, an `/api/events` upgrade with `Origin: https://flow.example.com` and a `Host` of the container's address is opened, and one with `Origin: http://elsewhere.example` is refused 403.
- Without the variable, every existing events-socket test passes unchanged.
- `design/system/flaiover-dashboard.md` says what the variable does, and `flai test` on the changed paths passes.

## Notes

Planned by the planner on 2026-10-23. If the operator saw a refusal with other words than the events socket's 403, the plan thread asks which; the story's agent then widens this task to that check.
